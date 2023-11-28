package main

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// deprecatedFields includes the fields that shoud always be ignored.
	deprecatedFields = []string{"tenant_id", "admin_state_up",
		// g42cloud does not support prepaid mode
		"charging_mode", "period", "period_unit", "auto_renew", "auto_pay",
	}

	specialProductNameMap = map[string]string{
		"GaussDBforopenGauss": "GaussDB(for openGauss)",
	}

	specialShortNameMap = map[string]string{
		"compute":    "ECS",
		"images":     "IMS",
		"identity":   "IAM",
		"networking": "VPC",
		"network":    "VPC",
		"vpcs":       "VPC", // specially for data.g42cloud_vpcs
		"lb":         "ELB",
		"api":        "APIG",
		"fgs":        "FunctionGraph",
		"enterprise": "EPS",
		"mapreduce":  "MRS",

		"antiddos":     "Anti-DDoS",
		"cloudtable":   "CloudTable",
		"codehub":      "CodeHub",
		"iotda":        "IoTDA",
		"gaussdb":      "GaussDB",
		"dataarts":     "DataArtsStudio",
		"modelarts":    "ModelArts",
		"projectman":   "ProjectMan",
		"servicestage": "ServiceStage",
		"secmaster":    "SecMaster",

		"live":      "Live",
		"meeting":   "Meeting",
		"workspace": "Workspace",
	}

	// if the resource name **equals** the key, then return the product name
	specialResourceMap = map[string]string{
		"g42cloud_vpc_eip":           "EIP",
		"g42cloud_vpc_eip_associate": "EIP",
		"g42cloud_vpc_bandwidth":     "EIP",
	}

	// if the resource name **contains** the key, then return the product name
	specialResourceKeyMap = map[string]string{
		"_cnad_advanced_": "AAD",
		"_dms_kafka_":     "Kafka",
		"_dms_rabbitmq_":  "RabbitMQ",
		"_dms_rocketmq_":  "RocketMQ",

		"_gaussdb_cassandra_": "GaussDBforNoSQL",
		"_gaussdb_influx_":    "GaussDBforNoSQL",
		"_gaussdb_mongo_":     "GaussDBforNoSQL",
		"_gaussdb_redis_":     "GaussDBforNoSQL",
		"_gaussdb_opengauss_": "GaussDBforopenGauss",
	}
)

func isInternalResource(resource *schema.Resource, key string) bool {
	// lb_xxx is not used any more
	if strings.Contains(key, "_lb_") {
		return true
	}

	if resource.Description != "" {
		// get extent attributes from description
		extent := parseExtentAttribute(resource.Description)
		if hasExtentAttribute(extent, "Internal") {
			return true
		}
	}

	return false
}

func isDeprecatedField(field string) bool {
	for _, key := range deprecatedFields {
		if field == key {
			return true
		}
	}
	return false
}
