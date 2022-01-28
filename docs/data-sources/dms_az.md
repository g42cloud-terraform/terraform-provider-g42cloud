---
subcategory: "Distributed Message Service (DMS)"
---

# g42cloud_dms_az

Use this data source to get the ID of an available G42Cloud dms az.

!> **WARNING:** It has been deprecated. This data source is used for the `available_zones` of the
`huaweicloud_dms_kafka_instance` and `huaweicloud_dms_rabbitmq_instance` resource.
Now argument `available_zones` has been deprecated, instead `availability_zones`,
this data source will no longer be used.

## Example Usage

```hcl

data "g42cloud_dms_az" "az1" {
  code = "ad-ae-1a"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dms az. If omitted, the provider-level region will be
  used.

* `code` - (Optional, String) Specifies the code of an AZ.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID in UUID format.

* `name` - Indicates the name of an AZ.

* `port` - Indicates the port number of an AZ.

* `ipv6_enabled` - Whether the IPv6 network is enabled.
