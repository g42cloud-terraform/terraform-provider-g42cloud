# G42Cloud Provider

The G42Cloud provider is used to interact with the
many resources supported by G42Cloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the G42Cloud Provider
provider "g42cloud" {
  region     = "ae-ad-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}

# Create a VPC
resource "g42cloud_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

## Authentication

The G42 Cloud provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials ###

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "g42cloud" {
  region     = "ae-ad-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via the `G42_ACCESS_KEY` and
`G42_SECRET_KEY`, environment variables, representing your G42
Cloud Username and Password, respectively.

```hcl
provider "g42cloud" {}
```

Usage:

```sh
$ export G42_ACCESS_KEY="user-name"
$ export G42_SECRET_KEY="password"
$ export G42_REGION_NAME="ae-ad-1"
$ terraform plan
```


## Configuration Reference

The following arguments are supported:

* `region` - (Required) This is the G42 Cloud region. It must be provided,
  but it can also be sourced from the `G42_REGION_NAME` environment variables.

* `account_name` - (Optional, Required for IAM resources) The
  of IAM to scope to. If omitted, the `G42_ACCOUNT_NAME` environment variable is used.

* `access_key` - (Optional) The access key of the G42Cloud to use.
  If omitted, the `G42_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the G42Cloud to use.
  If omitted, the `G42_SECRET_KEY` environment variable is used.

* `project_name` - (Optional) The Name of the Project to login with.
  If omitted, the `G42_PROJECT_NAME` environment variable are used.

* `auth_url` - (Optional) The Identity authentication URL. If omitted, the
  `G42_AUTH_URL` environment variable is used.

* `security_token` - (Optional) The security token to authenticate with a
  [temporary security credential](https://docs.g42cloud.com/usermanual/obs/obs_03_0208.html).
  If omitted, the `G42_SECURITY_TOKEN` environment variable is used.

* `cloud` - (Optional) The endpoint of the cloud provider. If omitted, the
  `G42_CLOUD` environment variable is used. Defaults to `g42cloud.com`.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `G42_INSECURE` environment variable is used.

* `max_retries` - (Optional) This is the maximum number of times an API
  call is retried, in the case where requests are being throttled or
  experiencing transient failures. The delay between the subsequent API
  calls increases exponentially. The default value is `5`.
  If omitted, the `G42_MAX_RETRIES` environment variable is used.

* `enterprise_project_id` - (Optional) Default Enterprise Project ID for supported resources.
  If omitted, the `G42_ENTERPRISE_PROJECT_ID` environment variable is used.

* `endpoints` - (Optional) Configuration block in key/value pairs for customizing service endpoints.
  The following endpoints support to be customized: autoscaling, ecs, vpc, evs, iam.
  An example provider configuration:

```hcl
provider "g42cloud" {
  ...
  endpoints = {
    ecs = "https://ecs-customizing-endpoint.com"
  }
}
```

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `G42_REGION_NAME` - The region in which to create resources.

* `G42_ACCESS_KEY` - The username to login with.

* `G42_SECRET_KEY` - The password to login with.

* `G42_ACCOUNT_NAME` - The IAM account name.


You should be able to use any G42Cloud environment to develop on as long as the
above environment variables are set.
