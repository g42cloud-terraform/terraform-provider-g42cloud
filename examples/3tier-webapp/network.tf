data "g42cloud_vpc" "vpc" {
  name = "vpc-default"
}

# Creates Subnet for web tier
resource "g42cloud_vpc_subnet" "subnet_web" {
  name       = "tf-subnet-web"
  cidr       = "192.168.150.0/24"
  gateway_ip = "192.168.150.1"
  vpc_id     = data.g42cloud_vpc.vpc.id
}

# Creates Subnet for app tier
resource "g42cloud_vpc_subnet" "subnet_app" {
  name       = "tf-subnet-app"
  cidr       = "192.168.155.0/24"
  gateway_ip = "192.168.155.1"
  vpc_id     = data.g42cloud_vpc.vpc.id
}

# Creates Subnet for db tier
resource "g42cloud_vpc_subnet" "subnet_db" {
  name       = "tf-subnet-app"
  cidr       = "192.168.156.0/24"
  gateway_ip = "192.168.156.1"
  vpc_id     = data.g42cloud_vpc.vpc.id
}
