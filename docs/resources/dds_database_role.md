---
subcategory: "Document Database Service (DDS)"
---

# g42cloud_dds_database_role

Manages a database role resource within G42Cloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "role_name" {}
variable "db_name" {}
variable "owned_role_name" {}
variable "owned_role_db_name" {}

resource "g42cloud_dds_database_role" "test" {
  instance_id = var.instance_id

  name    = var.role_name
  db_name = var.db_name

  roles {
    name    = var.owned_role_name
    db_name = var.owned_role_db_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DDS instance is located.
  Changing this parameter will create a new role.

* `instance_id` - (Required, String, ForceNew) Specifies the DDS instance ID to which the role belongs.
  Changing this parameter will create a new role.

* `name` - (Required, String, ForceNew) Specifies the role name.
  The name can contain `1` to `64` characters, only letters, digits, underscores (_), hyphens (-) and dots (.) are
  allowed. Changing this parameter will create a new role.

* `db_name` - (Required, String, ForceNew) Specifies the database name to which the role belongs.
  The name can contain `1` to `64` characters, only letters, digits and underscores (_) are allowed.
  Changing this parameter will create a new role.

  -> After a DDS instances is created, the default database is **admin**.

* `roles` - (Optional, List, ForceNew) Specifies the list of roles owned by the current role.
  The [roles](#dds_roles) object structure is documented below.
  Changing this parameter will create a new role.

<a name="dds_roles"></a>
The `roles` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of role owned by the current role.
  The name can contain `1` to `64` characters, only letters, digits, underscores (_), hyphens (-) and dots (.) are
  allowed. Changing this parameter will create a new role.

* `db_name` - (Required, String, ForceNew) Specifies the database name to which this owned role belongs.
  Changing this parameter will create a new role.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `privileges` - The list of database privileges owned by the current role.
  The [privileges](#dds_privileges) object structure is documented below.

* `inherited_privileges` - The list of database privileges owned by the current role, includes all privileges
  inherited by owned roles. The [inherited_privileges](#dds_privileges) object structure is documented below.

<a name="dds_privileges"></a>
The `privileges` and `inherited_privileges` block supports:

* `resources` - The details of the resource to which the privilege belongs.
  The [resources](#dds_resources) structure is documented below.

* `actions` - The operation permission list.

<a name="dds_resources"></a>
The `resources` block supports:

* `collection` - The database collection type.

* `db_name` - The database name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 2 minutes.
* `delete` - Default is 2 minutes.

## Import

Database roles can be imported using their `id` (combination of `instance_id`, `db_name` and `name`), separated by a
slash (/), e.g.

```shell
terraform import g42cloud_dds_database_role.test &ltinstance_id&gt/&ltdb_name&gt/&ltname&gt
```
