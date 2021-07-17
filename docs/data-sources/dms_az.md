---
subcategory: "Distributed Message Service (DMS)"
---

# g42cloud\_dms\_az

Use this data source to get the ID of an available G42Cloud dms az.

## Example Usage

```hcl

data "g42cloud_dms_az" "az1" {
  code = "ad-ae-1a"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dms az. If omitted, the provider-level region will be used.

* `name` - (Optional, String) Indicates the name of an AZ.

* `code` - (Optional, String) Indicates the code of an AZ.

* `port` - (Optional, String) Indicates the port number of an AZ.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

