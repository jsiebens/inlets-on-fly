# inlets-on-fly

inlets-on-fly automates the task of creating an [inlets Pro](https://inlets.dev) HTTP exit-server (tunnel server) on the [fly.io](https://fly.io) platform.

This automation started as a bash script which you can find [here](https://gist.github.com/jsiebens/4cf66c135ecefe8638c06a16c488b201)

Read more in the blog post [Run an inlets Pro tunnel server for free on fly.io](https://inlets.dev/blog/2021/07/07/inlets-fly-tutorial.html)

## Pre-requisites

inlets-on-fly make HTTP requests to the Fly.io API, so make sure you grab an api token using the flyctl CLI tool.

- [Installing flyctl](https://fly.io/docs/hands-on/install-flyctl/)
- [Login To Fly](https://fly.io/docs/getting-started/log-in-to-fly/)
- [flyctl auth token](https://fly.io/docs/flyctl/auth-token/)

## Example

``` bash
$ export FLY_API_TOKEN=$(fly auth token)
$ inlets-on-fly create
Host: smooth-honeybee-3106/3d8dd15ce95589, status: created
[1/500] Host: smooth-honeybee-3106/3d8dd15ce95589, status: active
inlets Pro HTTPS (0.9.9) server summary:
  IP: smooth-honeybee-3106.fly.dev
  HTTPS Domains: .fly.dev
  Auth-token: sNAZK3MxnAPdcOfIGIsR2yIZpEB6qa5IYRlJ9G4kiB9D8AYydFSoiY2QFox9Hwym

Command:

# Obtain a license at https://inlets.dev/pricing
# Store it at $HOME/.inlets/LICENSE or use --help for more options

# Where to route traffic from the inlets server
export UPSTREAM="http://127.0.0.1:8000"

inlets-pro http client --url "wss://smooth-honeybee-3106.fly.dev:8123" \
--token "sNAZK3MxnAPdcOfIGIsR2yIZpEB6qa5IYRlJ9G4kiB9D8AYydFSoiY2QFox9Hwym" \
--upstream $UPSTREAM

To delete:
  inlets-on-fly delete --id "smooth-honeybee-3106/3d8dd15ce95589"
```
