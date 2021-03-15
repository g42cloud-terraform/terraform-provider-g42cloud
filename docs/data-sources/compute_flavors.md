---
subcategory: "Elastic Cloud Server (ECS)"
---

# g42cloud\_compute\_flavors

Use this data source to get the ID of the available Compute Flavors.

## Example Usage

```hcl
data "g42cloud_compute_flavors" "flavors" {
  availability_zone = "ae-ad-1a"
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

# Create ECS instance with the first matched flavor

resource "g42cloud_compute_instance" "instance" {
  flavor_id = data.g42cloud_compute_flavors.flavors.ids[0]

  # Other properties...
}
```

## Argument Reference

In addition to all arguments above, the following attributes are exported:

* `region` - (Optional, String) The region in which to obtain the flavors. If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the AZ name.

* `performance_type` - (Optional, String) Specifies the ECS flavor type.

* `generation` - (Optional, String) Specifies the generation of an ECS type.

* `cpu_core_count` - (Optional, Int) Specifies the number of vCPUs in the ECS flavor.

* `memory_size` - (Optional, Int) Specifies the memory size(GB) in the ECS flavor.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `ids` - A list of flavor IDs.
