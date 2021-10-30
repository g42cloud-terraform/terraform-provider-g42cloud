## 1.2.0 (October 30, 2021)

FEATURES:

* **New Resource:** `g42cloud_networking_eip_associate` ([#42](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/42))

ENHANCEMENTS:

* Update the reference of sdk and huaweicloud ([#39](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/39))
* Update dcs resources ([#41](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/41))

## 1.1.0 (September 30, 2021)

FEATURES:

* **New Data Source:** `g42cloud_obs_bucket_object` ([#33](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/33))
* **New Data Source:** `g42cloud_rds_flavors` ([#34](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/34))
* **New Data Source:** `g42cloud_dms_az` ([#35](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/35))
* **New Data Source:** `g42cloud_dms_maintainwindow` ([#35](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/35))
* **New Data Source:** `g42cloud_dms_product` ([#35](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/35))
* **New Resource:** `g42cloud_obs_bucket_object` ([#33](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/33))
* **New Resource:** `g42cloud_obs_bucket_policy` ([#33](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/33))
* **New Resource:** `g42cloud_obs_bucket` ([#33](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/33))
* **New Resource:** `g42cloud_rds_configuration` ([#34](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/34))
* **New Resource:** `g42cloud_rds_instance` ([#34](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/34))
* **New Resource:** `g42cloud_rds_read_replica_instance` ([#34](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/34))
* **New Resource:** `g42cloud_dms_instance` ([#35](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/35))
* **New Resource:** `g42cloud_dli_queue` ([#36](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/36))
* **New Resource:** `g42cloud_ces_alarmrule` ([#37](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/37))

## 1.0.0 (July 8, 2021)

FEATURES:

* **New Data Source:** `g42cloud_kms_key` ([#26](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/26))
* **New Data Source:** `g42cloud_kms_data_key` ([#26](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/26))
* **New Data Source:** `g42cloud_dcs_az` ([#27](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/27))
* **New Data Source:** `g42cloud_dcs_maintainwindow` ([#27](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/27))
* **New Data Source:** `g42cloud_dcs_product` ([#27](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/27))
* **New Data Source:** `g42cloud_dds_flavors` ([#28](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/28))
* **New Resource:** `g42cloud_fgs_function` ([#25](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/25))
* **New Resource:** `g42cloud_kms_key` ([#26](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/26))
* **New Resource:** `g42cloud_dcs_instance` ([#27](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/27))
* **New Resource:** `g42cloud_dds_instance` ([#28](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/28))
* **New Resource:** `g42cloud_smn_topic` ([#29](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/29))
* **New Resource:** `g42cloud_smn_subscription` ([#29](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/29))

ENHANCEMENTS:

* provider: Make log message configurable ([#24](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/24))

## 0.4.0 (May 31, 2021)

FEATURES:

* **New Data Source:** `g42cloud_cce_cluster` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))
* **New Data Source:** `g42cloud_cce_node` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))
* **New Data Source:** `g42cloud_cce_node_pool` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))
* **New Data Source:** `g42cloud_cce_addon_template` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))
* **New Resource:** `g42cloud_identity_acl` ([#22](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/22))
* **New Resource:** `g42cloud_identity_agency` ([#22](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/22))
* **New Resource:** `g42cloud_identity_role` ([#22](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/22))
* **New Resource:** `g42cloud_cce_cluster` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))
* **New Resource:** `g42cloud_cce_node` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))
* **New Resource:** `g42cloud_cce_node_pool` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))
* **New Resource:** `g42cloud_cce_addon` ([#21](https://github.com/g42cloud-terraform/terraform-provider-g42cloud/pull/21))

## 0.3.0 (March 2, 2021)

ENHANCEMENTS:

* provider: Add custom cloud and endpoints support (#10)

## 0.2.0 (February 9, 2021)

FEATURES:

* **New Resource:** `g42cloud_as_configuration`
* **New Resource:** `g42cloud_as_group`
* **New Resource:** `g42cloud_as_policy`
* **New Resource:** `g42cloud_network_acl`
* **New Resource:** `g42cloud_network_acl_rule`

## 0.1.0 (December 29, 2020)

FEATURES:

* **New Data Source:** `g42cloud_availability_zones`
* **New Data Source:** `g42cloud_compute_flavors`
* **New Data Source:** `g42cloud_identity_role`
* **New Data Source:** `g42cloud_images_image`
* **New Data Source:** `g42cloud_networking_port`
* **New Data Source:** `g42cloud_networking_secgroup`
* **New Data Source:** `g42cloud_vpc`
* **New Data Source:** `g42cloud_vpc_bandwidth`
* **New Data Source:** `g42cloud_vpc_route`
* **New Data Source:** `g42cloud_vpc_subnet`
* **New Data Source:** `g42cloud_vpc_subnet_ids`
* **New Resource:** `g42cloud_dns_recordset`
* **New Resource:** `g42cloud_dns_zone`
* **New Resource:** `g42cloud_identity_role_assignment`
* **New Resource:** `g42cloud_identity_user`
* **New Resource:** `g42cloud_identity_group`
* **New Resource:** `g42cloud_identity_group_membership`
* **New Resource:** `g42cloud_images_image`
* **New Resource:** `g42cloud_compute_instance`
* **New Resource:** `g42cloud_compute_interface_attach`
* **New Resource:** `g42cloud_compute_keypair`
* **New Resource:** `g42cloud_compute_servergroup`
* **New Resource:** `g42cloud_compute_eip_associate`
* **New Resource:** `g42cloud_compute_volume_attach`
* **New Resource:** `g42cloud_evs_snapshot`
* **New Resource:** `g42cloud_evs_volume`
* **New Resource:** `g42cloud_lb_certificate`
* **New Resource:** `g42cloud_lb_l7policy`
* **New Resource:** `g42cloud_lb_l7rule`
* **New Resource:** `g42cloud_lb_listener`
* **New Resource:** `g42cloud_lb_loadbalancer`
* **New Resource:** `g42cloud_lb_member`
* **New Resource:** `g42cloud_lb_monitor`
* **New Resource:** `g42cloud_lb_pool`
* **New Resource:** `g42cloud_lb_whitelist`
* **New Resource:** `g42cloud_networking_secgroup`
* **New Resource:** `g42cloud_networking_secgroup_rule`
* **New Resource:** `g42cloud_vpc`
* **New Resource:** `g42cloud_vpc_eip`
* **New Resource:** `g42cloud_vpc_subnet`
* **New Resource:** `g42cloud_vpc_route`
* **New Resource:** `g42cloud_vpc_peering_connection`
