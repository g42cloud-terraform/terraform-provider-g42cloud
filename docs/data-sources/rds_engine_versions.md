---
subcategory: "Relational Database Service (RDS)"
---

# g42cloud_rds_engine_versions

Use this data source to obtain all version information of the specified engine type of G42Cloud RDS.

## Example Usage

```hcl
data "g42cloud_rds_engine_versions" "test" {
  type = "SQLServer"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the RDS engine versions.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the RDS engine type.
  The valid values are **MySQL**, **PostgreSQL**, **SQLServer** , default to **MySQL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID in hashcode format.

* `versions` - List of RDS versions. The [versions](#rds_versions) object structure is documented below.

<a name="rds_versions"></a>
The `versions` block contains:

* `id` - Version ID.

* `name` - Version name.
