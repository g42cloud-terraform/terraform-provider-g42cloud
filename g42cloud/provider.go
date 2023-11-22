package g42cloud

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/antiddos"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/mutexkv"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/as"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/bms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ces"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ecs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/swr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/tms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpcep"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

func init() {
	waf.PaidType = "postPaid"
	ecs.SystemDiskType = "SAS"
}

// Provider returns a schema.Provider for G42Cloud.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("G42_ACCESS_KEY", nil),
				Description:  descriptions["access_key"],
				RequiredWith: []string{"secret_key"},
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("G42_SECRET_KEY", nil),
				Description:  descriptions["secret_key"],
				RequiredWith: []string{"access_key"},
			},

			"security_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["security_token"],
				RequiredWith: []string{"access_key"},
				DefaultFunc:  schema.EnvDefaultFunc("G42_SECURITY_TOKEN", nil),
			},

			"auth_url": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"G42_AUTH_URL", "https://iam.ae-ad-1.g42cloud.com/v3"),
				Description: descriptions["auth_url"],
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["cloud"],
				DefaultFunc: schema.EnvDefaultFunc(
					"G42_CLOUD", "g42cloud.com"),
			},

			"endpoints": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoints"],
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  descriptions["region"],
				DefaultFunc:  schema.EnvDefaultFunc("G42_REGION_NAME", nil),
				InputDefault: "ae-ad-1",
			},

			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("G42_USERNAME", ""),
				Description:  descriptions["user_name"],
				RequiredWith: []string{"password", "account_name"},
			},

			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"G42_PROJECT_NAME",
				}, ""),
				Description: descriptions["project_name"],
			},

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("G42_PASSWORD", ""),
				Description:  descriptions["password"],
				RequiredWith: []string{"user_name", "account_name"},
			},

			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"G42_ACCOUNT_NAME",
				}, ""),
				Description:  descriptions["account_name"],
				RequiredWith: []string{"password", "user_name"},
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("G42_INSECURE", false),
				Description: descriptions["insecure"],
			},

			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["enterprise_project_id"],
				DefaultFunc: schema.EnvDefaultFunc("G42_ENTERPRISE_PROJECT_ID", ""),
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["max_retries"],
				DefaultFunc: schema.EnvDefaultFunc("G42_MAX_RETRIES", 5),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"g42cloud_apig_environments": apig.DataSourceEnvironments(),

			"g42cloud_as_configurations": as.DataSourceASConfigurations(),
			"g42cloud_as_groups":         as.DataSourceASGroups(),

			"g42cloud_availability_zones": huaweicloud.DataSourceAvailabilityZones(),

			"g42cloud_bms_flavors": bms.DataSourceBmsFlavors(),

			"g42cloud_cbr_vaults": cbr.DataSourceVaults(),

			"g42cloud_cce_addon_template": cce.DataSourceAddonTemplate(),
			"g42cloud_cce_cluster":        cce.DataSourceCCEClusterV3(),
			"g42cloud_cce_clusters":       cce.DataSourceCCEClusters(),
			"g42cloud_cce_node":           cce.DataSourceNode(),
			"g42cloud_cce_node_pool":      cce.DataSourceCCENodePoolV3(),
			"g42cloud_cce_nodes":          cce.DataSourceNodes(),

			"g42cloud_compute_flavors":      ecs.DataSourceEcsFlavors(),
			"g42cloud_compute_instance":     ecs.DataSourceComputeInstance(),
			"g42cloud_compute_instances":    ecs.DataSourceComputeInstances(),
			"g42cloud_compute_servergroups": ecs.DataSourceComputeServerGroups(),

			"g42cloud_css_flavors": css.DataSourceCssFlavors(),

			"g42cloud_dcs_flavors":        dcs.DataSourceDcsFlavorsV2(),
			"g42cloud_dcs_maintainwindow": dcs.DataSourceDcsMaintainWindow(),
			"g42cloud_dcs_product":        deprecated.DataSourceDcsProductV1(),
			"g42cloud_dcs_az":             deprecated.DataSourceDcsAZV1(),

			"g42cloud_dds_flavors": dds.DataSourceDDSFlavorV3(),

			"g42cloud_kms_key":      dew.DataSourceKmsKey(),
			"g42cloud_kms_data_key": dew.DataSourceKmsDataKeyV1(),

			"g42cloud_dms_product":         dms.DataSourceDmsProduct(),
			"g42cloud_dms_maintainwindow":  dms.DataSourceDmsMaintainWindow(),
			"g42cloud_dms_kafka_flavors":   dms.DataSourceKafkaFlavors(),
			"g42cloud_dms_kafka_instances": dms.DataSourceDmsKafkaInstances(),
			"g42cloud_dms_az":              deprecated.DataSourceDmsAZ(),

			"g42cloud_elb_certificate": elb.DataSourceELBCertificateV3(),
			"g42cloud_elb_flavors":     elb.DataSourceElbFlavorsV3(),

			"g42cloud_enterprise_project": eps.DataSourceEnterpriseProject(),

			"g42cloud_identity_role": iam.DataSourceIdentityRole(),

			"g42cloud_images_image":  ims.DataSourceImagesImageV2(),
			"g42cloud_images_images": ims.DataSourceImagesImages(),

			"g42cloud_lb_loadbalancer": lb.DataSourceELBV2Loadbalancer(),
			"g42cloud_lb_certificate":  lb.DataSourceLBCertificateV2(),
			"g42cloud_lb_pools":        lb.DataSourcePools(),

			"g42cloud_modelarts_datasets":         modelarts.DataSourceDatasets(),
			"g42cloud_modelarts_dataset_versions": modelarts.DataSourceDatasetVerions(),
			"g42cloud_modelarts_notebook_images":  modelarts.DataSourceNotebookImages(),

			"g42cloud_nat_gateway": nat.DataSourcePublicGateway(),

			"g42cloud_networking_port":     vpc.DataSourceNetworkingPortV2(),
			"g42cloud_networking_secgroup": vpc.DataSourceNetworkingSecGroup(),

			"g42cloud_obs_bucket_object": obs.DataSourceObsBucketObject(),

			"g42cloud_rds_flavors": rds.DataSourceRdsFlavor(),

			"g42cloud_rms_policy_definitions": rms.DataSourcePolicyDefinitions(),

			"g42cloud_servicestage_component_runtimes": servicestage.DataSourceComponentRuntimes(),

			"g42cloud_sms_source_servers": sms.DataSourceServers(),

			"g42cloud_vpc_bandwidth": eip.DataSourceBandWidth(),
			"g42cloud_vpc_eip":       eip.DataSourceVpcEip(),
			"g42cloud_vpc_eips":      eip.DataSourceVpcEips(),

			"g42cloud_vpc":             vpc.DataSourceVpcV1(),
			"g42cloud_vpc_subnet":      vpc.DataSourceVpcSubnetV1(),
			"g42cloud_vpc_subnet_ids":  vpc.DataSourceVpcSubnetIdsV1(),
			"g42cloud_vpc_route":       vpc.DataSourceVpcRouteV2(),
			"g42cloud_vpc_route_table": vpc.DataSourceVPCRouteTable(),

			"g42cloud_vpcep_public_services": vpcep.DataSourceVPCEPPublicServices(),

			"g42cloud_waf_certificate":         waf.DataSourceWafCertificateV1(),
			"g42cloud_waf_policies":            waf.DataSourceWafPoliciesV1(),
			"g42cloud_waf_dedicated_instances": waf.DataSourceWafDedicatedInstancesV1(),
			"g42cloud_waf_reference_tables":    waf.DataSourceWafReferenceTablesV1(),

			// Legacy
			"g42cloud_identity_role_v3": iam.DataSourceIdentityRole(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"g42cloud_aom_alarm_rule":             aom.ResourceAlarmRule(),
			"g42cloud_aom_service_discovery_rule": aom.ResourceServiceDiscoveryRule(),
			"g42cloud_aom_event_alarm_rule":       aom.ResourceEventAlarmRule(),
			"g42cloud_aom_alarm_action_rule":      aom.ResourceAlarmActionRule(),
			"g42cloud_aom_alarm_silence_rule":     aom.ResourceAlarmSilenceRule(),

			"g42cloud_antiddos_basic": antiddos.ResourceCloudNativeAntiDdos(),

			"g42cloud_apig_api":                         apig.ResourceApigAPIV2(),
			"g42cloud_apig_api_publishment":             apig.ResourceApigApiPublishment(),
			"g42cloud_apig_instance":                    apig.ResourceApigInstanceV2(),
			"g42cloud_apig_application":                 apig.ResourceApigApplicationV2(),
			"g42cloud_apig_custom_authorizer":           apig.ResourceApigCustomAuthorizerV2(),
			"g42cloud_apig_environment":                 apig.ResourceApigEnvironmentV2(),
			"g42cloud_apig_group":                       apig.ResourceApigGroupV2(),
			"g42cloud_apig_response":                    apig.ResourceApigResponseV2(),
			"g42cloud_apig_throttling_policy_associate": apig.ResourceThrottlingPolicyAssociate(),
			"g42cloud_apig_throttling_policy":           apig.ResourceApigThrottlingPolicyV2(),
			"g42cloud_apig_vpc_channel":                 deprecated.ResourceApigVpcChannelV2(),

			"g42cloud_as_bandwidth_policy": as.ResourceASBandWidthPolicy(),
			"g42cloud_as_configuration":    as.ResourceASConfiguration(),
			"g42cloud_as_group":            as.ResourceASGroup(),
			"g42cloud_as_instance_attach":  as.ResourceASInstanceAttach(),
			"g42cloud_as_lifecycle_hook":   as.ResourceASLifecycleHook(),
			"g42cloud_as_notification":     as.ResourceAsNotification(),
			"g42cloud_as_policy":           as.ResourceASPolicy(),

			"g42cloud_bms_instance": bms.ResourceBmsInstance(),

			"g42cloud_cbr_policy": cbr.ResourcePolicy(),
			"g42cloud_cbr_vault":  cbr.ResourceVault(),

			"g42cloud_cce_addon":     cce.ResourceAddon(),
			"g42cloud_cce_cluster":   cce.ResourceCluster(),
			"g42cloud_cce_namespace": cce.ResourceCCENamespaceV1(),
			"g42cloud_cce_node_pool": cce.ResourceNodePool(),
			"g42cloud_cce_node":      cce.ResourceNode(),
			"g42cloud_cce_pvc":       cce.ResourceCcePersistentVolumeClaimsV1(),

			"g42cloud_ces_alarmrule":      ces.ResourceAlarmRule(),
			"g42cloud_ces_alarm_template": ces.ResourceCesAlarmTemplate(),
			"g42cloud_ces_resource_group": ces.ResourceResourceGroup(),

			"g42cloud_compute_instance":         ecs.ResourceComputeInstance(),
			"g42cloud_compute_interface_attach": ecs.ResourceComputeInterfaceAttach(),
			"g42cloud_compute_servergroup":      ecs.ResourceComputeServerGroup(),
			"g42cloud_compute_eip_associate":    ecs.ResourceComputeEIPAssociate(),
			"g42cloud_compute_volume_attach":    ecs.ResourceComputeVolumeAttach(),
			"g42cloud_compute_keypair":          huaweicloud.ResourceComputeKeypairV2(),

			"g42cloud_csms_secret": dew.ResourceCsmsSecret(),

			"g42cloud_css_cluster":       css.ResourceCssCluster(),
			"g42cloud_css_snapshot":      css.ResourceCssSnapshot(),
			"g42cloud_css_thesaurus":     css.ResourceCssthesaurus(),
			"g42cloud_css_configuration": css.ResourceCssConfiguration(),

			"g42cloud_cts_tracker":      cts.ResourceCTSTracker(),
			"g42cloud_cts_data_tracker": cts.ResourceCTSDataTracker(),

			"g42cloud_dc_virtual_gateway":   dc.ResourceVirtualGateway(),
			"g42cloud_dc_virtual_interface": dc.ResourceVirtualInterface(),

			"g42cloud_dcs_instance": dcs.ResourceDcsInstance(),

			"g42cloud_dds_instance": dds.ResourceDdsInstanceV3(),

			"g42cloud_dli_queue": dli.ResourceDliQueue(),

			"g42cloud_dms_instance":          deprecated.ResourceDmsInstancesV1(),
			"g42cloud_dms_kafka_instance":    dms.ResourceDmsKafkaInstance(),
			"g42cloud_dms_kafka_topic":       dms.ResourceDmsKafkaTopic(),
			"g42cloud_dms_kafka_user":        dms.ResourceDmsKafkaUser(),
			"g42cloud_dms_kafka_permissions": dms.ResourceDmsKafkaPermissions(),
			"g42cloud_dms_rabbitmq_instance": dms.ResourceDmsRabbitmqInstance(),

			"g42cloud_dns_ptrrecord": dns.ResourceDNSPtrRecord(),
			"g42cloud_dns_recordset": dns.ResourceDNSRecordSetV2(),
			"g42cloud_dns_zone":      dns.ResourceDNSZone(),

			"g42cloud_dws_cluster": dws.ResourceDwsCluster(),

			"g42cloud_elb_certificate":  elb.ResourceCertificateV3(),
			"g42cloud_elb_l7policy":     elb.ResourceL7PolicyV3(),
			"g42cloud_elb_l7rule":       elb.ResourceL7RuleV3(),
			"g42cloud_elb_listener":     elb.ResourceListenerV3(),
			"g42cloud_elb_loadbalancer": elb.ResourceLoadBalancerV3(),
			"g42cloud_elb_monitor":      elb.ResourceMonitorV3(),
			"g42cloud_elb_ipgroup":      elb.ResourceIpGroupV3(),
			"g42cloud_elb_pool":         elb.ResourcePoolV3(),
			"g42cloud_elb_member":       elb.ResourceMemberV3(),

			"g42cloud_enterprise_project": eps.ResourceEnterpriseProject(),
			"g42cloud_evs_snapshot":       evs.ResourceEvsSnapshotV2(),
			"g42cloud_evs_volume":         evs.ResourceEvsVolume(),

			"g42cloud_fgs_function": fgs.ResourceFgsFunctionV2(),

			"g42cloud_identity_role_assignment":  iam.ResourceIdentityGroupRoleAssignment(),
			"g42cloud_identity_user":             iam.ResourceIdentityUser(),
			"g42cloud_identity_group":            iam.ResourceIdentityGroup(),
			"g42cloud_identity_group_membership": iam.ResourceIdentityGroupMembership(),
			"g42cloud_identity_acl":              iam.ResourceIdentityACL(),
			"g42cloud_identity_agency":           iam.ResourceIAMAgencyV3(),
			"g42cloud_identity_project":          iam.ResourceIdentityProject(),
			"g42cloud_identity_role":             iam.ResourceIdentityRole(),

			"g42cloud_images_image": ims.ResourceImsImage(),

			"g42cloud_kms_key": dew.ResourceKmsKey(),

			"g42cloud_lb_certificate":  lb.ResourceCertificateV2(),
			"g42cloud_lb_l7policy":     lb.ResourceL7PolicyV2(),
			"g42cloud_lb_l7rule":       lb.ResourceL7RuleV2(),
			"g42cloud_lb_listener":     lb.ResourceListener(),
			"g42cloud_lb_loadbalancer": lb.ResourceLoadBalancer(),
			"g42cloud_lb_member":       lb.ResourceMemberV2(),
			"g42cloud_lb_monitor":      lb.ResourceMonitorV2(),
			"g42cloud_lb_pool":         lb.ResourcePoolV2(),
			"g42cloud_lb_whitelist":    lb.ResourceWhitelistV2(),

			"g42cloud_lts_group":  lts.ResourceLTSGroup(),
			"g42cloud_lts_stream": lts.ResourceLTSStream(),

			"g42cloud_mapreduce_cluster": mrs.ResourceMRSClusterV2(),
			"g42cloud_mapreduce_job":     mrs.ResourceMRSJobV2(),

			"g42cloud_modelarts_dataset":                modelarts.ResourceDataset(),
			"g42cloud_modelarts_dataset_version":        modelarts.ResourceDatasetVersion(),
			"g42cloud_modelarts_notebook":               modelarts.ResourceNotebook(),
			"g42cloud_modelarts_notebook_mount_storage": modelarts.ResourceNotebookMountStorage(),

			"g42cloud_nat_dnat_rule": nat.ResourcePublicDnatRule(),
			"g42cloud_nat_gateway":   nat.ResourcePublicGateway(),
			"g42cloud_nat_snat_rule": nat.ResourcePublicSnatRule(),

			"g42cloud_network_acl":              huaweicloud.ResourceNetworkACL(),
			"g42cloud_network_acl_rule":         huaweicloud.ResourceNetworkACLRule(),
			"g42cloud_networking_eip_associate": eip.ResourceEIPAssociate(),
			"g42cloud_networking_secgroup":      vpc.ResourceNetworkingSecGroup(),
			"g42cloud_networking_secgroup_rule": vpc.ResourceNetworkingSecGroupRule(),
			"g42cloud_networking_vip":           vpc.ResourceNetworkingVip(),
			"g42cloud_networking_vip_associate": vpc.ResourceNetworkingVIPAssociateV2(),

			"g42cloud_obs_bucket":        obs.ResourceObsBucket(),
			"g42cloud_obs_bucket_object": obs.ResourceObsBucketObject(),
			"g42cloud_obs_bucket_policy": obs.ResourceObsBucketPolicy(),

			"g42cloud_rms_policy_assignment":                  rms.ResourcePolicyAssignment(),
			"g42cloud_rms_resource_aggregator":                rms.ResourceAggregator(),
			"g42cloud_rms_resource_aggregation_authorization": rms.ResourceAggregationAuthorization(),
			"g42cloud_rms_resource_recorder":                  rms.ResourceRecorder(),

			"g42cloud_rds_instance":              ResourceRdsInstanceV3(),
			"g42cloud_rds_parametergroup":        rds.ResourceRdsConfiguration(),
			"g42cloud_rds_read_replica_instance": rds.ResourceRdsReadReplicaInstance(),

			"g42cloud_servicestage_application":                 servicestage.ResourceApplication(),
			"g42cloud_servicestage_component_instance":          servicestage.ResourceComponentInstance(),
			"g42cloud_servicestage_component":                   servicestage.ResourceComponent(),
			"g42cloud_servicestage_environment":                 servicestage.ResourceEnvironment(),
			"g42cloud_servicestage_repo_token_authorization":    servicestage.ResourceRepoTokenAuth(),
			"g42cloud_servicestage_repo_password_authorization": servicestage.ResourceRepoPwdAuth(),

			"g42cloud_sfs_turbo": sfs.ResourceSFSTurbo(),

			"g42cloud_smn_subscription": smn.ResourceSubscription(),
			"g42cloud_smn_topic":        smn.ResourceTopic(),

			"g42cloud_sms_server_template": sms.ResourceServerTemplate(),
			"g42cloud_sms_task":            sms.ResourceMigrateTask(),

			"g42cloud_swr_organization":             swr.ResourceSWROrganization(),
			"g42cloud_swr_organization_permissions": swr.ResourceSWROrganizationPermissions(),
			"g42cloud_swr_repository":               swr.ResourceSWRRepository(),
			"g42cloud_swr_repository_sharing":       swr.ResourceSWRRepositorySharing(),

			"g42cloud_tms_tags": tms.ResourceTmsTag(),

			"g42cloud_vpc_bandwidth":     eip.ResourceVpcBandWidthV2(),
			"g42cloud_vpc_eip":           eip.ResourceVpcEIPV1(),
			"g42cloud_vpc_eip_associate": eip.ResourceEIPAssociate(),

			"g42cloud_vpc":                             vpc.ResourceVirtualPrivateCloudV1(),
			"g42cloud_vpc_route":                       vpc.ResourceVPCRouteTableRoute(),
			"g42cloud_vpc_route_table":                 vpc.ResourceVPCRouteTable(),
			"g42cloud_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"g42cloud_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"g42cloud_vpc_subnet":                      vpc.ResourceVpcSubnetV1(),
			"g42cloud_vpc_address_group":               vpc.ResourceVpcAddressGroup(),

			"g42cloud_vpcep_approval": vpcep.ResourceVPCEndpointApproval(),
			"g42cloud_vpcep_endpoint": vpcep.ResourceVPCEndpoint(),
			"g42cloud_vpcep_service":  vpcep.ResourceVPCEndpointService(),

			"g42cloud_waf_certificate":                waf.ResourceWafCertificateV1(),
			"g42cloud_waf_domain":                     waf.ResourceWafDomainV1(),
			"g42cloud_waf_policy":                     waf.ResourceWafPolicyV1(),
			"g42cloud_waf_rule_blacklist":             waf.ResourceWafRuleBlackListV1(),
			"g42cloud_waf_rule_data_masking":          waf.ResourceWafRuleDataMaskingV1(),
			"g42cloud_waf_rule_web_tamper_protection": waf.ResourceWafRuleWebTamperProtectionV1(),
			"g42cloud_waf_dedicated_instance":         waf.ResourceWafDedicatedInstance(),
			"g42cloud_waf_dedicated_domain":           waf.ResourceWafDedicatedDomainV1(),
			"g42cloud_waf_reference_table":            waf.ResourceWafReferenceTableV1(),
			"g42cloud_waf_rule_cc_protection":         waf.ResourceRuleCCProtection(),
			"g42cloud_waf_rule_precise_protection":    waf.ResourceRulePreciseProtection(),

			// Legacy
			"g42cloud_identity_role_assignment_v3":  iam.ResourceIdentityGroupRoleAssignment(),
			"g42cloud_identity_user_v3":             iam.ResourceIdentityUser(),
			"g42cloud_identity_group_v3":            iam.ResourceIdentityGroup(),
			"g42cloud_identity_group_membership_v3": iam.ResourceIdentityGroupMembership(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"cloud": "The endpoint of cloud provider, defaults to g42cloud.com",

		"endpoints": "The custom endpoints used to override the default endpoint URL.",

		"region": "The G42Cloud region to connect to.",

		"access_key": "The access key of the G42Cloud to use.",

		"secret_key": "The secret key of the G42Cloud to use.",

		"security_token": "The security token to authenticate with a temporary security credential.",

		"user_name": "Username to login with.",

		"project_name": "The name of the Project to login with.",

		"password": "Password to login with.",

		"account_name": "The name of the Account to login with.",

		"insecure": "Trust self-signed certificates.",
	}
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	var project_name string

	region := d.Get("region").(string)

	// Use region as project_name if it's not set
	if v, ok := d.GetOk("project_name"); ok && v.(string) != "" {
		project_name = v.(string)
	} else {
		project_name = region
	}

	config := config.Config{
		AccessKey:           d.Get("access_key").(string),
		SecretKey:           d.Get("secret_key").(string),
		SecurityToken:       d.Get("security_token").(string),
		DomainName:          d.Get("account_name").(string),
		IdentityEndpoint:    d.Get("auth_url").(string),
		Insecure:            d.Get("insecure").(bool),
		Password:            d.Get("password").(string),
		Region:              d.Get("region").(string),
		TenantName:          project_name,
		Username:            d.Get("user_name").(string),
		TerraformVersion:    terraformVersion,
		Cloud:               d.Get("cloud").(string),
		MaxRetries:          d.Get("max_retries").(int),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		RegionClient:        true,
		RegionProjectIDMap:  make(map[string]string),
		RPLock:              new(sync.Mutex),
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	if config.HwClient != nil && config.HwClient.ProjectID != "" {
		config.RegionProjectIDMap[config.Region] = config.HwClient.ProjectID
	}

	// get custom endpoints
	endpoints, err := flattenProviderEndpoints(d)
	if err != nil {
		return nil, err
	}

	// set default endpoints
	if _, ok := endpoints["sms"]; !ok {
		endpoints["sms"] = fmt.Sprintf("https://sms.%s.%s/", region, config.Cloud)
	}

	config.Endpoints = endpoints

	return &config, nil
}

func flattenProviderEndpoints(d *schema.ResourceData) (map[string]string, error) {
	endpoints := d.Get("endpoints").(map[string]interface{})
	epMap := make(map[string]string)

	for key, val := range endpoints {
		endpoint := strings.TrimSpace(val.(string))
		// check empty string
		if endpoint == "" {
			return nil, fmt.Errorf("the value of customer endpoint %s must be specified", key)
		}

		// add prefix "https://" and suffix "/"
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = fmt.Sprintf("%s/", endpoint)
		}
		epMap[key] = endpoint
	}

	// unify the endpoint which has multi types
	if endpoint, ok := epMap["iam"]; ok {
		epMap["identity"] = endpoint
	}
	if endpoint, ok := epMap["ecs"]; ok {
		epMap["ecsv11"] = endpoint
		epMap["ecsv21"] = endpoint
	}
	if endpoint, ok := epMap["cce"]; ok {
		epMap["cce_addon"] = endpoint
	}
	if endpoint, ok := epMap["evs"]; ok {
		epMap["volumev2"] = endpoint
	}
	if endpoint, ok := epMap["vpc"]; ok {
		epMap["networkv2"] = endpoint
		epMap["security_group"] = endpoint
	}

	log.Printf("[DEBUG] customer endpoints: %+v", epMap)
	return epMap, nil
}
