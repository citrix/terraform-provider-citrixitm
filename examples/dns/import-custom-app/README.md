# Import Custom App

In this example we use `terraform import` to allow an existing custom DNS routing app to be managed by Terraform.

## Example Files

The example module includes the following files:

| Filename       | Description |
| ---            | --- |
| README.md      | The file you're reading now. |
| main.tf        | The main Terraform config file for this module. |
| vars.tf        | A file specifying input variables that must be set when you run Terraform commands against this module. |
| outputs.tf     | A file specifying output variables that will be set automatically when you run Terraform commands that update infrastructure. |
| app.js | A JavaScript file containing a sample custom Citrix ITM DNS routing application. |


## Running the Example

As you go through this example, you'll be asked to run various Terraform commands, and the instructions often remind you to execute the command from within the example module directory. This simply means that you should be in the directory containing the example's \*.tf files when executing the command. Terraform acts on files starting from the current working directory, which it considers to be the "root" module.

### Authenticating

Like all of the examples in this repository, you'll need a client ID and secret pair, which can be generated on the Cedexis Portal [OAuth configuration](https://portal.cedexis.com/ui/api/oauth) page.

For convenience we assume that you have two environment variables set to these.

Example:

```
CITRIXITM_CLIENT_ID=<your client ID>
CITRIXITM_CLIENT_SECRET=<your client secret>
```

These environment variables are used in place of the client ID and secret in the commands that follow.

## Setup the App Using Citrix ITM Portal

Before we can import the app, it must first be created using the Citrix ITM Portal UI. Go to the [Openmix applications page](https://portal.cedexis.com/ui/openmix/applications) now and create a new application. You can use the app.js file included in this example as the application code file.

Give the file an name that's easy to remember.

<img src="./res/new-app-to-import.png" width="400px" alt="New App Basic Settings">

Set the app fallback CNAME. For purposes of this example, you can set it to anything.

<img src="./res/choose-fallback-cname.png" width="400px" alt="Choose Fallback CNAME">

Choose the application file by clicking on the FOO button and navigating to the app.js file found in this example.

<img src="./res/choose-app-file.png" width="400px" alt="Choose app file">

Click Publish.

<img src="./res/click-publish.png" width="400px" alt="Publish the app">

Take note of the application ID.

<img src="./res/note-app-id.png" width="400px" alt="Note App ID">

### Import the DNS Routing App

The first step is to initialize the module. From within the example module directory, execute:

```bash
$ terraform init
```
