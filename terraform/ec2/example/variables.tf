variable "vpc_id" {}
variable "zone_id" {}
variable "key_name" {}
variable "private_subnets" {}

variable "count" {
  default = 1
}

variable "version" {
  default = "0.1.8"
}

variable "instance_type" {
  default = "t2.micro"
}
