---
subcategory: "Virtual Private Cloud (VPC)"
---

# g42cloud_vpc_route

Manages a VPC route resource within G42Cloud.

## Example Usage

### Add route to the default route table

```hcl
variable "vpc_id" {}
variable "nexthop" {}

resource "g42cloud_vpc_route" "vpc_route" {
  vpc_id      = var.vpc_id
  destination = "192.168.0.0/16"
  type        = "peering"
  nexthop     = var.nexthop
}
```

### Add route to a custom route table

```hcl
variable "vpc_id" {}
variable "nexthop" {}

data "g42cloud_vpc_route_table" "rtb" {
  vpc_id = var.vpc_id
  name   = "demo"
}

resource "g42cloud_vpc_route" "vpc_route" {
  vpc_id         = var.vpc_id
  route_table_id = data.g42cloud_vpc_route_table.rtb.id
  destination    = "172.16.8.0/24"
  type           = "ecs"
  nexthop        = var.nexthop
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC route. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `vpc_id` (Required, String, ForceNew) - Specifies the VPC for which a route is to be added. Changing this creates a
  new resource.

* `destination` (Required, String, ForceNew) - Specifies the destination address in the CIDR notation format,
  for example, 192.168.200.0/24. The destination of each route must be unique and cannot overlap with any
  subnet in the VPC. Changing this creates a new resource.

* `type` (Required, String, ForceNew) - Specifies the route type. Currently, the value can only be: **peering**.
  Changing this creates a new resource.

* `nexthop` (Required, String, ForceNew) - Specifies the next hop.
  + If the route type is **peering**, the value is a VPC peering connection ID.

  Changing this creates a new resource.

* `description` (Optional, String, ForceNew) - Specifies the supplementary information about the route.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).
  Changing this creates a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The route ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.
