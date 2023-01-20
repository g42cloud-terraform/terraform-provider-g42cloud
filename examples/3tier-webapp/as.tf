################################################
# Web Tier AS
################################################

# Create as configuration for web tier
resource "g42cloud_as_configuration" "as_config_web" {
  scaling_configuration_name = "tf_as_config_web"

  instance_config {
    flavor   = var.ecs_flavor_id
    image    = data.g42cloud_images_image.image.id
    key_name = g42cloud_compute_keypair.keypair.name

    disk {
      size        = 40
      volume_type = "SAS"
      disk_type   = "SYS"
    }

    disk {
      size        = 100
      volume_type = "SAS"
      disk_type   = "DATA"
      kms_id      = data.g42cloud_kms_key.kms.id
    }
  }
}

# Create as group for web tier
resource "g42cloud_as_group" "as_group_web" {
  scaling_group_name       = "tf_as_group_web"
  scaling_configuration_id = g42cloud_as_configuration.as_config_web.id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = data.g42cloud_vpc.vpc.id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = g42cloud_vpc_subnet.subnet_web.id
  }
  security_groups {
    id = g42cloud_networking_secgroup.secgroup_web.id
  }

  tags = {
    owner = "AutoScaling"
  }

  enterprise_project_id = var.eps_id
}

################################################
# App Tier AS + ELB
################################################

# Create as configuration for app tier
resource "g42cloud_as_configuration" "as_config_app" {
  scaling_configuration_name = "tf_as_config_app"

  instance_config {
    flavor   = var.ecs_flavor_id
    image    = data.g42cloud_images_image.image.id
    key_name = g42cloud_compute_keypair.keypair.name

    disk {
      size        = 40
      volume_type = "SAS"
      disk_type   = "SYS"
    }

    disk {
      size        = 100
      volume_type = "SAS"
      disk_type   = "DATA"
      kms_id      = data.g42cloud_kms_key.kms.id
    }
  }
}

# Create as group for app tier
resource "g42cloud_as_group" "as_group_app" {
  scaling_group_name       = "tf_as_group_app"
  scaling_configuration_id = g42cloud_as_configuration.as_config_app.id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = data.g42cloud_vpc.vpc.id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = g42cloud_vpc_subnet.subnet_app.id
  }
  security_groups {
    id = g42cloud_networking_secgroup.secgroup_app.id
  }
  lbaas_listeners {
    pool_id = g42cloud_elb_pool.pool.id
    protocol_port = g42cloud_elb_listener.listener.protocol_port
  }

  tags = {
    owner = "AutoScaling"
  }

  enterprise_project_id = var.eps_id
}
