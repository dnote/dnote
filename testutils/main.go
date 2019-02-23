// Package testutils provides utilities used in tests
package testutils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
)

// InitEnv sets up a test env and returns a new dnote context
func InitEnv(t *testing.T, relPath string, relFixturePath string, migrated bool) infra.DnoteCtx {
	path, err := filepath.Abs(relPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "pasrsing path"))
	}

	os.Setenv("DNOTE_HOME_DIR", path)
	ctx, err := infra.NewCtx("", "")
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting new ctx"))
	}

	// set up directory
	if err := os.MkdirAll(ctx.DnoteDir, 0755); err != nil {
		t.Fatal(err)
	}

	// set up db
	b := ReadFileAbs(relFixturePath)
	setupSQL := string(b)

	db := ctx.DB
	if _, err := db.Exec(setupSQL); err != nil {
		t.Fatal(errors.Wrap(err, "running schema sql"))
	}

	if migrated {
		// mark migrations as done. When adding new migrations, bump the numbers here.
		if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", infra.SystemSchema, 9); err != nil {
			t.Fatal(errors.Wrap(err, "inserting schema"))
		}

		if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", infra.SystemRemoteSchema, 1); err != nil {
			t.Fatal(errors.Wrap(err, "inserting remote schema"))
		}
	}

	return ctx
}

// Login simulates a logged in user by inserting credentials in the local database
func Login(t *testing.T, ctx *infra.DnoteCtx) {
	db := ctx.DB

	MustExec(t, "inserting sessionKey", db, "INSERT INTO system (key, value) VALUES (?, ?)", infra.SystemSessionKey, "someSessionKey")
	MustExec(t, "inserting sessionKeyExpiry", db, "INSERT INTO system (key, value) VALUES (?, ?)", infra.SystemSessionKeyExpiry, time.Now().Add(24*time.Hour).Unix())
	MustExec(t, "inserting cipherKey", db, "INSERT INTO system (key, value) VALUES (?, ?)", infra.SystemCipherKey, "QUVTMjU2S2V5LTMyQ2hhcmFjdGVyczEyMzQ1Njc4OTA=")

	ctx.SessionKey = "someSessionKey"
	ctx.SessionKeyExpiry = time.Now().Add(24 * time.Hour).Unix()
	ctx.CipherKey = []byte("AES256Key-32Characters1234567890")
}

// TeardownEnv cleans up the test env represented by the given context
func TeardownEnv(ctx infra.DnoteCtx) {
	ctx.DB.Close()

	if err := os.RemoveAll(ctx.DnoteDir); err != nil {
		panic(err)
	}
}

// CopyFixture writes the content of the given fixture to the filename inside the dnote dir
func CopyFixture(ctx infra.DnoteCtx, fixturePath string, filename string) {
	fp, err := filepath.Abs(fixturePath)
	if err != nil {
		panic(err)
	}
	dp, err := filepath.Abs(filepath.Join(ctx.DnoteDir, filename))
	if err != nil {
		panic(err)
	}

	err = utils.CopyFile(fp, dp)
	if err != nil {
		panic(err)
	}
}

// WriteFile writes a file with the given content and  filename inside the dnote dir
func WriteFile(ctx infra.DnoteCtx, content []byte, filename string) {
	dp, err := filepath.Abs(filepath.Join(ctx.DnoteDir, filename))
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(dp, content, 0644); err != nil {
		panic(err)
	}
}

// ReadFile reads the content of the file with the given name in dnote dir
func ReadFile(ctx infra.DnoteCtx, filename string) []byte {
	path := filepath.Join(ctx.DnoteDir, filename)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return b
}

// ReadFileAbs reads the content of the file with the given file path by resolving
// it as an absolute path
func ReadFileAbs(relpath string) []byte {
	fp, err := filepath.Abs(relpath)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	return b
}

func checkEqual(a interface{}, b interface{}, message string) (bool, string) {
	if a == b {
		return true, ""
	}

	var m string
	if len(message) == 0 {
		m = fmt.Sprintf("%v != %v", a, b)
	} else {
		m = message
	}
	errorMessage := fmt.Sprintf("%s. Actual: %+v. Expected: %+v.", m, a, b)

	return false, errorMessage
}

// AssertEqual errors a test if the actual does not match the expected
func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Error(m)
	}
}

// AssertEqualf fails a test if the actual does not match the expected
func AssertEqualf(t *testing.T, a interface{}, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Fatal(m)
	}
}

// AssertNotEqual fails a test if the actual matches the expected
func AssertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v == %v", a, b)
	}
	t.Errorf("%s. Actual: %+v. Expected: %+v.", message, a, b)
}

// AssertDeepEqual fails a test if the actual does not deeply equal the expected
func AssertDeepEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if reflect.DeepEqual(a, b) {
		return
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Errorf("%s.\nActual:   %+v.\nExpected: %+v.", message, a, b)
}

// ReadJSON reads JSON fixture to the struct at the destination address
func ReadJSON(path string, destination interface{}) {
	var dat []byte
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, "Failed to load fixture payload"))
	}
	if err := json.Unmarshal(dat, destination); err != nil {
		panic(errors.Wrap(err, "Failed to get event"))
	}
}

// IsEqualJSON deeply compares two JSON byte slices
func IsEqualJSON(s1, s2 []byte) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	if err := json.Unmarshal(s1, &o1); err != nil {
		return false, errors.Wrap(err, "unmarshalling first  JSON")
	}
	if err := json.Unmarshal(s2, &o2); err != nil {
		return false, errors.Wrap(err, "unmarshalling second JSON")
	}

	return reflect.DeepEqual(o1, o2), nil
}

// MustExec executes the given SQL query and fails a test if an error occurs
func MustExec(t *testing.T, message string, db *infra.DB, query string, args ...interface{}) sql.Result {
	result, err := db.Exec(query, args...)
	if err != nil {
		t.Fatal(errors.Wrap(errors.Wrap(err, "executing sql"), message))
	}

	return result
}

// MustScan scans the given row and fails a test in case of any errors
func MustScan(t *testing.T, message string, row *sql.Row, args ...interface{}) {
	err := row.Scan(args...)
	if err != nil {
		t.Fatal(errors.Wrap(errors.Wrap(err, "scanning a row"), message))
	}
}

// NewDnoteCmd returns a new Dnote command and a pointer to stderr
func NewDnoteCmd(ctx infra.DnoteCtx, binaryName string, arg ...string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer, error) {
	var stderr, stdout bytes.Buffer

	binaryPath, err := filepath.Abs(binaryName)
	if err != nil {
		return &exec.Cmd{}, &stderr, &stdout, errors.Wrap(err, "getting the absolute path to the test binary")
	}

	cmd := exec.Command(binaryPath, arg...)
	cmd.Env = []string{fmt.Sprintf("DNOTE_DIR=%s", ctx.DnoteDir), fmt.Sprintf("DNOTE_HOME_DIR=%s", ctx.HomeDir)}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	return cmd, &stderr, &stdout, nil
}

// RunDnoteCmd runs a dnote command
func RunDnoteCmd(t *testing.T, ctx infra.DnoteCtx, binaryName string, arg ...string) {
	t.Logf("running: %s %s", binaryName, strings.Join(arg, " "))

	cmd, stderr, stdout, err := NewDnoteCmd(ctx, binaryName, arg...)
	if err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrap(err, "getting command").Error())
	}

	cmd.Env = append(cmd.Env, "DNOTE_DEBUG=1")

	if err := cmd.Run(); err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrapf(err, "running command %s", stderr.String()))
	}

	// Print stdout if and only if test fails later
	t.Logf("\n%s", stdout)
}

// WaitDnoteCmd runs a dnote command and waits until the command is exited
func WaitDnoteCmd(t *testing.T, ctx infra.DnoteCtx, runFunc func(io.WriteCloser) error, binaryName string, arg ...string) {
	t.Logf("running: %s %s", binaryName, strings.Join(arg, " "))

	cmd, stderr, stdout, err := NewDnoteCmd(ctx, binaryName, arg...)
	if err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrap(err, "getting command").Error())
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrap(err, "getting stdin %s"))
	}
	defer stdin.Close()

	// Start the program
	err = cmd.Start()
	if err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrap(err, "starting command"))
	}

	err = runFunc(stdin)
	if err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrap(err, "running with stdin"))
	}

	err = cmd.Wait()
	if err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrapf(err, "running command %s", stderr.String()))
	}

	// Print stdout if and only if test fails later
	t.Logf("\n%s", stdout)
}

// UserConfirm simulates confirmation from the user by writing to stdin
func UserConfirm(stdin io.WriteCloser) error {
	// confirm
	if _, err := io.WriteString(stdin, "y\n"); err != nil {
		return errors.Wrap(err, "indicating confirmation in stdin")
	}

	return nil
}

// MustMarshalJSON marshalls the given interface into JSON.
// If there is any error, it fails the test.
func MustMarshalJSON(t *testing.T, v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("%s: marshalling data", t.Name())
	}

	return b
}

// MustUnmarshalJSON marshalls the given interface into JSON.
// If there is any error, it fails the test.
func MustUnmarshalJSON(t *testing.T, data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		t.Fatalf("%s: unmarshalling data", t.Name())
	}
}
