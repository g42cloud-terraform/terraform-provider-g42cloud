#####################################################
# Security Group for web tier
####################################################

resource "g42cloud_networking_secgroup" "secgroup_web" {
  name = "tf_sg_web"
  enterprise_project_id = var.eps_id
}

resource "g42cloud_networking_secgroup_rule" "allow_http" {
  description       = "allow http access"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "80"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = g42cloud_networking_secgroup.secgroup_web.id
}

#####################################################
# Security Group for app tier
####################################################

resource "g42cloud_networking_secgroup" "secgroup_app" {
  name = "tf_sg_app"
  enterprise_project_id = var.eps_id
}

resource "g42cloud_networking_secgroup_rule" "allow_web_app" {
  description       = "allow traffic from web"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = g42cloud_networking_secgroup.secgroup_web.id
  security_group_id = g42cloud_networking_secgroup.secgroup_app.id
}

#####################################################
# Security Group for db tier
####################################################

resource "g42cloud_networking_secgroup" "secgroup_db" {
  name = "tf_sg_db"
  enterprise_project_id = var.eps_id
}

resource "g42cloud_networking_secgroup_rule" "allow_web_db" {
  description       = "allow traffic from web"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = g42cloud_networking_secgroup.secgroup_app.id
  security_group_id = g42cloud_networking_secgroup.secgroup_db.id
}
