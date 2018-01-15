package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var watchPath = flag.String("watch-path", os.Getenv("LNW_WATCH_PATH"), "the configuration path to watch")

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage: %s [Flags] app-to-launch [app-arg1 app-arg2 ...]\n", os.Args[0])
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(*watchPath) < 0 {
		log.Println("Missing path to watch. Must be defined by either LNW_WATCH_PATH environment variable or --watch-path flag")
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
				if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
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
		log.Fatal("Failed to watch", *watchPath, err)
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
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	pid := cmd.Process.Pid
	log.Println("Launched", flag.Arg(0), "pid", pid, "with params:", flag.Args()[1:])

	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Println("Error waiting for Cmd", err)
			os.Exit(1)
		}
	}()

	return pid
}
