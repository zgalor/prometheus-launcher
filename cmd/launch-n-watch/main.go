package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var watchPath = flag.String("watch-path", "", "the configuration path to watch")

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage: %s --watch-path=/path/to/watch app-to-launch [app-arg1 app-arg2 ...]\n", os.Args[0])
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(*watchPath) < 1 {
		log.Println("Missing path to watch")
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		log.Println("Missing app to launch")
		flag.Usage()
		os.Exit(1)
	}

	pid := launchApp()

	// define watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create != 0 {
					log.Println("modified file:", event.Name, "Sending SIGHUP to app: ", flag.Arg(0))
					// when watch notifies change send SIGHUP to app
					if e := syscall.Kill(pid, syscall.SIGHUP); e != nil {
						log.Println("error:", e)
					}
				}
			case watchErr := <-watcher.Errors:
				log.Println("error:", watchErr)
			}
		}
	}()

	// start watch on path
	err = watcher.Add(*watchPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	<-done
}

// launchApp launches the application.
// it forks it as a child process and redirect stdin and stdout to/from it in a separate go routine
// the func returns the child process pid
func launchApp() int {
	// launch app with args - save pid
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var stdoutBuf, stderrBuf bytes.Buffer
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	pid := cmd.Process.Pid
	log.Println("Launched", flag.Arg(0), "pid", pid, "with params:", flag.Args()[1:])

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Println("Error waiting for Cmd", err)
			os.Exit(1)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}
		outStr, errStr := stdoutBuf.String(), stderrBuf.String()
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	}()

	return pid
}
