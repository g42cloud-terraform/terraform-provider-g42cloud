---
subcategory: "Identity and Access Management (IAM)"
---

# g42cloud\_identity\_user\_v3

Manages a User resource within G42Cloud IAM service.

Note: You _must_ have admin privileges in your G42Cloud cloud to use
this resource.

## Example Usage

```hcl
resource "g42cloud_identity_user_v3" "user_1" {
  name        = "user_1"
  description = "A user"
  password    = "password123!"

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the user. The user name consists of 5 to 32
     characters. It can contain only uppercase letters, lowercase letters, 
     digits, spaces, and special characters (-_) and cannot start with a digit.

* `description` - (Optional, String) A description of the user.

* `default_project_id` - (Optional, String) The default project this user belongs to.

* `domain_id` - (Optional, String) The domain this user belongs to.

* `enabled` - (Optional, Bool) Whether the user is enabled or disabled. Valid
    values are `true` and `false`.

* `password` - (Optional, String) The password for the user. It must contain at least 
     two of the following character types: uppercase letters, lowercase letters, 
     digits, and special characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Users can be imported using the `id`, e.g.

```
$ terraform import g42cloud_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
