/* Copyright (C) 2019 Monomax Software Pty Ltd
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
	"log"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/job/repetition"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

// Job is a configuration for job
type Job struct {
	DB           *gorm.DB
	Clock        clock.Clock
	EmailTmpl    mailer.Templates
	EmailBackend mailer.Backend
	WebURL       string
}

func (j *Job) validate() error {
	if j.DB == nil {
		return errors.New("DB is not provided")
	}
	if j.Clock == nil {
		return errors.New("Clock is not provided")
	}
	if j.EmailTmpl == nil {
		return errors.New("EmailTmpl is not provided")
	}
	if j.EmailBackend == nil {
		return errors.New("EmailBackend is not provided")
	}
	if j.WebURL == "" {
		return errors.New("WebURL is not provided")
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

func (j *Job) schedule(ch chan error) {
	// Schedule jobs
	cr := cron.New()
	scheduleJob(cr, "* * * * *", func() { j.DoRepetition() })
	cr.Start()

	ch <- nil

	// Block forever
	select {}
}

// Run starts the background tasks in a separate goroutine that runs forever
func (j *Job) Run() error {
	if err := j.validate(); err != nil {
		return errors.Wrap(err, "validating job configurations")
	}

	ch := make(chan error)
	go j.schedule(ch)
	if err := <-ch; err != nil {
		return errors.Wrap(err, "scheduling jobs")
	}

	log.Println("Started background tasks")

	return nil
}

// DoRepetition creates spaced repetitions and delivers the results based on the rules
func (j *Job) DoRepetition() error {
	p := repetition.Params{
		DB:           j.DB,
		Clock:        j.Clock,
		EmailTmpl:    j.EmailTmpl,
		EmailBackend: j.EmailBackend,
	}
	if err := repetition.Do(p); err != nil {
		return errors.Wrap(err, "performing repetition job")
	}

	return nil
}
