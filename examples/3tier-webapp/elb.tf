data "g42cloud_elb_flavors" "l7flavors" {
  type            = "L7"
  max_connections = 200000
  cps             = 2000
  bandwidth       = 50
}

resource "g42cloud_elb_loadbalancer" "elb_web" {
  name            = var.elbweb_name
  ipv4_subnet_id  = g42cloud_vpc_subnet.subnet_web.subnet_id
  l7_flavor_id    = data.g42cloud_elb_flavors.l7flavors.ids[0]
  ipv4_eip_id     = g42cloud_vpc_eip.eip.id

  availability_zone = [
    "ae-ad-1a",
    "ae-ad-1b"
  ]

  enterprise_project_id = var.eps_id
}

resource "g42cloud_elb_listener" "listener_web" {
  name            = var.elbweb_name
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_elb_loadbalancer.elb_web.id

  forward_eip = true
}

resource "g42cloud_elb_pool" "pool_web" {
  name        = var.elbweb_name
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = g42cloud_elb_listener.listener_web.id
}

# Create eip for the loadbalancer
resource "g42cloud_vpc_eip" "eip" {
  name = "${var.elbweb_name}-eip"
  bandwidth {
    name        = "${var.elbweb_name}-bandwidth"
    size        = var.eip_bandwidth_size
    charge_mode = "bandwidth"
    share_type  = "PER"
  }

  publicip {
    type = "5_bgp"
  }

  enterprise_project_id = var.eps_id
}

data "g42cloud_elb_flavors" "l4flavors" {
  type            = "L4"
  max_connections = 1000000
  cps             = 20000
  bandwidth       = 100
}

resource "g42cloud_elb_loadbalancer" "elb" {
  name            = var.elbapp_name
  ipv4_subnet_id  = g42cloud_vpc_subnet.subnet_app.subnet_id
  l4_flavor_id    = data.g42cloud_elb_flavors.l4flavors.ids[0]
  ipv4_eip_id     = g42cloud_vpc_eip.eip.id

  availability_zone = [
    "ae-ad-1a",
    "ae-ad-1b"
  ]

  enterprise_project_id = var.eps_id
}

resource "g42cloud_elb_listener" "listener" {
  name            = var.elbapp_name
  protocol        = "TCP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_elb_loadbalancer.elb.id

  forward_eip = true
}

resource "g42cloud_elb_pool" "pool_app" {
  name        = var.elbapp_name
  protocol    = "TCP"
  lb_method   = "ROUND_ROBIN"
  listener_id = g42cloud_elb_listener.listener.id

}
