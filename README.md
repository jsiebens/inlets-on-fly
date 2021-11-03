# inlets-on-fly

inlets-on-fly automates the task of creating an [inlets-pro](https://inlets.dev) exit-server (tunnel server) on the [fly.io](https://fly.io) platform.

## prerequisites

inlets-on-fly is actually a little wrapper around flyctl, so make sure you have that CLI installed and that you are authenticated.
- [Installing flyctl](https://fly.io/docs/getting-started/installing-flyctl/)
- [Login To Fly](https://fly.io/docs/getting-started/login-to-fly/)

## example

``` bash
$ inlets-on-fly create --region ams --ports 8080:80,8080:443,5432:10032
Temp dir name: /tmp/inletsfly-885552054

Selected App Name: vast-goblin-6527


New app created: vast-goblin-6527
Region Pool: 
ams
Backup Region: 
fra
lhr
Secrets are staged for the first deployment
Deploying vast-goblin-6527
==> Validating app configuration
--> Validating app configuration done
Services
TCP 10023 ⇢ 8123
TCP 10032 ⇢ 5432
TCP 80/443 ⇢ 8080
Waiting for remote builder fly-builder-crimson-dust-9853...
==> Creating build context
--> Creating build context done
==> Building image with Docker
--> docker host: 20.10.8 linux x86_64
Sending build context to Docker daemon     422B
Step 1/2 : FROM ghcr.io/inlets/inlets-pro:0.9.1
 ---> 68840e710735
Step 2/2 : CMD ["tcp", "server", "--auto-tls-san=vast-goblin-6527.fly.dev", "--token-env=TOKEN"]
 ---> Running in c148498ae9b3
 ---> 39a98585fc9b
Successfully built 39a98585fc9b
Successfully tagged registry.fly.io/vast-goblin-6527:deployment-1636030837
--> Building image done
==> Pushing image to fly
The push refers to repository [registry.fly.io/vast-goblin-6527]
8345a2e5488b: Preparing
c0d270ab7e0d: Preparing
8345a2e5488b: Mounted from amused-stag-5259
c0d270ab7e0d: Mounted from amused-stag-5259
deployment-1636030837: digest: sha256:5f0d03afb7044731670ba0cd1e20fa5793b4f4d286997b6aea51cb6d5879545c size: 738
--> Pushing image done
Image: registry.fly.io/vast-goblin-6527:deployment-1636030837
Image size: 19 MB
==> Creating release
Release v2 created

You can detach the terminal anytime without stopping the deployment
Monitoring Deployment

v0 is being deployed
ad3148e3: ams running healthy [health checks: 1 total, 1 passing]
--> v0 deployed successfully
==================================================================
inlets PRO TCP (vast-goblin-6527) server summary:

  URL: wss://vast-goblin-6527.fly.dev:10023/connect
  Auth-token: XuGk0bSsLuf9q3Q2gDcydohXUyOwuwyl1WzU3ep4KkpLA9cjWh0MLpNEtdWEP9ra

Command:

# Obtain a license at https://inlets.dev
# Store it at $HOME/.inlets/LICENSE or use --help for more options

export PORTS="5432,8080,"
export UPSTREAM="localhost"

inlets-pro tcp client \
  --url wss://vast-goblin-6527.fly.dev:10023/connect \
  --token XuGk0bSsLuf9q3Q2gDcydohXUyOwuwyl1WzU3ep4KkpLA9cjWh0MLpNEtdWEP9ra \
  --upstream $UPSTREAM \
  --ports $PORTS

To delete:
  flyctl destroy vast-goblin-6527

```