package web

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func TestInit(t *testing.T) {
	mockIndexHTML := []byte("<html></html>")
	mockRobotsTxt := []byte("Allow: *")
	mockServiceWorkerJs := []byte("function() {}")
	mockStaticFileSystem := http.Dir(".")

	testCases := []struct {
		ctx         Context
		expectedErr error
	}{
		{
			ctx: Context{
				DB:               &gorm.DB{},
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: nil,
		},
		{
			ctx: Context{
				DB:               nil,
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyDB,
		},
		{
			ctx: Context{
				DB:               &gorm.DB{},
				IndexHTML:        nil,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyIndexHTML,
		},
		{
			ctx: Context{
				DB:               &gorm.DB{},
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        nil,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyRobotsTxt,
		},
		{
			ctx: Context{
				DB:               &gorm.DB{},
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  nil,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyServiceWorkerJS,
		},
		{
			ctx: Context{
				DB:               &gorm.DB{},
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: nil,
			},
			expectedErr: ErrEmptyStaticFileSystem,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			_, err := Init(tc.ctx)

			assert.Equal(t, errors.Cause(err), tc.expectedErr, "error mismatch")
		})
	}
}
