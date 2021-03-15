---
subcategory: "Elastic Cloud Server (ECS)"
---

# g42cloud\_compute\_eip\_associate

Associate an EIP to an instance.

## Example Usage

### Automatically detect the correct network

```hcl
resource "g42cloud_compute_instance" "myinstance" {
  name              = "instance"
  image_id          = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id         = "s6.small.1"
  key_pair          = "my_key_pair_name"
  security_groups   = ["default"]
  availability_zone = "ae-ad-1a"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
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

resource "g42cloud_compute_eip_associate" "associated" {
  public_ip   = g42cloud_vpc_eip.myeip.address
  instance_id = g42cloud_compute_instance.myinstance.id
}
```

### Explicitly set the network to attach to

```hcl
resource "g42cloud_compute_instance" "myinstance" {
  name              = "instance"
  image_id          = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id         = "s6.small.1"
  key_pair          = "my_key_pair_name"
  security_groups   = ["default"]
  availability_zone = "ae-ad-1a"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }
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

resource "g42cloud_compute_eip_associate" "associated" {
  public_ip   = g42cloud_vpc_eip.myeip.address
  instance_id = g42cloud_compute_instance.myinstance.id
  fixed_ip    = g42cloud_compute_instance.myinstance.network.1.fixed_ip_v4
}
```

## Argument Reference

The following arguments are supported:

* `public_ip` - (Required, String, ForceNew) The EIP to associate.

* `instance_id` - (Required, String, ForceNew) The instance to associte the EIP with.

* `fixed_ip` - (Optional, String, ForceNew) The specific IP address to direct traffic to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.


## Import

This resource can be imported by specifying all three arguments, separated
by a forward slash:

```
$ terraform import g42cloud_compute_eip_associate.eip_1 <eip>/<instance_id>/<fixed_ip>
```
