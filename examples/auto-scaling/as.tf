# Fetching the vpc with name vpc-default
data "g42cloud_vpc" "vpc" {
  name = "vpc-default"
}

# Fetching the subnet with name subnet-default
data "g42cloud_vpc_subnet" "subnet" {
  name   = "subnet-default"
}

# Fetching the availability zones of the region
data "g42cloud_availability_zones" "az" {}

# Feching the images with the specificed filters
data "g42cloud_images_image" "image" {
  name = "CentOS 7.9 64bit"
  visibility = "public"
  most_recent = true
}

# Create keypair
resource "g42cloud_compute_keypair" "keypair" {
  name       = "tf_key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB"
}

# Create kms for encrypting disk
#resource "g42cloud_kms_key" "kms" {
#  key_alias    = "tf_kms"
#  pending_days = "7"
#
#  enterprise_project_id = var.eps_id
#}

data "g42cloud_kms_key" "kms" {
  key_alias = "tf_kms"
}

# Create as configuration
resource "g42cloud_as_configuration" "as_config" {
  scaling_configuration_name = "tf_as_config"

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

# Create as group
resource "g42cloud_as_group" "as_group" {
  scaling_group_name       = "tf_as_group"
  scaling_configuration_id = g42cloud_as_configuration.as_config.id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10
  vpc_id                   = data.g42cloud_vpc.vpc.id
  delete_publicip          = true
  delete_instances         = "yes"

  networks {
    id = data.g42cloud_vpc_subnet.subnet.id
  }
  security_groups {
    id = g42cloud_networking_secgroup.secgroup.id
  }

  tags = {
    owner = "AutoScaling"
  }

  enterprise_project_id = var.eps_id
}
