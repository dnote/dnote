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
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

func command(binary string, args []string, entryPoint string) *exec.Cmd {
	cmd := exec.Command(binary, args...)

	// Notice this change.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Dir = entryPoint
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Using Start and not Run.
	err := cmd.Start()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return cmd
}

func execCmd(task string, watchDir string) *exec.Cmd {
	parts := strings.Fields(task)

	return command(parts[0], parts[1:], watchDir)
}

func main() {
	w := watcher.New()
	w.IgnoreHiddenFiles(true)
	w.SetMaxEvents(1)

	watchDir := os.Args[1]
	task := os.Args[2]

	var e *exec.Cmd

	go func() {
		for {
			select {
			case <-w.Event:
				log.Println("Change detected. Restarting server...")

				// Killing the process here.
				if e != nil {
					pgid, err := syscall.Getpgid(e.Process.Pid)
					if err == nil {
						syscall.Kill(-pgid, syscall.SIGKILL)
					}
					e.Wait()
				}

				// Starting it again here or starting for the first time.
				e = execCmd(task, watchDir)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(watchDir); err != nil {
		log.Fatalln(errors.Wrap(err, "watching the given pattern"))
	}

	e = execCmd(task, watchDir)

	log.Printf("watching %d files", len(w.WatchedFiles()))
	if err := w.Start(time.Millisecond * 1000); err != nil {
		log.Fatalln(errors.Wrap(err, "starting watcher"))
	}
}
