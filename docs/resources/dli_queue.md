---
subcategory: "Data Lake Insight (DLI)"
---

# g42cloud_dli_queue

Manages a DLI queue resource within G42Cloud.

## Example Usage

### create a queue

```hcl
resource "g42cloud_dli_queue" "queue" {
  name     = "terraform_dli_queue_test"
  cu_count = 4
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DLI queue resource.
  If omitted, the provider-level region will be used. Changing this creates a new DLI Queue resource.

* `cu_count` - (Required, Int, ForceNew) Specifies the minimum number of CUs that are bound to a queue. The value can be
  4, 16, or 64. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of a queue. The name can contain only digits, letters, and
  underscores (_), but cannot contain only digits or start with an
  underscore (_). Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of a queue.
  Changing this parameter will create a new resource.

* `management_subnet_cidr` - (Optional, String, ForceNew) Specifies the CIDR of the management subnet.
  Changing this parameter will create a new resource.

* `subnet_cidr` - (Optional, String, ForceNew) Specifies the subnet CIDR. Changing this parameter will create a new resource.

* `vpc_cidr` - (Optional, String, ForceNew) Specifies the VPC CIDR. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `create_time` -  Time when a queue is created.
