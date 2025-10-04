/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package job

import (
	slog "log"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/mailer"
	"gorm.io/gorm"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

var (
	// ErrEmptyDB is an error for missing database connection in the app configuration
	ErrEmptyDB = errors.New("No database connection was provided")
	// ErrEmptyClock is an error for missing clock in the app configuration
	ErrEmptyClock = errors.New("No clock was provided")
	// ErrEmptyWebURL is an error for missing WebURL content in the app configuration
	ErrEmptyWebURL = errors.New("No WebURL was provided")
	// ErrEmptyEmailTemplates is an error for missing EmailTemplates content in the app configuration
	ErrEmptyEmailTemplates = errors.New("No EmailTemplate store was provided")
	// ErrEmptyEmailBackend is an error for missing EmailBackend content in the app configuration
	ErrEmptyEmailBackend = errors.New("No EmailBackend was provided")
)

// Runner is a configuration for job
type Runner struct {
	DB           *gorm.DB
	Clock        clock.Clock
	EmailTmpl    mailer.Templates
	EmailBackend mailer.Backend
	Config       config.Config
}

// NewRunner returns a new runner
func NewRunner(db *gorm.DB, c clock.Clock, t mailer.Templates, b mailer.Backend, config config.Config) (Runner, error) {
	ret := Runner{
		DB:           db,
		EmailTmpl:    t,
		EmailBackend: b,
		Clock:        c,
		Config:       config,
	}

	if err := ret.validate(); err != nil {
		return Runner{}, errors.Wrap(err, "validating runner configuration")
	}

	return ret, nil
}

func (r *Runner) validate() error {
	if r.DB == nil {
		return ErrEmptyDB
	}
	if r.Clock == nil {
		return ErrEmptyClock
	}
	if r.EmailTmpl == nil {
		return ErrEmptyEmailTemplates
	}
	if r.EmailBackend == nil {
		return ErrEmptyEmailBackend
	}
	if r.Config.WebURL == "" {
		return ErrEmptyWebURL
	}

	return nil
}

func scheduleJob(c *cron.Cron, spec string, cmd func()) {
	s, err := cron.ParseStandard(spec)
	if err != nil {
		panic(errors.Wrap(err, "parsing schedule"))
	}

	c.Schedule(s, cron.FuncJob(cmd))
}

func (r *Runner) schedule(ch chan error) {
	// Schedule jobs
	cr := cron.New()
	cr.Start()

	ch <- nil

	// Block forever
	select {}
}

// Do starts the background tasks in a separate goroutine that runs forever
func (r *Runner) Do() error {
	// validate
	if err := r.validate(); err != nil {
		return errors.Wrap(err, "validating job configurations")
	}

	ch := make(chan error)
	go r.schedule(ch)
	if err := <-ch; err != nil {
		return errors.Wrap(err, "scheduling jobs")
	}

	slog.Println("Started background tasks")

	return nil
}
