---
subcategory: "Identity and Access Management (IAM)"
---

# g42cloud_identity_group_membership

Manages a User Group Membership resource within G42Cloud IAM service.

Note: You *must* have admin privileges in your G42Cloud cloud to use this resource.

## Example Usage

```hcl
resource "g42cloud_identity_group" "group_1" {
  name        = "group1"
  description = "This is a test group"
}

resource "g42cloud_identity_user" "user_1" {
  name     = "user1"
  enabled  = true
  password = "password12345!"
}

resource "g42cloud_identity_user" "user_2" {
  name     = "user2"
  enabled  = true
  password = "password12345!"
}

resource "g42cloud_identity_group_membership" "membership_1" {
  group = g42cloud_identity_group.group_1.id
  users = [
    g42cloud_identity_user.user_1.id,
    g42cloud_identity_user.user_2.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String, ForceNew) The group ID of this membership.

* `users` - (Required, List) A List of user IDs to associate to the group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
