---
subcategory: "Virtual Private Cloud (VPC)"
---

# g42cloud_networking_vip_associate

Manages a Vip associate resource within G42Cloud.

## Example Usage

```hcl
data "g42cloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

resource "g42cloud_networking_vip" "myvip" {
  network_id = data.g42cloud_vpc_subnet.mynet.id
}

resource "g42cloud_networking_vip_associate" "vip_associated" {
  vip_id   = g42cloud_networking_vip.myvip.id
  port_ids = [
    var.port_1,
    var.port_2
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, ForceNew) The region in which to create the vip associate resource. If omitted, the
  provider-level region will be used.

* `vip_id` - (Required, ForceNew) The ID of vip to attach the ports to.

* `port_ids` - (Required, List) An array of one or more IDs of the ports to attach the vip to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `vip_subnet_id` - The ID of the subnet this vip connects to.
* `vip_ip_address` - The IP address in the subnet for this vip.
* `ip_addresses` - The IP addresses of ports to attach the vip to.

## Import

Vip associate can be imported using the `vip_id` and port IDs separated by slashes (no limit on the number of
port IDs), e.g.

```
$ terraform import g42cloud_networking_vip_associate.vip_associated vip_id/port1_id/port2_id
```
