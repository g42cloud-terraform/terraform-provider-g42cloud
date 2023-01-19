variable "eps_id" {
  description = "The id of the enterprise project"
  default = "fc04bc67-2b65-4898-ab71-ce0d63fbc267"
}

variable "security_group_name" {
  description = "The name of the security group"
  default = "high-available-web-demo"
}

variable "rds_name" {
  description = "The name of the rds instance"
  default = "high-available-web-demo"
}

variable "rds_flavor_id" {
  description = "The flavor id used to create the rds instance"
  default = "rds.pg.m6.large.8.ha"
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
  default = "NiuzhenguoTest@123"
  sensitive = true
}
