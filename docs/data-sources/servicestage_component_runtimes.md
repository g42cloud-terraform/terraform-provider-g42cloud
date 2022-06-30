---
subcategory: "ServiceStage"
---

# G42cloud_servicestage_component_runtimes

Use this data source to query available runtimes within G42Cloud.

## Example Usage

```hcl
data "G42cloud_servicestage_component_runtimes" "test" {}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the component runtimes.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the runtime name to use for filtering.
  For the runtime names corresponding to each type of component, please refer to the [document](https://docs.g42cloud.com/usermanual/servicestage/servicestage_user_0411.html).

* `default_port` - (Optional, Int) Specifies the default container port to use for filtering.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `runtimes` - The list of runtime details.

The `runtimes` block contains:

* `name` - The runtime name.

* `default_port` - The default container port.

* `description` - The runtime description.
