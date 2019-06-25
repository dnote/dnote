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

package main

import (
	"flag"
	"log"
	"os"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job/digest"
	"github.com/dnote/dnote/pkg/server/mailer"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

var (
	emailTemplateDir = flag.String("emailTemplateDir", "../mailer/templates/src", "the path to the template directory")
)

func init() {
	// Load env
	if os.Getenv("GO_ENV") == "PRODUCTION" {
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	} else {
		err := godotenv.Load(".env.dev")
		if err != nil {
			panic(err)
		}
	}
}

func scheduleJob(c *cron.Cron, spec string, cmd func()) {
	s, err := cron.ParseStandard(spec)
	if err != nil {
		panic(errors.Wrap(err, "parsing schedule"))
	}

	c.Schedule(s, cron.FuncJob(cmd))
}

func main() {
	flag.Parse()

	mailer.InitTemplates(*emailTemplateDir)

	database.InitDB()
	defer database.CloseDB()

	// Run jobs on initial start
	log.Println("Job is running")

	// Schedule jobs
	c := cron.New()

	scheduleJob(c, "0 20 * * 5", func() { digest.Send() })

	c.Start()

	// Block forever
	select {}
}
