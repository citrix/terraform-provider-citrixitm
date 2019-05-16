# Using Docker

Docker can be used to run and test the Citrix ITM provider within a containerized environment. The Dockerfile included in the project creates an image provisioned with everything needed for most development and testing purposes. GNUmakefile contains targets allowing you start and manage a long-running container, within which you can execute Terraform commands and run the test suite. All of these tasks are explained below. 

* [Build a Docker Image](#build-a-docker-image)
* [Create a Container](#create-a-container)
* [Restart a Stopped Container](#restart-a-stopped-container)
* [Attach to a Running Container](#attach-to-a-running-container)
* [Run Tests](#run-tests)
* [Ad hoc testing](#ad-hoc-testing)
* [Viewing Logs](#viewing-logs)
* [Starting Over](#starting-over)

## Build a Docker Image

A Docker image is like a template that is used to run containers. The specification for the image is found in the project's Dockerfile.

From within the project's root directory on the Docker host, execute the following Make target:

```bash
$ make docker-build
```

This prints a lot of information as it executes the steps defined in the Dockerfile. Once complete, there will be a new imaged tagged citrixitm-terraform:latest on the Docker host. You can see it using the following command:

```bash
$ docker images --filter=reference='citrixitm*'
REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
citrixitm-terraform   latest              4a93e5c833fd        5 seconds ago       580MB
```

This step generally needs to be repeated occasionally, such as when you delete the Docker image from your system for any reason, or when the Dockerfile is modified. The latter might occur when we update the version of Terraform included, for example.

## Create a Container

The purpose of a Docker image is to give Docker a known state from which to create containers. Within a running container, we'll be able to run Terraform and test the provider. The make docker-run target creates a running container with a Bash shell that you can use for a variety of purpose. There are some arguments to know about:

| Argument | Required | Description |
| -------- | -------- | ----------- |
| ITM_BASE_URL | YES | This tells the provider what Portal instance to access when making API requests. Examples include https://portal.cedexis.com/api and https://portalha.dev.cedexis.com/api. For development and testing purposes, you'll almost always use portalha.dev. See Citrix ITM Terraform Provider - Portal Instances. |
| ITM_CLIENT_ID | YES | The Oauth 2.0 client id for the credentials you created at https://portalha.dev.cedexis.com/ui/api/oauth |
| ITM_CLIENT_SECRET | YES | The Oauth 2.0 client secret for the credentials you created at https://portalha.dev.cedexis.com/ui/api/oauth |
| ITM_HOST_MODULE_DIR | NO | The path to a Terraform module on the Docker host that you would like to test within the container. If this argument is given, this directory is made available as a bind mount to /terraform-module within the container. You can omit this argument if you don't plan to run any Terraform commands to exercise the provider. For example, you may only want to run the unit and acceptance tests, which don't actually involve running Terraform and don't require an external Terraform module to operate on. |

Example usage:

```bash
$ make docker-run ITM_BASE_URL=https://portalha.dev.cedexis.com/api ITM_CLIENT_ID=<your client ID> ITM_CLIENT_SECRET=<your client secret> ITM_HOST_MODULE_DIR=<path to test module on Docker host>
```

If you'd rather not supply command line arguments, make docker-run can pick these up from the environment if you have them set. For example, this works also:

```bash
$ export ITM_BASE_URL=https://portalha.dev.cedexis.com/api
$ export ITM_CLIENT_ID=<your client ID>
$ export ITM_CLIENT_SECRET=<your client secret>
$ export ITM_HOST_MODULE_DIR=<path to test module on Docker host>
$ make docker-run
```

Once you execute make docker-run, you are left at a Bash prompt in the /terraform-provider-citrixitm directory:

```bash
bash-4.4# pwd
[container] /terraform-provider-citrixitm $
bash-4.4#
```

Here you can check to see that the environment variables are set within the container:

```bash
[container] /terraform-provider-citrixitm $ env | grep ITM
ITM_CLIENT_ID=<your client id>
ITM_CLIENT_SECRET=<your client secret>
ITM_BASE_URL=https://portalha.dev.cedexis.com/api
```

## Restart a Stopped Container

The container runs as long as the original Bash session that started, so it's important to understand which Make targets do this:

| Command | Description |
| ------- | ----------- |
| make docker-run | Creates a new container named citrixitm_tf_dev_container. The command fails if a container with this name already exists (see below). |
| make docker-start | Restarts a stopped container. |

Container names are unique, so Docker complains if you try to create a new container with a name already being used by an existing container. We use this to allow you to completely shutdown a container with artifacts from development and testing (e.g. Terraform state and log files) and resume later. This provides a lot of flexibility in how you can use the containerized environment.

Example

This is a common use case when you're first starting out (you've just built the Docker image for the first time and need to start your first container, for example), or after you've deleted the citrixitm_td_dev_container container for any reason (see Starting Over).

For the sake of this example, we'll assume that the ITM_* variables are exported to the environment as described above, and don't need to be supplied to make docker-run.

First create the container with make docker-run:

```bash
$ make docker-run
ITM variables...
ITM_BASE_URL: https://portalha.dev.cedexis.com/api
ITM_CLIENT_ID: <your client id>
ITM_CLIENT_SECRET: <your client secret>
ITM_HOST_MODULE_DIR:
docker run -it --name citrixitm_tf_dev_container --env ITM_BASE_URL --env ITM_CLIENT_ID --env ITM_CLIENT_SECRET --mount type=bind,readonly=1,src=/Users/jacob/Documents/repos/cedexis/terraform-provider-citrixitm,dst=/terraform-provider-citrixitm citrixitm-terraform /bin/bash
[container] /terraform-provider-citrixitm $
```

Leave something in the container we can look for later:

```bash
[container] /terraform-provider-citrixitm $ echo Foo > /tmp/foo
[container] /terraform-provider-citrixitm $ cat /tmp/foo
Foo
```

Then exit the Bash session within the container:

```bash
[container] /terraform-provider-citrixitm $ exit
exit
$
```

At this point, the container still exists, but is stopped, as you can see:

```bash
$ docker ps --all --filter=name='citrixitm_tf'
CONTAINER ID        IMAGE                 COMMAND             CREATED             STATUS                     PORTS               NAMES
5a974e670915        citrixitm-terraform   "/bin/bash"         9 minutes ago       Exited (0) 4 minutes ago                       citrixitm_tf_dev_container
```

Trying to execute make docker-run again produces an error message similar to this:

```bash
$ make docker-run
ITM variables...
ITM_BASE_URL: https://portalha.dev.cedexis.com/api
ITM_CLIENT_ID: <your client id>
ITM_CLIENT_SECRET: <your client secret>
ITM_HOST_MODULE_DIR:
docker: Error response from daemon: Conflict. The container name "/citrixitm_tf_dev_container" is already in use by container "5a974e670915ff5eb0ee67719dbddd497d73ae77108d41fddd18fc88795a7ed0". You have to remove (or rename) that container to be able to reuse that name.
See 'docker run --help'.
make: *** [docker-run] Error 125
```

Simply run make docker-start instead:

```bash
$ make docker-start
[container] /terraform-provider-citrixitm $
```

This has same effect in that it drops you into an interactive Bash session within the container, but it re-uses the stopped container, allowing you to resume from where you left off:

```bash
[container] /terraform-provider-citrixitm $ cat /tmp/foo
Foo
```

## Attach to a Running Container

The container is running as long as the original command executed within it. In our case, this is /bin/bash from either the make docker-run or make docker-start Make targets. For example:

```bash
$ make docker-start
[container] /terraform-provider-citrixitm $ ps aux
PID   USER     TIME  COMMAND
    1 root      0:00 /bin/bash
   11 root      0:00 ps aux
```

You may wish to run another Bash session using the same container. This is what the make docker-exec-bash target is for. It starts another bash session in the container that's already running. You can now see that there are (at least) two Bash session running:

```bash
$ make docker-exec-bash
[container] /terraform-provider-citrixitm $ ps aux
PID   USER     TIME  COMMAND
    1 root      0:00 /bin/bash
   12 root      0:00 /bin/bash
   26 root      0:00 ps aux
```

It's important to keep in mind that Docker still shuts down the container when the original command used to start it terminates, so if you exit the first Bash shell for any reason, it'll boot you out of any attached shells as well. It's a good practice to only use the original shell to execute Terraform commands or run unit tests, and use auxiliary attached shells to inspect artifacts as they are create, e.g. watching the Terraform log file, etc.

## Run Tests

From within /terraform-provider-citrixitm, you can run the unit and acceptance test suites:

```bash
[container] /terraform-provider-citrixitm $ make test
...
[container] /terraform-provider-citrixitm $ make testacc
...
```

Acceptance tests involve making real API requests using the ITM_* variables defined when the container was created initially, so take care to use the intended Portal instance and Oauth credentials.

Since the /terraform-provider-citrix is bind mounted to the project repo on the Docker host machine, you can make changes to the code on the host machine "locally" and immediately see those changes reflected by running tests within the container. 

## Ad hoc testing

You may also wish to use the container to exercise the Citrix ITM provider using Terraform, as an end-user would. The Dockerfile specifies instructions to download a recent copy of the Terraform executable when creating the Docker image, so any container created from it already has Terraform installed and ready to use.

To do this, you'll need to have specified the ITM_HOST_MODULE_DIR argument to make docker-run. If you didn't do that before, or you aren't yet sure what a Terraform module is, we'll go through a quick walkthrough now.

First of all, make docker-run creates a container named "citrixitm_tf_dev_container". As long as this container exists, we won't be able to execute make docker-run again with different arguments, so let's delete it from the Docker host:

```bash
$ docker rm citrixitm_tf_dev_container
```

Now we need to have a Terraform module to test. A Terraform module is simply a directory containing Terraform configuration files. The current working directory when you run Terraform is called the root module. Modules are a way to create parameterized configurations that are reusable, but for purposes of testing the Citrix ITM provider, a simple, solitary root module will be sufficient. We'll create one for demonstration purposes now. You can follow these instructions step by step to do this from the command line. We'll assume that the new module directory will be created at $(HOME)/Documents/terraform_test_module, but you can adapt the instructions to put it elsewhere if you choose.

```bash
$ cd ~/Documents
 
# Create the module directory and make it the current working directory
$ mkdir terraform_test_module && cd terraform_test_module
 
# Create main.tf
$ cat << 'EOF' > main.tf
terraform {
    required_version = ">= 0.11, < 0.12"
}
 
provider "citrixitm" {
 
}
 
resource "citrixitm_dns_app" "docker_test_app" {
    name            = "Testing Docker container"
    description     = "A sample DNS app for use in testing"
    app_data        = "${file("${path.module}/files/app.js")}"
    fallback_cname  = "origin.example.com"
}
EOF
 
# Create app.js (the Citrix ITM DNS app code)
$ mkdir files
$ cat << EOF > files/app.js
function init(config) {}
 
function onRequest(request, response) {
    response.addCName('foo.example.com');
    response.setTTL(20);
}
EOF
 
# Change back to the Citrix ITM repo to continue
$ cd <path to Citrix ITM repo>
```

Be sure that you changed back to the Citrix ITM provider repo in the last step.

Depending on whether you chose to execute make docker-run by passing it arguments, or by exporting ITM_* variables to the shell environment, re-run make docker-run as you did in the Run a Containerized Bash Session section, but supply the path to the terraform_test_module directory as the value for ITM_HOST_MODULE_DIR.

For example,

```bash
$ make docker-run ITM_BASE_URL=https://portalha.dev.cedexis.com/api ITM_CLIENT_ID=<your client ID> ITM_CLIENT_SECRET=<your client secret> ITM_HOST_MODULE_DIR=~/Documents/terraform_test_module
```

Or...

```bash
$ export ITM_BASE_URL=https://portalha.dev.cedexis.com/api
$ export ITM_CLIENT_ID=<your client ID>
$ export ITM_CLIENT_SECRET=<your client secret>
$ export ITM_HOST_MODULE_DIR=~/Documents/terraform_test_module
$ make docker-run
```

Since the ITM_HOST_MODULE_DIR variable is defined, the docker-run recipe additionally bind mounts the specified host directory to /terraform-module within the container.

You should now be back at the Bash prompt within the newly created container. You can now change to the /terraform-module directory and begin using Terraform to manage the resources defined in the demo root module.

```bash
[container] /terraform-provider-citrixitm $ cd /terraform-module/
```

Try running terraform init:

```bash
[container] /terraform-module $ terraform init
 
Initializing provider plugins...
- Checking for available provider plugins on https://releases.hashicorp.com...
 
Provider "citrixitm" not available for installation.
 
A provider named "citrixitm" could not be found in the official repository.
 
This may result from mistyping the provider name, or the given provider may
be a third-party provider that cannot be installed automatically.
 
In the latter case, the plugin must be installed manually by locating and
downloading a suitable distribution package and placing the plugin's executable
file in the following directory:
    terraform.d/plugins/linux_amd64
 
Terraform detects necessary plugins by inspecting the configuration and state.
To view the provider versions requested by each module, run
"terraform providers".
```

Oh no! What happened?

This occurred because we haven't actually built the provider binary yet. There are two ways to do this. You can do it within the container manually, or using a Make target executed on the Docker host.

To build the provider binary manually, execute make build from within the /terraform-citrixitm-provider directory inside the running container. Then copy the binary into the Terraform plugins directory mentioned in the earlier error message.

```bash
[container] /terraform-module $ cd /terraform-citrixitm-provider
[container] /terraform-citrixitm-provider $ make build
...
[container] /terraform-citrixitm-provider $ mkdir -p ~/.terraform.d/plugins/linux_amd64
[container] /terraform-citrixitm-provider $ cp /go/bin/terraform-provider-citrixitm ~/.terraform.d/plugins/linux_amd64/
```

Alternatively, there's a Make target that can be executed on the Docker host to perform the same steps described above.

For this to work, you should keep the container running in one shell on the Docker host and open another to execute make docker-exec-install-plugin. That's because it uses docker exec internally, which acts on a running container. 

```bash
$ make docker-exec-install-plugin
docker exec -it citrixitm_tf_dev_container scripts/docker_install_binary.sh
Building the Terraform provider plugin binary...
==> Checking that code complies with gofmt requirements...
go install
...
```

Using one of the methods described above, you should now have built the binary and moved it into the root user's Terraform plugins directory.

You can re-run terraform init and terraform providers to see that Terraform now recognizes the provider plugin:

```bash
[container] /terraform-provider-citrixitm $ cd /terraform-module/
[container] /terraform-module $ terraform init

Initializing provider plugins...
 
Terraform has been successfully initialized!
 
You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.
 
If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
[container] /terraform-module $ terraform providers
.
└── provider.citrixitm
```

Now you can do things like terraform plan, terraform apply and terraform destroy.

```bash
[container] /terraform-module $ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.
 
 
------------------------------------------------------------------------
 
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create
 
Terraform will perform the following actions:
 
  + citrixitm_dns_app.docker_test_app
      id:             <computed>
      app_data:       "function init(config) {}\n\nfunction onRequest(request, response) {\n    response.addCName('foo.example.com');\n    response.setTTL(20);\n}\n"
      cname:          <computed>
      description:    "A sample DNS app for use in testing"
      fallback_cname: "origin.example.com"
      fallback_ttl:   "20"
      name:           "Testing Docker container"
      version:        <computed>
 
 
Plan: 1 to add, 0 to change, 0 to destroy.
 
------------------------------------------------------------------------
 
Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.

[container] /terraform-module $ terraform apply
 
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create
 
Terraform will perform the following actions:
 
  + citrixitm_dns_app.docker_test_app
      id:             <computed>
      app_data:       "function init(config) {}\n\nfunction onRequest(request, response) {\n    response.addCName('foo.example.com');\n    response.setTTL(20);\n}\n"
      cname:          <computed>
      description:    "A sample DNS app for use in testing"
      fallback_cname: "origin.example.com"
      fallback_ttl:   "20"
      name:           "Testing Docker container"
      version:        <computed>
 
 
Plan: 1 to add, 0 to change, 0 to destroy.
 
Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.
 
  Enter a value: yes
 
citrixitm_dns_app.docker_test_app: Creating...
  app_data:       "" => "function init(config) {}\n\nfunction onRequest(request, response) {\n    response.addCName('foo.example.com');\n    response.setTTL(20);\n}\n"
  cname:          "" => "<computed>"
  description:    "" => "A sample DNS app for use in testing"
  fallback_cname: "" => "origin.example.com"
  fallback_ttl:   "" => "20"
  name:           "" => "Testing Docker container"
  version:        "" => "<computed>"
citrixitm_dns_app.docker_test_app: Creation complete after 2s (ID: 50)
 
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

[container] /terraform-module $ terraform destroy
citrixitm_dns_app.docker_test_app: Refreshing state... (ID: 50)
 
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy
 
Terraform will perform the following actions:
 
  - citrixitm_dns_app.docker_test_app
 
 
Plan: 0 to add, 0 to change, 1 to destroy.
 
Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.
 
  Enter a value: yes
 
citrixitm_dns_app.docker_test_app: Destroying... (ID: 50)
citrixitm_dns_app.docker_test_app: Destruction complete after 0s
 
Destroy complete! Resources: 1 destroyed.
```

## Viewing Logs

Within the container, the TF_LOG and TF_LOG_PATH environment variables are set at the image level:

```bash
[container] /terraform-module $ env | grep TF_LOG
TF_LOG_PATH=/var/log/terraform.log
TF_LOG=TRACE
```

These variables influence how much information Terraform logs and where. You can see that it's set to write to a file at /var/log/terraform.log. Inspect this file to see information recorded when you run tests or when you execute Terraform commands directly.

## Starting Over

The method described above calls for a long-running container that persists after its Bash session terminates. This is to support restarting the container without losing various artifacts, such as downloaded Go module dependencies and Terraform log files. But sometimes you'd like to start over with a fresh container. In practice, this is mainly when you've done something to change the underlying Docker image, such as when the project's Dockerfile changes for any reason, or when you're moving between Git branches.

To start over with a black slate, delete any existing container named "citrixitm_tf_dev_container":

```bash
$ docker rm citrixitm_tf_dev_container
citrixitm_tf_dev_container
```
Then create a new container, as described in Create a Container. 
