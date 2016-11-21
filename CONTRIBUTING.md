CONTRIBUTING
============

## Prerequisites

These are the prequesents required for testing

* [Docker](https://docs.docker.com/engine/installation/)
* [Golang](https://golang.org/dl/)
* [Glide](https://github.com/Masterminds/glide) ||  `go get github.com/Masterminds/glide`

## Packaging

To build `challenge.tar.gz` run `make package`

## Testing

### Local

Local has limitations as far as distributed testing goes.Running unit
tests locally (`make test`) is preferred for your own velocity.

* Run unit tests with `make test`
* Run `make cover` to view coverage stats
  NOTE: use `make test cover` if you want to re-run tests before
  generate the coverfile
* Build the binary specifically for your system (i.e. OS X) `make bins`
* Build the binary specifically for **linux** with `make linux`
* Execute the binary `./bin/challenge-executable`

### Docker

Docker containers are building using the `Dockerfile`, which exposes
port `7777` and forwards that your local ports `123{4,5,6}`

* Deploy to local docker containers using `make deploy`
  Tear down the containers with `make stop` or `make remove`
* View debug-request tracing with `make open-debug-requests`
  NOTE: This may not work outside of OS X, so open `http://localhost:1234/debug/requests` in your browser for
  any any port
* Execute the integration script
  Load peer configs with `./tests/live-requests.sh config`
```bash
➜  challenge git:(master) ✗ ./tests/live-requests.sh config
UPDATING CONFIG
HTTP/1.1 200 OK
Date: Mon, 21 Nov 2016 16:42:55 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8

NODE: 1234
VALUE: 92
CONSISTENT VALUE: 108

NODE: 1235
VALUE: 108
CONSISTENT VALUE: 108

NODE: 1236
VALUE: 108
CONSISTENT VALUE: 108
```
  Simulate requests using the same script (without args) as needed `./tests/live-requests.sh`

```bash
➜  challenge git:(master) ✗ ./tests/live-requests.sh
NODE: 1234
VALUE: 223
CONSISTENT VALUE: 303

NODE: 1235
VALUE: 303
CONSISTENT VALUE: 303

NODE: 1236
VALUE: 303
CONSISTENT VALUE: 303
```
* Pause a container and rerun the simulation requests
  NOTE: the output may look weird on the paused container port

```bash
➜  challenge git:(master) ✗ ./tests/live-requests.sh 2> /dev/null
NODE: 1234
VALUE: 407
CONSISTENT VALUE: {"method": "GET", "error": "context deadline exceeded", "best-guess": "413"}

NODE: 1235
VALUE: 413
CONSISTENT VALUE: {"method": "GET", "error": "context deadline exceeded", "best-guess": "413"}

NODE: 1236
VALUE: CONSISTENT VALUE:
```