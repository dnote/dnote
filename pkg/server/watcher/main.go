package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

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

func execCmd() *exec.Cmd {
	return command("go", []string{"run", "main.go", "start", "-port", "3000"}, "..")
}

func main() {
	w := watcher.New()
	w.IgnoreHiddenFiles(true)
	w.SetMaxEvents(1)

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
				e = execCmd()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(".."); err != nil {
		log.Fatalln(err)
	}

	e = execCmd()

	log.Printf("watching %d files", len(w.WatchedFiles()))
	if err := w.Start(time.Millisecond * 500); err != nil {
		log.Fatalln(err)
	}
}
