---
subcategory: "NAT Gateway (NAT)"
---

# g42cloud_nat_snat_rule

Manages a SNAT rule resource within G42Cloud.

## Example Usage

### SNAT rule in VPC scenario

```hcl
resource "g42cloud_nat_snat_rule" "snat_1" {
  nat_gateway_id = var.natgw_id
  floating_ip_id = var.publicip_id
  subnet_id      = var.subent_id
}
```

### SNAT rule in Direct Connect scenario

```hcl
resource "g42cloud_nat_snat_rule" "snat_2" {
  nat_gateway_id = var.natgw_id
  floating_ip_id = var.publicip_id
  source_type    = 1
  cidr           = "192.168.10.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the snat rule resource. If omitted, the
  provider-level region will be used. Changing this creates a new snat rule resource.

* `nat_gateway_id` - (Required, String, ForceNew) ID of the nat gateway this snat rule belongs to. Changing this creates
  a new snat rule.

* `floating_ip_id` - (Required, String) Specifies the EIP ID this snat rule connects to.
  Multiple EIPs are separated using commas (,). The number of EIP IDs cannot exceed 20.

* `subnet_id` - (Optional, String, ForceNew) ID of the subnet this snat rule connects to. This parameter and `cidr` are
  alternative. Changing this creates a new snat rule.

* `cidr` - (Optional, String, ForceNew) Specifies CIDR, which can be in the format of a network segment or a host IP
  address. This parameter and `subnet_id` are alternative. Changing this creates a new snat rule.

* `source_type` - (Optional, Int, ForceNew) Specifies the scenario. The valid value is 0 (VPC scenario) and 1 (Direct
  Connect scenario). Defaults to 0, only `cidr` can be specified over a Direct Connect connection. Changing this creates
  a new snat rule.

* `description` - (Optional, String) Specifies the description of the snat rule.
  The value is a string of no more than 255 characters, and angle brackets (<>) are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `floating_ip_address` - The actual floating IP address.
* `status` - The status of the snat rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

SNAT rules can be imported using the following format:

```
$ terraform import g42cloud_nat_snat_rule.snat_1 9e0713cb-0a2f-484e-8c7d-daecbb61dbe4
```
