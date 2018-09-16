// Package testutils provides utilities used in tests
package testutils

import (
	"database/sql"
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

// InitEnv sets up a test env and returns a new dnote context
func InitEnv(relPath string, relFixturePath string) infra.DnoteCtx {
	path, err := filepath.Abs(relPath)
	if err != nil {
		panic(errors.Wrap(err, "pasrsing path").Error())
	}

	os.Setenv("DNOTE_HOME_DIR", path)
	ctx, err := infra.NewCtx("", "")
	if err != nil {
		panic(errors.Wrap(err, "getting new ctx").Error())
	}

	// set up directory and db
	if err := os.MkdirAll(ctx.DnoteDir, 0755); err != nil {
		panic(err)
	}

	b := ReadFileAbs(relFixturePath)
	setupSQL := string(b)

	db := ctx.DB
	_, err = db.Exec(setupSQL)
	if err != nil {
		panic(errors.Wrap(err, "running schema sql").Error())
	}

	return ctx
}

// TeardownEnv cleans up the test env represented by the given context
func TeardownEnv(ctx infra.DnoteCtx) {
	ctx.DB.Close()

	if err := os.RemoveAll(ctx.DnoteDir); err != nil {
		panic(err)
	}
}

// WriteFile writes the content of the given fixture to the filename inside the dnote dir
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

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	errorMessage := fmt.Sprintf("%s. Actual: %+v. Expected: %+v.", message, a, b)

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
func MustExec(t *testing.T, message string, db *sql.DB, query string, args ...interface{}) sql.Result {
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
