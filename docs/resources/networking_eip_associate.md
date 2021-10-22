---
subcategory: "Virtual Private Cloud (VPC)"
---

# g42cloud_networking_eip_associate

Associates an EIP to a port.

## Example Usage

```hcl
variable network_id {}
variable fixed_ip {}

data "g42cloud_networking_port" "myport" {
  network_id = var.network_id
  fixed_ip   = var.fixed_ip
}

resource "g42cloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "g42cloud_networking_eip_associate" "associated" {
  public_ip = g42cloud_vpc_eip.myeip.address
  port_id   = data.g42cloud_networking_port.myport.id
}
```

## Argument Reference

The following arguments are supported:

* `public_ip` - (Required, String, ForceNew) Specifies the EIP address to associate.

* `port_id` - (Required, String, ForceNew) Specifies an existing port ID to associate with this EIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Import

EIP associations can be imported using the `id` of the EIP, e.g.

```
$ terraform import g42cloud_networking_eip_associate.eip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
