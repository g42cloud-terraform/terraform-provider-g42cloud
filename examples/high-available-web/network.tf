data "g42cloud_vpc" "vpc" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "subnet" {
  name   = "subnet-default"
}

resource "g42cloud_networking_secgroup" "secgroup" {
  name = var.security_group_name

  enterprise_project_id = var.eps_id
}

resource "g42cloud_networking_secgroup_rule" "allow_ssh_linux" {
  description       = "allow ssh access"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "22"
  protocol          = "tcp"
  remote_ip_prefix  = "192.168.10.0/24"
  security_group_id = g42cloud_networking_secgroup.secgroup.id
}

resource "g42cloud_networking_secgroup_rule" "allow_http" {
  description       = "allow http access"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "80"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = g42cloud_networking_secgroup.secgroup.id
}

resource "g42cloud_networking_secgroup_rule" "allow_https" {
  description       = "allow https access"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "443"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = g42cloud_networking_secgroup.secgroup.id
}
