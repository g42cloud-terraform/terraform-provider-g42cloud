variable "eps_id" {
  description = "The id of the enterprise project"
  default = "fc04bc67-2b65-4898-ab71-ce0d63fbc267"
}

variable "security_group_name" {
  description = "The name of the security group"
  default = "high-available-web-demo"
}

variable "ecs_flavor_id" {
  description = "The flavor id used to create the ecs instance"
  default = "s6.medium.2"
}
