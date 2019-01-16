provider "aws" {
  region = "us-east-1"
  alias = "east"
}

provider "aws" {
  region = "us-west-2"
  alias = "west"
}

variable "version" {
  default = "0.1.0"
}

resource "null_resource" "download" {
  provisioner "local-exec" {
    command = "wget https://github.com/shawncatz/automagical/releases/download/v${var.version}/automagical-${var.version}.zip"
  }
}

module "setup" {
  source = "github.com/shawncatz/automagical//terraform/setup"
}

module "east" {
  source = "github.com/shawncatz/automagical//terraform/region"
  file = "./automagical-${var.version}.zip"
  role = "${module.setup.role}"
  providers = {
    aws = "aws.east"
  }
}

module "west" {
  source = "github.com/shawncatz/automagical//terraform/region"
  file = "./automagical-${var.version}.zip"
  role = "${module.setup.role}"
  providers = {
    aws = "aws.west"
  }
}
