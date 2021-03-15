---
subcategory: "Identity and Access Management (IAM)"
---

# g42cloud\_identity\_role\_assignment_v3

Manages a V3 Role assignment within group on G42Cloud IAM Service.

Note: You _must_ have admin privileges in your G42Cloud cloud to use
this resource. 

## Example Usage: Assign Role On Project Level

```hcl
data "g42cloud_identity_role_v3" "role_1" {
  name = "system_all_4" #ECS admin
}

resource "g42cloud_identity_group_v3" "group_1" {
  name = "group_1"
}


resource "g42cloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id   = g42cloud_identity_group_v3.group_1.id
  project_id = var.project_id
  role_id    = data.g42cloud_identity_role_v3.role_1.id
}
```

## Example Usage: Assign Role On Domain Level

```hcl

data "g42cloud_identity_role_v3" "role_1" {
  name = "secu_admin" #security admin
}

resource "g42cloud_identity_group_v3" "group_1" {
  name = "group_1"
}

resource "g42cloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id  = g42cloud_identity_group_v3.group_1.id
  domain_id = var.domain_id
  role_id   = data.g42cloud_identity_role_v3.role_1.id
}

```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required, String, ForceNew) Specifies the role to assign.

* `group_id` - (Required, String, ForceNew) Specifies the group to assign the role to.

* `domain_id` - (Optional, String, ForceNew; Required if `project_id` is empty) Specifies the domain to assign the role in.

* `project_id` - (Optional, String, ForceNew; Required if `domain_id` is empty) Specifies the project to assign the role in.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
