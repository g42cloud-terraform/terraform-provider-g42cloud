---
subcategory: "Distributed Message Service (DMS)"
---

# g42cloud_dms_product

Use this data source to get the ID of an available G42Cloud dms product.

## Example Usage

```hcl

data "g42cloud_dms_product" "product1" {
  engine            = "kafka"
  version           = "1.1.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dms products. If omitted, the provider-level region
  will be used.

* `engine` - (Required, String) Indicates the name of a message engine. The valid values are **kafka**, **rabbitmq**.

* `instance_type` - (Required, String) Indicates an instance type. The valid values are **single** and **cluster**.

* `version` - (Optional, String) Indicates the version of a message engine.

* `availability_zones` - (Optional, List) Indicates the list of availability zones with available resources.

* `vm_specification` - (Optional, String) Indicates VM specifications.

* `storage` - (Optional, String) Indicates the storage capacity of the resource.
  The default value is the storage capacity of the product.

* `bandwidth` - (Optional, String) Indicates the baseline bandwidth of a DMS instance.
  The valid values are **100MB**, **300MB**, **600MB** and **1200MB**.

* `partition_num` - (Optional, String) Indicates the maximum number of topics that can be created for a Kafka instance.
  The valid values are **300**, **900** and **1800**.

* `storage_spec_code` - (Optional, String) Indicates an I/O specification.
  The valid values are **dms.physical.storage.high** and **dms.physical.storage.ultra**.

* `node_num` - (Optional, String) Indicates the number of nodes in a cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.
