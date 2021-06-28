---
subcategory: "Distributed Cache Service"
---

# g42cloud\_dcs\_product

Use this data source to get the ID of an available DCS product.

## Example Usage

```hcl

data "g42cloud_dcs_product" "product1" {
  spec_code = "dcs.single_node"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dcs products. If omitted, the provider-level region will be used.

* `spec_code` - (Optional, String) DCS instance specification code. For details, see
[Querying Service Specifications](https://docs.g42cloud.com/api/dcs/dcs-api-0312040.html).


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.
