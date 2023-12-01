---
subcategory: "Distributed Database Middleware (DDM)"
---

# g42cloud_ddm_engines

Use this data source to get the list of DDM engines.

## Example Usage

```hcl
data "g42cloud_ddm_engines" test {
  version = "3.0.8.2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `version` - (Optional, String) Specifies the engine version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `engines` - Indicates the list of DDM engine.
  The [engines](#DdmEngines_Engine) object structure is documented below.

<a name="DdmEngines_Engine"></a>
The `engines` block supports:

* `id` - Indicates the ID of the engine.

* `version` - Indicates the engine version.
