// Package testutils provides utilities used in tests
package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
)

func InitCtx(relPath string) infra.DnoteCtx {
	path, err := filepath.Abs(relPath)
	if err != nil {
		panic(err)
	}

	ctx := infra.DnoteCtx{
		HomeDir:  path,
		DnoteDir: fmt.Sprintf("%s/.dnote", path),
	}

	return ctx
}

func WriteFile(ctx infra.DnoteCtx, fixturePath string, filename string) {
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

func WriteFileWithContent(ctx infra.DnoteCtx, content []byte, filename string) {
	dp, err := filepath.Abs(filepath.Join(ctx.DnoteDir, filename))
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(dp, content, 0644); err != nil {
		panic(err)
	}
}

func ReadFile(ctx infra.DnoteCtx, filename string) []byte {
	path := filepath.Join(ctx.DnoteDir, filename)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return b
}

func ReadFileAbs(filename string) []byte {
	fp, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	return b
}

func SetupTmp(ctx infra.DnoteCtx) {
	if err := os.MkdirAll(ctx.DnoteDir, 0755); err != nil {
		panic(err)
	}
}

func ClearTmp(ctx infra.DnoteCtx) {
	if err := os.RemoveAll(ctx.DnoteDir); err != nil {
		panic(err)
	}
}

// AssertEqual fails a test if the actual does not match the expected
func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Errorf("%s. Actual: %+v. Expected: %+v.", message, a, b)
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
