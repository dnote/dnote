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
	"os"

	"github.com/dnote/dnote/pkg/server/job/ctx"
	"github.com/dnote/dnote/pkg/server/job/repetition"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

func scheduleJob(c *cron.Cron, spec string, cmd func()) {
	s, err := cron.ParseStandard(spec)
	if err != nil {
		panic(errors.Wrap(err, "parsing schedule"))
	}

	c.Schedule(s, cron.FuncJob(cmd))
}

func checkEnvironment() error {
	if os.Getenv("WebURL") == "" {
		return errors.New("WebURL is empty")
	}

	return nil
}

func schedule(c ctx.Ctx, ch chan error) {
	// Schedule jobs
	cr := cron.New()
	scheduleJob(cr, "* * * * *", func() { repetition.Do(c) })
	cr.Start()

	ch <- nil

	// Block forever
	select {}
}

// Run starts the background tasks in a separate goroutine that runs forever
func Run(c ctx.Ctx) error {
	if err := checkEnvironment(); err != nil {
		return errors.Wrap(err, "checking environment variables")
	}

	ch := make(chan error)
	go schedule(c, ch)
	if err := <-ch; err != nil {
		return errors.Wrap(err, "scheduling jobs")
	}

	log.Println("Started background tasks")

	return nil
}
