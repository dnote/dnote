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

	"github.com/dnote/dnote/pkg/clock"
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

func schedule(ch chan error) {
	cl := clock.New()

	// Schedule jobs
	c := cron.New()
	scheduleJob(c, "* * * * *", func() { repetition.Do(cl) })
	c.Start()

	ch <- nil

	// Block forever
	select {}
}

// Run starts the background tasks and blocks forever.
func Run() error {
	if err := checkEnvironment(); err != nil {
		return errors.Wrap(err, "checking environment variables")
	}

	ch := make(chan error)
	go schedule(ch)
	if err := <-ch; err != nil {
		return errors.Wrap(err, "scheduling jobs")
	}

	log.Println("Started background tasks")

	return nil
}
