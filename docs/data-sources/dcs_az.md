---
subcategory: "Deprecated"
---

# g42cloud_dcs_az

Use this data source to get the ID of an available G42cloud dcs az.

!> **WARNING:** It has been deprecated. This data source is used for the `available_zones` of the
`g42cloud_dcs_instance` resource. Now `available_zones` has been deprecated and this data source is no longer used.

## Example Usage

```hcl
data "g42cloud_dcs_az" "az1" {
  code = "ae-ad-1a"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dcs az. If omitted, the provider-level region will be
  used.

* `code` - (Required, String) Specifies the code of an AZ, e.g. "ae-ad-1a".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID in UUID format.
* `name` - Indicates the name of an AZ.
* `port` - Indicates the port number of an AZ.
