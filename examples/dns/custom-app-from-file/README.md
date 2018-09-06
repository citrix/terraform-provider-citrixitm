# Custom App from File

This example Terraform module demonstrates how to maintain a custom DNS app using a local file as the source.

## Example Files

The example module includes the following files:

| Filename       | Description |
| ---            | --- |
| README.md      | The file you're reading now. |
| main.tf        | The main Terraform config file for this module. |
| vars.tf        | A file specifying input variables that must be set when you run Terraform commands against this module. |
| outputs.tf     | A file specifying output variables that will be set automatically when you run Terraform commands that update infrastructure. |
| website.dns.js | A JavaScript file containing a sample custom Citrix ITM DNS routing application. |

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

### Create a DNS Routing App

The first step is to initialize the module. From within the example module directory, execute:

```bash
$ terraform init
```

Then run the Terraform `plan` command. This causes Terraform to determine what actions to perform to make the infrastructure consistent with the configuration. Note that we utilize the environment variables storing your client ID and secret these commands.

```bash
$ terraform plan -var itm_client_id=$CLIENT_ID -var itm_client_secret=$CLIENT_SECRET
```

This should produce output similar to the following:

```bash
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + citrixitm_dns_app.website
      id:             <computed>
      app_data:       "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    // ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n"
      cname:          <computed>
      description:    "DNS routing for the website"
      fallback_cname: "origin.example.com"
      fallback_ttl:   "20"
      name:           "Simple Round Robin"
      version:        <computed>


Plan: 1 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

If the plan looks reasonable, go ahead and run it with the Terraform `apply` command:

```bash
$ terraform apply -var itm_client_id=$CLIENT_ID -var itm_client_secret=$CLIENT_SECRET
```

You'll be prompted to type "yes" if you think the proposed plan looks reasonable. If it does, then do so and press Enter. This should produce output similar to the following:

```bash
An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + citrixitm_dns_app.website
      id:             <computed>
      app_data:       "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    // ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n"
      cname:          <computed>
      description:    "DNS routing for the website"
      fallback_cname: "origin.example.com"
      fallback_ttl:   "20"
      name:           "Simple Round Robin"
      version:        <computed>


Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

citrixitm_dns_app.website: Creating...
  app_data:       "" => "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    // ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n"
  cname:          "" => "<computed>"
  description:    "" => "DNS routing for the website"
  fallback_cname: "" => "origin.example.com"
  fallback_ttl:   "" => "20"
  name:           "" => "Simple Round Robin"
  version:        "" => "<computed>"
citrixitm_dns_app.website: Creation complete after 3s (ID: 23)

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

dns_app_cname = 2-01-d896-0017.cdx.cedexis.net
dns_app_id = 23
dns_app_version = 1
```

Note the output variables displayed at the end. These are the output variables that we specified in the module's outputs.tf file. Of course you will likely see different values than those shown above.

You can see the output variables again by running the Terraform `output` command from within the module directory:

```bash
$ terraform output
dns_app_cname = 2-01-d896-0017.cdx.cedexis.net
dns_app_id = 23
dns_app_version = 1
```

Now if you visit the Portal [Openmix Applications](https://portal.cedexis.com/ui/openmix/applications) page, you'll see a new application with the ID given by the `dns_app_id` output variable. You can view the application code in the Portal and confirm that it is identical to that of the website.dns.js file in this module.

You can test the DNS application itself at this point. For example, using the value of the `dns_app_cname` output variable:

```bash
$ dig +short <dns_app_cname>
foo.example.com.
```

This simple app returns a CNAME made using tokens listed in website.dns.js. For example, the first time the app is called it should return foo.example.com. The next time it should return bar.example.com, and so on. When it reaches the end of the tokens in the list, it starts over at the beginning.

One word of caution. When you call the app repatedly, say by repeating the `dig` command above over and over, that doesn't mean that the app itself is called each time. This is due to caching and it is influenced by the TTL (time to live) value set by the app. We set the response TTL for this app to 20 seconds, which means that you should generally expect to receive the same response for about 20 seconds.

To see how we can update an existing Citrix ITM app using Terraform, open the website.dns.js file in an editor and remove the `//` from the line starting with `// ,'quux'` and save the file.

Now when you execute Terraform `plan`, you'll see that Terraform notices the change and proposes a plan to update your infrastructure:

```bash
$ terraform plan -var itm_client_id=$(CLIENT_ID) -var itm_client_secret=$(CLIENT_SECRET)

Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

citrixitm_dns_app.website: Refreshing state... (ID: 23)

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  ~ citrixitm_dns_app.website
      app_data: "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    // ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n" => "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n"


Plan: 0 to add, 1 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

You can run the Terraform `apply` command as before, confirming the change when asked by entering "yes" and pressing Enter.

```bash
$ terraform apply -var itm_client_id=$(CLIENT_ID) -var itm_client_secret=$(CLIENT_SECRET)

citrixitm_dns_app.website: Refreshing state... (ID: 23)

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  ~ citrixitm_dns_app.website
      app_data: "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    // ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n" => "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n"


Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

citrixitm_dns_app.website: Modifying... (ID: 23)
  app_data: "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    // ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n" => "var prefixes = [\n    'foo'\n    ,'bar'\n    ,'baz'\n    ,'qux'\n    ,'quux'\n];\n\nvar lastIndex = 0;\n\nfunction init(config) {}\n\nfunction onRequest(request, response) {\n    if (lastIndex >= prefixes.length) {\n        lastIndex = 0;\n    }\n    var prefix = prefixes[lastIndex++];\n    response.addCName(prefix + '.example.com');\n    response.setTTL(20);\n}\n"
citrixitm_dns_app.website: Modifications complete after 1s (ID: 23)

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

dns_app_cname = 2-01-d896-0017.cdx.cedexis.net
dns_app_id = 23
dns_app_version = 2
```

Note that the value of the `dns_app_version` output variable has been incremented.

You can view the app code in the Portal and see that it has changed as expected.

You can run `dig +short <your dns_app_cname>` repeatedly until you obtain quux.example.com, which is the expected behavior.

Feel free to experiment at this point, changing the app code and running the `terraform apply ...` command again to update your Citrix ITM infrastructure. When you're finished, you _may_ want to clean up by running the Terraform `destroy` command:

```bash
$ terraform destroy -var itm_client_id=$(CLIENT_ID) -var itm_client_secret=$(CLIENT_SECRET)

citrixitm_dns_app.website: Refreshing state... (ID: 23)

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  - citrixitm_dns_app.website


Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

citrixitm_dns_app.website: Destroying... (ID: 23)
citrixitm_dns_app.website: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```

However you may not want to get in the habit of running the `destroy` command on a regular basis. That's because one of the goals of infrastucture-as-code is that your declarative configuration always represents your actual infrastructure as closely as possible. Terraform's `destroy` command helps clean up infrastructure created during development and testing of configuration changes, but in production it's better to remove the configuration for a resource that you wish gone, and then `apply` the change.
