---
subcategory: "Virtual Private Cloud (VPC)"
---

# g42cloud_networking_vip

Manages a network VIP resource within G42Cloud VPC.

## Example Usage

```hcl
resource "g42cloud_vpc_subnet" "test" {
  ...
}

resource "g42cloud_networking_vip" "test" {
  network_id = g42cloud_vpc_subnet.test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VIP.
  If omitted, the provider-level region will be used. Changing this will create a new VIP resource.

* `network_id` - (Required, String, ForceNew) Specifies the network ID of the VPC subnet to which the VIP belongs.
  Changing this will create a new VIP resource.

* `ip_version` - (Optional, Int, ForceNew) Specifies the IP version, either `4` (default) or `6`.
  Changing this will create a new VIP resource.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address desired in the subnet for this VIP.
  Changing this will create a new VIP resource.

* `name` - (Optional, String) Specifies a unique name for the VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VIP ID.

* `mac_address` - The MAC address of the VIP.

* `status` - The VIP status.

* `device_owner` - The device owner of the VIP.

* `subnet_id` - The subnet ID in which to allocate IP address for this VIP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 2 minute.
* `delete` - Default is 2 minute.

## Import

Network VIPs can be imported using their `id`, e.g.:

```
$ terraform import g42cloud_networking_vip.test ce595799-da26-4015-8db5-7733c6db292e
```
