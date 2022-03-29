---
subcategory: "Enterprise Project Management Service (EPS)"
---

# g42cloud_enterprise_project

Use this data source to get an enterprise project from G42Cloud

## Example Usage

```hcl
data "g42cloud_enterprise_project" "test" {
  name = "test"
}
```

## Argument Reference

* `name` - (Optional, String) Specifies the enterprise project name. Fuzzy search is supported.

* `id` - (Optional, String) Specifies the ID of an enterprise project. The value 0 indicates enterprise project default.

* `status` - (Optional, Int) Specifies the status of an enterprise project.
    + 1 indicates Enabled.
    + 2 indicates Disabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Provides supplementary information about the enterprise project.

* `created_at` - Specifies the time (UTC) when the enterprise project was created. Example: 2018-05-18T06:49:06Z

* `updated_at` - Specifies the time (UTC) when the enterprise project was modified. Example: 2018-05-28T02:21:36Z
