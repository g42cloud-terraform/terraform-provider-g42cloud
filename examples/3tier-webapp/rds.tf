# Create RDS PostgreSQL
resource "g42cloud_rds_instance" "rds" {
  availability_zone = [
    "ae-ad-1a",
    "ae-ad-1a"
  ]

  backup_strategy {
    keep_days = 1
    start_time = "00:00-01:00"
  }

  db {
    type     = "PostgreSQL"
    password = var.rds_password
    version  = "11"
    port     = "8635"
  }

  name      = var.rds_name
  flavor    = var.rds_flavor_id
  vpc_id    = data.g42cloud_vpc.vpc.id
  subnet_id = g42cloud_vpc_subnet.subnet_db.id

  ha_replication_mode = "async"
  security_group_id   = g42cloud_networking_secgroup.secgroup_db.id

  volume {
    size = var.rds_volume_size
    type = var.rds_volume_type

    disk_encryption_id = data.g42cloud_kms_key.kms.id
  }

  # RDS labels
  tags = {
    owner = "Terraform"
  }

  enterprise_project_id = var.eps_id
}
