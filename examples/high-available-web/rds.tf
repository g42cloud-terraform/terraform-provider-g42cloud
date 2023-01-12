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
