package ctx

import (
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/jinzhu/gorm"
)

// Ctx holds information that a job needs in order to perform
type Ctx struct {
	DB           *gorm.DB
	Clock        clock.Clock
	EmailTmpl    mailer.Templates
	EmailBackend mailer.Backend
}
