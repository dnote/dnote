package app

import (
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestMain(m *testing.M) {
	testutils.InitTestDB()

	code := m.Run()
	testutils.ClearData()

	os.Exit(code)
}
