resource "g42cloud_lb_loadbalancer" "elb" {
  name          = var.elb_name
  vip_subnet_id = data.g42cloud_vpc_subnet.subnet.subnet_id

  enterprise_project_id = var.eps_id
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

# Create eip for the loadbalancer
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

  enterprise_project_id = var.eps_id
}

# Bind the above eip to the specified loadbalancer
resource "g42cloud_vpc_eip_associate" "eip" {
  port_id = g42cloud_lb_loadbalancer.elb.vip_port_id
  public_ip = g42cloud_vpc_eip.eip.address
}
