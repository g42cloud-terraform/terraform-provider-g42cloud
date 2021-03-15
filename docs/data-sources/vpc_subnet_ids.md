---
subcategory: "Virtual Private Cloud (VPC)"
---

# g42cloud\_vpc\_subnet\_ids

`g42cloud_vpc_subnet_ids` provides a list of subnet ids for a vpc_id

This resource can be useful for getting back a list of subnet ids for a vpc.

## Example Usage

The following example shows outputing all cidr blocks for every subnet id in a vpc.

 ```hcl
data "g42cloud_vpc_subnet_ids" "subnet_ids" {
  vpc_id = var.vpc_id
}

data "g42cloud_vpc_subnet" "subnet" {
  count = length(data.g42cloud_vpc_subnet_ids.subnet_ids.ids)
  id    = tolist(data.g42cloud_vpc_subnet_ids.subnet_ids.ids)[count.index]
 }

output "subnet_cidr_blocks" {
  value = [for s in data.g42cloud_vpc_subnet.subnet: "${s.name}: ${s.id}: ${s.cidr}"]
}
 ```

## Argument Reference

The following arguments are supported:

* `vpc_id` (Required) - Specifies the VPC ID used as the query filter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `ids` - A set of all the subnet ids found. This data source will fail if none are found.
