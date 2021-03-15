---
subcategory: "Virtual Private Cloud (VPC)"
---

# g42cloud\_networking\_secgroup

Manages a V2 neutron security group resource within G42Cloud.
Unlike Nova security groups, neutron separates the group from the rules
and also allows an admin to target a specific tenant_id.

## Example Usage

```hcl
resource "g42cloud_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    security group.

* `name` - (Required, String) A unique name for the security group.

* `description` - (Optional, String) Description of the security group.

* `tenant_id` - (Optional, String, ForceNew) The owner of the security group. Required if admin
    wants to create a port for another tenant. Changing this creates a new
    security group.

* `delete_default_rules` - (Optional, Bool, ForceNew) Whether or not to delete the default
    egress security rules. This is `false` by default. See the below note
    for more information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Default Security Group Rules

In most cases, G42Cloud will create some egress security group rules for each
new security group. These security group rules will not be managed by
Terraform, so if you prefer to have *all* aspects of your infrastructure
managed by Terraform, set `delete_default_rules` to `true` and then create
separate security group rules such as the following:

```hcl
resource "g42cloud_networking_secgroup_rule" "secgroup_rule_v4" {
  direction         = "egress"
  ethertype         = "IPv4"
  security_group_id = g42cloud_networking_secgroup.secgroup.id
}

resource "g42cloud_networking_secgroup_rule" "secgroup_rule_v6" {
  direction         = "egress"
  ethertype         = "IPv6"
  security_group_id = g42cloud_networking_secgroup.secgroup.id
}
```

Please note that this behavior may differ depending on the configuration of
the G42Cloud cloud. The above illustrates the current default Neutron
behavior. Some G42Cloud clouds might provide additional rules and some might
not provide any rules at all (in which case the `delete_default_rules` setting
is moot).


## Timeouts
This resource provides the following timeouts configuration options:
- `delete` - Default is 10 minute.
## Import

Security Groups can be imported using the `id`, e.g.

```
$ terraform import g42cloud_networking_secgroup.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
