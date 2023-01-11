# Fetching the availability zones of the region
data "g42cloud_availability_zones" "az" {}

# Feching the images with the specificed filters
data "g42cloud_images_image" "image" {
  name = "CentOS 7.9 64bit"
  visibility = "public"
  most_recent = true
}

data "g42cloud_vpc" "vpc" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "subnet" {
  name   = "subnet-default"
}

resource "g42cloud_networking_secgroup" "secgroup" {
  name = var.security_group_name
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

# Create 2 ecs instances in different availability zones
resource "g42cloud_compute_instance" "ecs" {
  count = 2
  availability_zone  = data.g42cloud_availability_zones.az.names[count.index]
  security_group_ids = [g42cloud_networking_secgroup.secgroup.id]
  system_disk_size   = var.ecs_volume_size
  system_disk_type   = var.ecs_volume_type

  name       = "${var.ecs_name}_${count.index+1}"
  flavor_id  = var.ecs_flavor_id
  image_id   = data.g42cloud_images_image.image.id
  admin_pass = var.ecs_password

  network {
    uuid = data.g42cloud_vpc_subnet.subnet.id
  }
}

resource "g42cloud_rds_instance" "rds" {
  availability_zone = [
    data.g42cloud_availability_zones.az.names[0],
    data.g42cloud_availability_zones.az.names[1]
  ]

  backup_strategy {
    keep_days = 1
    start_time = "00:00-01:00"
  }

  db {
    type     = "MySQL"
    password = var.rds_password
    version  = "5.6"
    port     = "3306"
  }

  name      = var.rds_name
  flavor    = var.rds_flavor_id
  vpc_id    = data.g42cloud_vpc.vpc.id
  subnet_id = data.g42cloud_vpc_subnet.subnet.id

  ha_replication_mode = "async"
  security_group_id   = g42cloud_networking_secgroup.secgroup.id

  volume {
    size = var.rds_volume_size
    type = var.rds_volume_type
  }
}

resource "g42cloud_vpc_eip" "eip" {
  name = "${var.elb_name}-eip"
  bandwidth {
    name        = "${var.elb_name}-bandwidth"
    size        = var.eip_bandwidth_size
    charge_mode = "bandwidth"
    share_type  = "PER"
  }

  publicip {
    type = "5_bgp"
  }
}

resource "g42cloud_lb_loadbalancer" "elb" {
  name          = var.elb_name
  vip_subnet_id = data.g42cloud_vpc_subnet.subnet.subnet_id
}

resource "g42cloud_lb_listener" "listener_http" {
  protocol = "HTTP"
  protocol_port = "80"
  loadbalancer_id = g42cloud_lb_loadbalancer.elb.id
}

resource "g42cloud_lb_member" "member" {
  count = 2
  address = g42cloud_compute_instance.ecs[count.index].access_ip_v4
  pool_id = g42cloud_lb_pool.group_http.id
  protocol_port = var.elb_listen_port
  subnet_id = data.g42cloud_vpc_subnet.subnet.subnet_id
  weight = 1
}

resource "g42cloud_lb_monitor" "monitor_http" {
  delay = 5
  max_retries = 3
  pool_id = "${g42cloud_lb_pool.group_http.id}"
  timeout = 3
  type = "HTTP"
  url_path = "/"
}

resource "g42cloud_lb_pool" "group_http" {
  lb_method = "ROUND_ROBIN"
  listener_id = g42cloud_lb_listener.listener_http.id
  name = "group_http"
  protocol = "HTTP"
}

resource "g42cloud_vpc_eip_associate" "eip" {
  port_id = g42cloud_lb_loadbalancer.elb.vip_port_id
  public_ip = g42cloud_vpc_eip.eip.address
}
