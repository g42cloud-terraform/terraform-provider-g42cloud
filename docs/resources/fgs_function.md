---
subcategory: "FunctionGraph"
---

# g42cloud\_fgs\_function

Manages a Function resource within G42Cloud.

## Example Usage

### With base64 func code

```hcl
resource "g42cloud_fgs_function" "f_1" {
  name        = "func_1"
  app         = "default"
  agency      = "test"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
```

### With text code

```hcl
resource "g42cloud_fgs_function" "f_1" {
  name        = "func_1"
  app         = "default"
  agency      = "test"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler (event, context):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": json.dumps(event),
        "headers": {
            "Content-Type": "application/json"
        }
    }
EOF
}

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the Function resource.
  If omitted, the provider-level region will be used. Changing this creates a new Function resource.

* `name` - (Required, String, ForceNew) Specifies the name of the function.

* `app` - (Required, String) Specifies the group to which the function belongs.

* `handler` - (Required, String) Specifies the entry point of the function.

* `memory_size` - (Required, Int) Specifies the memory size(MB) allocated to the function.

* `runtime` - (Required, String) Specifies the environment for executing the function.

* `timeout` - (Required, Int) Specifies the timeout interval of the function, ranges from 3s to 900s.

* `code_type` - (Required, String) Specifies the function code type, which can be inline: inline code, zip: ZIP file,
	jar: JAR file or java functions, obs: function code stored in an OBS bucket.

* `func_code` - (Optional, String) Specifies the function code. When code_type is set to inline, zip, or jar, this parameter is mandatory,
	and the code can be encoded using Base64 or just with the text code.

* `code_url` - (Optional, String) Specifies the code url. This parameter is mandatory when code_type is set to obs.

* `code_filename` - (Optional, String) Specifies the name of a function file, This field is mandatory only when coe_type is
	set to jar or zip.

* `depend_list` - (Optional, String) Specifies the dependencies of the function.

* `user_data` - (Optional, String) Specifies the Key/Value information defined for the function. Key/value data might be parsed with [Terraform `jsonencode()` function]('https://www.terraform.io/docs/language/functions/jsonencode.html').

* `agency` - (Optional, String) Specifies the agency. This parameter is mandatory if the function needs to access other cloud services.

* `app_agency` - (Optional, String) Specifies An execution agency enables you to obtain a token or an AK/SK for accessing other cloud services.

* `description` - (Optional, String) Specifies the description of the function.

* `initializer_handler` - (Optional, String) Specifies the initializer of the function.

* `initializer_timeout` - (Optional, Int) Specifies the maximum duration the function can be initialized. Value range: 1s to 300s.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the function.
  Changing this creates a new function.

* `vpc_id`  - (Optional, String) Specifies the ID of VPC.

* `network_id`  - (Optional, String) Specifies the ID of subnet.

* `mount_user_id` - (Optional, String) Specifies the user ID, a non-0 integer from –1 to 65534. Default to -1.

* `mount_user_group_id` - (Optional, String) Specifies the user group ID, a non-0 integer from –1 to 65534. Default to -1.

* `func_mounts` - (Optional, List) Specifies the file system list. The `func_mounts` object structure is documented below.

The `func_mounts` block supports:

* `mount_type` - (Required, String) Specifies the mount type. Options: sfs, sfsTurbo, and ecs.

* `mount_resource` - (Required, String) Specifies the ID of the mounted resource (corresponding cloud service).

* `mount_share_path` - (Required, String) Specifies the remote mount path. Example: 192.168.0.12:/data.

* `local_mount_path` - (Required, String) Specifies the function access path.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `func_mounts/status` - The status of file system.
* `urn` - Uniform Resource Name
* `version` - The version of the function

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

Functions can be imported using the `id`, e.g.

```
$ terraform import g42cloud_fgs_function.my-func 7117d38e-4c8f-4624-a505-bd96b97d024c
```
