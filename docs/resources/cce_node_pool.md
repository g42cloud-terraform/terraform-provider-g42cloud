---
subcategory: "Cloud Container Engine (CCE)"
---

# g42cloud_cce_node_pool

Add a node pool to a container cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "key_pair" {}
variable "availability_zone" {}

resource "g42cloud_cce_node_pool" "node_pool" {
  cluster_id               = var.cluster_id
  name                     = "testpool"
  os                       = "EulerOS 2.5"
  initial_node_count       = 2
  flavor_id                = "s3.large.4"
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair
  scall_enable             = true
  min_node_count           = 1
  max_node_count           = 10
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CCE pool resource. If omitted, the
  provider-level region will be used. Changing this creates a new CCE node pool resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the node pool name.

* `initial_node_count` - (Required, Int) Specifies the initial number of expected nodes in the node pool.
  This parameter can be also used to manually scale the node count afterwards.

* `flavor_id` - (Required, String, ForceNew) Specifies the flavor ID. Changing this parameter will create a new
  resource.

* `type` - (Optional, String, ForceNew) Specifies the node pool type. Possible values are: **vm** and **ElasticBMS**.

* `availability_zone` - (Optional, String, ForceNew) Specifies the name of the available partition (AZ). Default value
  is random to create nodes in a random AZ in the node pool. Changing this parameter will create a new resource.

* `os` - (Optional, String, ForceNew) Specifies the operating system of the node.
  Changing this parameter will create a new resource.

* `key_pair` - (Optional, String, ForceNew) Specifies the key pair name when logging in to select the key pair mode.
  This parameter and `password` are alternative. Changing this parameter will create a new resource.

* `password` - (Optional, String, ForceNew) Specifies the root password when logging in to select the password mode.
  This parameter can be plain or salted and is alternative to `key_pair`.
  Changing this parameter will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) Specifies the ID of the subnet to which the NIC belongs.
  Changing this parameter will create a new resource.

* `max_pods` - (Optional, Int, ForceNew) Specifies the maximum number of instances a node is allowed to create.
  Changing this parameter will create a new resource.

* `preinstall` - (Optional, String, ForceNew) Specifies the script to be executed before installation.
  The input value can be a Base64 encoded string or not. Changing this parameter will create a new resource.

* `postinstall` - (Optional, String, ForceNew) Specifies the script to be executed after installation.
  The input value can be a Base64 encoded string or not. Changing this parameter will create a new resource.

* `extend_param` - (Optional, Map, ForceNew) Specifies the extended parameter.
  Changing this parameter will create a new resource.
  The available keys are as follows:
  + **agency_name**: The agency name to provide temporary credentials for CCE node to access other cloud services.
  + **alpha.cce/NodeImageID**: The custom image ID used to create the BMS nodes.
  + **dockerBaseSize**: The available disk space of a single docker container on the node in device mapper mode.
  + **DockerLVMConfigOverride**: Specifies the data disk configurations of Docker.
  
  The following is an example default configuration:

```hcl
extend_param = {
  DockerLVMConfigOverride = "dockerThinpool=vgpaas/90%VG;kubernetesLV=vgpaas/10%VG;diskType=evs;lvType=linear"
}
```

* `scall_enable` - (Optional, Bool) Specifies whether to enable auto scaling.
  If Autoscaler is enabled, install the autoscaler add-on to use the auto scaling feature.

* `min_node_count` - (Optional, Int) Specifies the minimum number of nodes allowed if auto scaling is enabled.

* `max_node_count` - (Optional, Int) Specifies the maximum number of nodes allowed if auto scaling is enabled.

* `scale_down_cooldown_time` - (Optional, Int) Specifies the time interval between two scaling operations, in minutes.

* `priority` - (Optional, Int) Specifies the weight of the node pool.
  A node pool with a higher weight has a higher priority during scaling.

* `ecs_group_id` - (Optional, String, ForceNew) Specifies the ECS group ID. If specified, the node will be created under
  the cloud server group. Changing this parameter will create a new resource.

* `runtime` - (Optional, String, ForceNew) Specifies the runtime of the CCE node pool. Valid values are *docker* and
  *containerd*. Changing this creates a new resource.

* `labels` - (Optional, Map) Specifies the tags of a Kubernetes node, key/value pair format.

* `tags` - (Optional, Map) Specifies the tags of a VM node, key/value pair format.

* `root_volume` - (Required, List, ForceNew) Specifies the configuration of the system disk.
  The [root_volume](#cce_root_volume)object structure is documented below.
  Changing this parameter will create a new resource.

* `data_volumes` - (Required, List, ForceNew) Specifies the configuration of the data disks.
  The [data_volumes](#cce_data_volumes) object structure is documented below.
  Changing this parameter will create a new resource.

* `storage` - (Optional, List, ForceNew) Specifies the disk initialization management parameter.
  If omitted, disks are managed based on the DockerLVMConfigOverride parameter in extendParam.
  This parameter is supported for clusters of v1.15.11 and later. The [storage](#cce_storage) object structure is
  documented below. Changing this parameter will create a new resource.

* `taints` - (Optional, List) Specifies the taints configuration of the nodes to set anti-affinity.
  The [taints](#cce_taints) object structure is documented below.

* `security_groups` - (Optional, List, ForceNew) Specifies the list of custom security group IDs for the node pool.
  If specified, the nodes will be put in these security groups. When specifying a security group, do not modify
  the rules of the port on which CCE running depends.

* `pod_security_groups` - (Optional, List, ForceNew) Specifies the list of security group IDs for the pod.
  Only supported in CCE Turbo clusters of v1.19 and above. Changing this parameter will create a new resource.

* `initialized_conditions` - (Optional, List) Specifies the custom initialization flags.

<a name="cce_root_volume"></a>
The `root_volume` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size in GB. Changing this parameter will create a new resource.

* `volumetype` - (Required, String, ForceNew) Specifies the disk type. Changing this parameter will create a new resource.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the KMS key ID. This is used to encrypt the volume.
  Changing this parameter will create a new resource.

* `extend_params` - (Optional, Map, ForceNew) Specifies the disk expansion parameters.
  Changing this parameter will create a new resource.

<a name="cce_data_volumes"></a>
The `data_volumes` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size in GB. Changing this parameter will create a new resource.

* `volumetype` - (Required, String, ForceNew) Specifies the disk type. Changing this parameter will create a new resource.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the KMS key ID. This is used to encrypt the volume.
  Changing this parameter will create a new resource.

* `extend_params` - (Optional, Map, ForceNew) Specifies the disk expansion parameters.
  Changing this parameter will create a new resource.

<a name="cce_taints"></a>
The `taints` block supports:

* `key` - (Required, String) A key must contain 1 to 63 characters starting with a letter or digit. Only letters,
  digits, hyphens (-), underscores (_), and periods (.) are allowed. A DNS subdomain name can be used as the
  prefix of a key.

* `value` - (Required, String) A value must start with a letter or digit and can contain a maximum of 63 characters,
  including letters, digits, hyphens (-), underscores (_), and periods (.).

* `effect` - (Required, String) Available options are NoSchedule, PreferNoSchedule, and NoExecute.

<a name="cce_storage"></a>
The `storage` block supports:

* `selectors` - (Required, List, ForceNew) Specifies the disk selection.
  Matched disks are managed according to match labels and storage type. The [selectors](#cce_selectors) object
  structure is documented below. Changing this parameter will create a new resource.

* `groups` - (Required, List, ForceNew) Specifies the storage group consists of multiple storage devices.
  This is used to divide storage space. The [groups](#cce_groups) object structure is documented below.
  Changing this parameter will create a new resource.

<a name="cce_selectors"></a>
The `selectors` block supports:

* `name` - (Required, String, ForceNew) Specifies the selector name, used as the index of `selector_names` in storage group.
  The name of each selector must be unique. Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the storage type. Currently, only **evs (EVS volumes)** is supported.
  The default value is **evs**. Changing this parameter will create a new resource.

* `match_label_size` - (Optional, String, ForceNew) Specifies the matched disk size. If omitted,
  the disk size is not limited. Example: 100. Changing this parameter will create a new resource.

* `match_label_volume_type` - (Optional, String, ForceNew) Specifies the EVS disk type. Currently,
  **SSD**, **GPSSD**, and **SAS** are supported. If omitted, the disk type is not limited.
  Changing this parameter will create a new resource.

* `match_label_metadata_encrypted` - (Optional, String, ForceNew) Specifies the disk encryption identifier.
  Values can be: **0** indicates that the disk is not encrypted and **1** indicates that the disk is encrypted.
  If omitted, whether the disk is encrypted is not limited. Changing this parameter will create a new resource.

* `match_label_metadata_cmkid` - (Optional, String, ForceNew) Specifies the customer master key ID of an encrypted
  disk. Changing this parameter will create a new resource.

* `match_label_count` - (Optional, String, ForceNew) Specifies the number of disks to be selected. If omitted,
  all disks of this type are selected. Changing this parameter will create a new resource.

<a name="cce_groups"></a>
The `groups` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of a virtual storage group. Each group name must be unique.
  Changing this parameter will create a new resource.

* `cce_managed` - (Optional, Bool, ForceNew) Specifies whether the storage space is for **kubernetes** and
  **runtime** components. Only one group can be set to true. The default value is **false**.
  Changing this parameter will create a new resource.

* `selector_names` - (Required, List, ForceNew) Specifies the list of names of selectors to match.
  This parameter corresponds to name in `selectors`. A group can match multiple selectors,
  but a selector can match only one group. Changing this parameter will create a new resource.

* `virtual_spaces` - (Required, List, ForceNew) Specifies the detailed management of space configuration in a group.
  The [virtual_spaces](#cce_virtual_spaces) object structure is documented below.
  Changing this parameter will create a new resource.

<a name="cce_virtual_spaces"></a>
The `virtual_spaces` block supports:

* `name` - (Required, String, ForceNew) Specifies the virtual space name. Currently, only **kubernetes**, **runtime**,
  and **user** are supported. Changing this parameter will create a new resource.

* `size` - (Required, String, ForceNew) Specifies the size of a virtual space. Only an integer percentage is supported.
  Example: 90%. Note that the total percentage of all virtual spaces in a group cannot exceed 100%.
  Changing this parameter will create a new resource.

* `lvm_lv_type` - (Optional, String, ForceNew) Specifies the LVM write mode, values can be **linear** and **striped**.
  This parameter takes effect only in **kubernetes** and **user** configuration. Changing this parameter will create
  a new resource.

* `lvm_path` - (Optional, String, ForceNew) Specifies the absolute path to which the disk is attached.
  This parameter takes effect only in **user** configuration. Changing this parameter will create a new resource.

* `runtime_lv_type` - (Optional, String, ForceNew) Specifies the LVM write mode, values can be **linear** and **striped**.
  This parameter takes effect only in **runtime** configuration. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - Node status information.

* `billing_mode` - Billing mode of a node.

* `current_node_count` - The current number of the nodes.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

CCE node pool can be imported using the cluster ID and node pool ID separated by a slash, e.g.:

```shell
terraform import g42cloud_cce_node_pool.my_node_pool 5c20fdad-7288-11eb-b817-0255ac10158b/e9287dff-7288-11eb-b817-0255ac10158b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `subnet_id`, `preinstall`, `posteinstall`, `taints` and `initial_node_count`.
It is generally recommended running `terraform plan` after importing a node pool.
You can then decide if changes should be applied to the node pool, or the resource
definition should be updated to align with the node pool. Also you can ignore changes as below.

```
resource "g42cloud_cce_node_pool" "my_node_pool" {
    ...

  lifecycle {
    ignore_changes = [
      password, subnet_id,
    ]
  }
}
```
