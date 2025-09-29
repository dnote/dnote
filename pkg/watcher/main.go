/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

// splitCommandParts splits the given commad string at space, except
// when inside a double quotation mark.
func splitCommandParts(cmd string) []string {
	re := regexp.MustCompile(`\r?\n`)
	s := re.ReplaceAllString(cmd, " ")

	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '

	fields, err := r.Read()
	if err != nil {
		panic(err)
	}

	return fields
}

func command(binary string, args []string, entryPoint string) *exec.Cmd {
	log.Printf("executing command: %s %s", binary, args)
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
	parts := splitCommandParts(task)

	return command(parts[0], parts[1:], watchDir)
}

var task, context, ignore string

func init() {
	flag.StringVar(&task, "task", "", "the command to execute")
	flag.StringVar(&context, "context", ".", "the file or directory from which to execute the task")
	flag.StringVar(&ignore, "ignore", ".", "the file or directory to ignore")

	flag.Parse()

	if task == "" {
		log.Println("task was not provided. Exiting the watcher...")
		os.Exit(1)
	}
}

func killCmdProcess(cmd *exec.Cmd) {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGKILL)
	}
}

func main() {
	w := watcher.New()
	w.IgnoreHiddenFiles(true)
	w.SetMaxEvents(1)

	targets := flag.Args()

	var e *exec.Cmd

	go func() {
		for {
			select {
			case <-w.Event:
				log.Println("Change detected. Restarting server...")

				// Killing the process here.
				if e != nil {
					killCmdProcess(e)
					e.Wait()
				}

				// Starting it again here or starting for the first time.
				e = execCmd(task, context)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if ignore != "" {
		files := strings.Split(ignore, ",")
		for _, file := range files {
			if err := w.Ignore(file); err != nil {
				log.Fatalln(errors.Wrapf(err, "ignoring %s", file))
			}
		}
	}

	for _, target := range targets {
		if err := w.AddRecursive(target); err != nil {
			log.Fatalln(errors.Wrap(err, "watching the given pattern"))
		}
	}

	e = execCmd(task, context)

	// watch for quit signals and kill the child process
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChan
		killCmdProcess(e)
		os.Exit(0)
	}()

	log.Printf("watching %d files", len(w.WatchedFiles()))
	if err := w.Start(time.Millisecond * 1000); err != nil {
		log.Fatalln(errors.Wrap(err, "starting watcher"))
	}
}
