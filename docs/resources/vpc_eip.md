---
subcategory: "Elastic IP (EIP)"
---

# g42cloud\_vpc\_eip

Manages a V1 EIP resource within G42Cloud VPC.

## Example Usage

```hcl
resource "g42cloud_vpc_eip" "eip_1" {
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
```

## Example Usage of Share Bandwidth

```hcl
resource "g42cloud_vpc_bandwidth" "bandwidth_1" {
	name = "bandwidth_1"
	size = 5
}

resource "g42cloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    id = g42cloud_vpc_bandwidth.bandwidth_1.id
    share_type = "WHOLE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the eip. If omitted,
    the `region` argument of the provider is used. Changing this creates a new eip.

* `publicip` - (Required, List) The elastic IP address object.

* `bandwidth` - (Required, List) The bandwidth object.


The `publicip` block supports:

* `type` - (Required, String, ForceNew) The value must be a type supported by the system. Only
    `5_bgp` supported now. Changing this creates a new eip.

* `ip_address` - (Optional, String, ForceNew) The value must be a valid IP address in the available
    IP address segment. Changing this creates a new eip.

* `port_id` - (Optional, String) The port id which this eip will associate with. If the value
    is "" or this not specified, the eip will be in unbind state.


The `bandwidth` block supports:

* `name` - (Optional, String) The bandwidth name, which is a string of 1 to 64 characters
    that contain letters, digits, underscores (_), and hyphens (-).

* `size` - (Optional, Int) The bandwidth size. The value ranges from 1 to 300 Mbit/s.

* `id` - (Optional, String, ForceNew) The share bandwidth id. Changing this creates a new eip.

* `share_type` - (Required, String, ForceNew) Whether the bandwidth is shared or exclusive. Changing
    this creates a new eip.

* `charge_mode` - (Optional, String, ForceNew) This is a reserved field. If the system supports charging
    by traffic and this field is specified, then you are charged by traffic for elastic
    IP addresses. Changing this creates a new eip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `address` - The IP address of the eip.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.


## Import

EIPs can be imported using the `id`, e.g.

```
$ terraform import g42cloud_vpc_eip.eip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
