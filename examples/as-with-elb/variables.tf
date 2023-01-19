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
