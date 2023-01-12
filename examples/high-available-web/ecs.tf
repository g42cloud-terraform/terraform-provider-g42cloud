# Fetching the availability zones of the region
data "g42cloud_availability_zones" "az" {}

# Feching the images with the specificed filters
data "g42cloud_images_image" "image" {
  name = "CentOS 7.9 64bit"
  visibility = "public"
  most_recent = true
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

  enterprise_project_id = var.eps_id
}
