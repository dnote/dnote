package handlers

import (
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestMain(m *testing.M) {
	testutils.InitTestDB()
	templatePath := os.Getenv("DNOTE_TEST_EMAIL_TEMPLATE_DIR")
	mailer.InitTemplates(&templatePath)

	code := m.Run()
	testutils.ClearData()

	os.Exit(code)
}
