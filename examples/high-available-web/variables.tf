variable "security_group_name" {
  description = "The name of the security group"
  default = "high-available-web-demo"
}

variable "ecs_name" {
  description = "The name of the ecs instance"
  default = "high-available-web-demo"
}

variable "ecs_flavor_id" {
  description = "The flavor id used to create the ecs instance"
  default = "s6.large.2"
}

variable "ecs_password" {
  description = "The password of the ecs"
  default = "Test@123"
  sensitive = true
}

variable "ecs_volume_size" {
  description = "The volume size of the ecs system disk"
  default = 100
}

variable "ecs_volume_type" {
  description = "The volume type of the ecs system disk"
  default = "SAS"
}

variable "rds_name" {
  description = "The name of the rds instance"
  default = "high-available-web-demo"
}

variable "rds_flavor_id" {
  description = "The flavor id used to create the rds instance"
  default = "rds.mysql.m6.xlarge.8.ha"
}

variable "rds_volume_size" {
  description = "The volume size of the rds instance"
  default = 100
}

variable "rds_volume_type" {
  description = "The volume type of the rds instance"
  default = "ULTRAHIGH"
}

variable "rds_password" {
  description = "The password of the rds instance"
  default = "NiuzhenguoFABTest@123"
  sensitive = true
}

variable "eip_bandwidth_size" {
  description = "The bandwidth size of the eip"
  default = 5
}

variable "elb_name" {
  description = "The name of the elb"
  default = "high-available-web-demo"
}

variable "elb_listen_port" {
  description = "The listen port of the elb"
  default = "80"
}
