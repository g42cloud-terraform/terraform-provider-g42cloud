---
subcategory: "Distributed Message Service (DMS)"
---

# g42cloud_dms_rabbitmq_flavors

Use this data source to get the list of RabbitMQ available flavor details within G42Cloud.

## Example Usage

### Query the list of RabbitMQ flavors by cluster type

```hcl
data "g42cloud_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}
```

### Query the list of RabbitMQ flavors by flavor ID

```hcl
data "g42cloud_dms_rabbitmq_flavors" "test" {
  flavor_id = "c6.2u4g.cluster"
}
```

### Query the list of RabbitMQ flavors by availability zone

```hcl
variable "az1" {}
variable "az2" {}

data "g42cloud_dms_rabbitmq_flavors" "test" {
  availability_zones = [
    var.az1,
    var.az2,
  ]
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the dms RabbitMQ flavors.
  If omitted, the provider-level region will be used.

* `flavor_id` - (Optional, String) Specifies the DMS flavor ID, e.g. **c6.2u4g.cluster**.

* `storage_spec_code` - (Optional, String) Specifies the disk IO encoding.
  + **dms.physical.storage.high.v2**: Type of the disk that uses high I/O.
  + **dms.physical.storage.ultra.v2**: Type of the disk that uses ultra-high I/O.

* `type` - (Optional, String) Specifies flavor type. The valid values are **single** and **cluster**.

* `arch_type` - (Optional, String) Specifies the type of CPU architecture, e.g. **X86**.

* `availability_zones` - (Optional, List) Specifies the list of availability zones with available resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the data source ID.

* `versions` - Indicates the supported flavor versions.

* `flavors` - Indicates the list of flavor details.
  The [flavors](#dms_rabbitmq_flavors) structure is documented below.

<a name="dms_rabbitmq_flavors"></a>
The `flavors` block supports:

* `id` - Indicates the flavor ID.

* `type` - Indicates the flavor type.

* `vm_specification` - Indicates the underlying VM specification.

* `arch_types` - Indicates the list of supported CPU architectures.

* `ios` - Indicates the list of supported disk IO types.
  The [ios](#dms_rabbitmq_flavor_ios) structure is documented below.

* `support_features` - Indicates the list of features supported by the current specification.
  The [support_features](#dms_rabbitmq_flavor_support_features) structure is documented below.

* `properties` - Indicates the properties of the current specification.
  The [properties](#dms_rabbitmq_flavor_properties) structure is documented below.

<a name="dms_rabbitmq_flavor_ios"></a>
The `ios` block supports:

* `storage_spec_code` - Indicates the disk IO encoding.

* `type` - The disk type.

* `availability_zones` - Indicates the list of availability zones with available resources.

* `unavailability_zones` - Indicates the list of unavailability zones with available resources.

<a name="dms_rabbitmq_flavor_support_features"></a>
The `support_features` block supports:

* `name` - Indicates the function name, e.g. **connector_obs**.

* `properties` - Indicates the function property details.
  The [properties](#dms_rabbitmq_flavor_support_feature_properties) structure is documented below.

<a name="dms_rabbitmq_flavor_support_feature_properties"></a>
The `properties` block supports:

* `max_task` - Indicates the maximum number of tasks for the dump function.

* `min_task` - Indicates the minimum number of tasks for the dump function.

* `max_node` - Indicates the maximum number of nodes for the dump function.

* `min_node` - Indicates the minimum number of nodes for the dump function.

<a name="dms_rabbitmq_flavor_properties"></a>
The `properties` block supports:

* `max_broker` - Indicates the maximum number of brokers.

* `min_broker` - Indicates the minimum number of brokers.

* `max_bandwidth_per_broker` - Indicates the maximum bandwidth per broker.

* `max_consumer_per_broker` - Indicates the maximum number of consumers per broker.

* `max_partition_per_broker` - Indicates the maximum number of partitions per broker.

* `max_tps_per_broker` - Indicates the maximum TPS per broker.

* `max_storage_per_node` - Indicates the maximum storage per node. The unit is GB.

* `min_storage_per_node` - Indicates the minimum storage per node. The unit is GB.

* `flavor_alias` - Indicates the flavor ID alias.
