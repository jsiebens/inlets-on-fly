# inlets-on-fly

inlets-on-fly automates the task of creating an [inlets Pro](https://inlets.dev) exit-server (tunnel server) on the [fly.io](https://fly.io) platform.

This automation started as a bash script which you can find [here](https://gist.github.com/jsiebens/4cf66c135ecefe8638c06a16c488b201)

Read more in the blog post [Run an inlets Pro tunnel server for free on fly.io](https://inlets.dev/blog/2021/07/07/inlets-fly-tutorial.html)

## Pre-requisites

inlets-on-fly is actually a little wrapper around flyctl, so make sure you have that CLI installed and that you are authenticated.
- [Installing flyctl](https://fly.io/docs/hands-on/install-flyctl/)
- [Login To Fly](https://fly.io/docs/getting-started/log-in-to-fly/)

## Example

``` bash
$ inlets-on-fly create tcp --region ams --ports 5432,6379
Temp dir name: /tmp/inletsfly-486952897

Selected App Name: enjoyed-goldfish-6269


New app created: enjoyed-goldfish-6269
Region Pool: 
ams
Backup Region: 
fra
lhr
Secrets are staged for the first deployment
Deploying enjoyed-goldfish-6269
==> Validating app configuration
--> Validating app configuration done
Services
TCP 8123 ⇢ 8123
TCP 5432 ⇢ 5432
TCP 6379 ⇢ 6379
==> Creating build context
--> Creating build context done
==> Building image with Docker
--> docker host: 20.10.7 linux x86_64
Sending build context to Docker daemon  3.072kB
Step 1/2 : FROM ghcr.io/inlets/inlets-pro:0.9.1
 ---> 68840e710735
Step 2/2 : CMD ["tcp", "server", "--auto-tls-san=enjoyed-goldfish-6269.fly.dev", "--token-env=TOKEN"]
 ---> Running in 16d95f141336
 ---> a1b634947187
Successfully built a1b634947187
Successfully tagged registry.fly.io/enjoyed-goldfish-6269:deployment-1641658112
--> Building image done
==> Pushing image to fly
The push refers to repository [registry.fly.io/enjoyed-goldfish-6269]
8345a2e5488b: Preparing
c0d270ab7e0d: Preparing
c0d270ab7e0d: Mounted from grand-tortoise-6149
8345a2e5488b: Mounted from grand-tortoise-6149
deployment-1641658112: digest: sha256:9cfddad45c6b112714e8c607219ed1907e14d4a5dc5611d6dc4d07f402bf9507 size: 738
--> Pushing image done
Image: registry.fly.io/enjoyed-goldfish-6269:deployment-1641658112
Image size: 19 MB
==> Creating release
Release v2 created

You can detach the terminal anytime without stopping the deployment
Monitoring Deployment

v0 is being deployed
c7a5572c: ams pending
c7a5572c: ams running unhealthy [health checks: 1 total]
c7a5572c: ams running healthy [health checks: 1 total, 1 passing]
--> v0 deployed successfully
==================================================================
inlets PRO TCP (enjoyed-goldfish-6269) server summary:

  URL: wss://enjoyed-goldfish-6269.fly.dev:8123/connect
  Auth-token: aDdfDfl5uu4pmYzTxG066uLcosnkyyF27OQFCryZTWaAA2r8qCTbob8CErCnJOYm

Command:

# Obtain a license at https://inlets.dev
# Store it at $HOME/.inlets/LICENSE or use --help for more options


export PORTS="5432,6379,"
export UPSTREAM="localhost"

inlets-pro tcp client \
  --url wss://enjoyed-goldfish-6269.fly.dev:8123 \
  --token aDdfDfl5uu4pmYzExG066uLcosnkyRF27OQFCryZTWaAA2r8qCTbob8CErCnJOYm \
  --upstream $UPSTREAM \
  --ports $PORTS


To delete:
  flyctl destroy enjoyed-goldfish-6269
```
