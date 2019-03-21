package login

import (
	"encoding/base64"
	"strconv"

	"github.com/dnote/cli/client"
	"github.com/dnote/cli/core"
	"github.com/dnote/cli/crypt"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
  dnote login`

// NewCmd returns a new login command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to dnote server",
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

// Do dervies credentials on the client side and requests a session token from the server
func Do(ctx infra.DnoteCtx, email, password string) error {
	presigninResp, err := client.GetPresignin(ctx, email)
	if err != nil {
		return errors.Wrap(err, "getting presiginin")
	}

	masterKey, authKey, err := crypt.MakeKeys([]byte(password), []byte(email), presigninResp.Iteration)
	if err != nil {
		return errors.Wrap(err, "making keys")
	}

	authKeyB64 := base64.StdEncoding.EncodeToString(authKey)
	signinResp, err := client.Signin(ctx, email, authKeyB64)
	if err != nil {
		return errors.Wrap(err, "requesting session")
	}

	cipherKeyDec, err := crypt.AesGcmDecrypt(masterKey, signinResp.CipherKeyEnc)
	if err != nil {
		return errors.Wrap(err, "decrypting cipher key")
	}

	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if err := core.UpsertSystem(tx, infra.SystemCipherKey, cipherKeyDec); err != nil {
		return errors.Wrap(err, "saving enc key")
	}
	if err := core.UpsertSystem(tx, infra.SystemSessionKey, signinResp.Key); err != nil {
		return errors.Wrap(err, "saving session key")
	}
	if err := core.UpsertSystem(tx, infra.SystemSessionKeyExpiry, strconv.FormatInt(signinResp.ExpiresAt, 10)); err != nil {
		return errors.Wrap(err, "saving session key")
	}

	tx.Commit()

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var email, password string
		if err := utils.PromptInput("email", &email); err != nil {
			return errors.Wrap(err, "getting email input")
		}
		if email == "" {
			return errors.New("Email is empty")
		}

		if err := utils.PromptPassword("password", &password); err != nil {
			return errors.Wrap(err, "getting password input")
		}
		if password == "" {
			return errors.New("Password is empty")
		}

		err := Do(ctx, email, password)
		if errors.Cause(err) == client.ErrInvalidLogin {
			log.Error("wrong login\n")
			return nil
		} else if err != nil {
			return errors.Wrap(err, "logging in")
		}

		log.Success("logged in\n")

		return nil
	}

}
