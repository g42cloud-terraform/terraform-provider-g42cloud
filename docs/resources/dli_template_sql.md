---
subcategory: "Data Lake Insight (DLI)"
---

# g42cloud_dli_template_sql

Manages a DLI SQL template resource within G42Cloud.

## Example Usage

```hcl
variable "sql" {}

resource "g42cloud_dli_template_sql" "test" {
  name        = "demo"
  sql         = var.sql
  group       = "test"
  description = "This is a demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the SQL template.

* `sql` - (Required, String) The statement of the SQL template.

* `description` - (Optional, String) The description of the SQL template.

* `group` - (Optional, String) The group of the SQL template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `owner` - The user ID of owner.

## Import

The SQL template can be imported using the `id`, e.g.

```bash
$ terraform import g42cloud_dli_template_sql.test 0ce123456a00f2591fabc00385ff1234
```
