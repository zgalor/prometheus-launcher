# prometheus-launcher
Launch Prometheus as child process, watch configmap path for changes, and send SIGHUP for reload

## Initial Dependencies Setup

`make setup`

## Build

`make build`

## Run

`bin/lnw [--watch_path=/path/to/watch] app-to-launch [app_arg1, ..., app_argn]`

### Arguments

`--watch-path` 

The path to watch for modification. 
File creation in that path will trigger a SIGHUP to the launched app
This can also be set using the `LNW_WATCH_PATH` environment variable. using `--watch-path` flag will override the environment variable.

* This argument is mandatory - must be set by env or flag

`app-to-launch` (required)

The name and location of the app to launch as a child process.
App output will be piped to stdout

`app_arg1,..,app_argn` (optional)

Application arguments to be passed to app


## Build prometheus-launcher docker image

`make container`

### Docker Image Notes
* prometheus arguments can be overriden externally.
* watch path defaults to `/etc/prometheus` but can be overridden using `LNW_WATCH_PATH`


