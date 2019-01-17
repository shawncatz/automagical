// west provider is default
provider "aws" {
  region = "us-west-2"
}

// common setup module
module "setup" {
  // source = "github.com/shawncatz/automagical//terraform/setup"
  source = "../setup"
}

module "west" {
  // source = "github.com/shawncatz/automagical//terraform/region"
  source       = "../region"
  file_version = "${var.version}"
  role         = "${module.setup.role}"

  // you can pass in environment variables
  environment = {
    "AUTOMAGICAL_VARIABLE" = "value"
  }
}

// for each additonal region you wish to have automagical running,
// use the region module.
//provider "aws" {
//  region = "us-east-1"
//  alias = "east"
//}
//module "east" {
//  //  source = "github.com/shawncatz/automagical//terraform/region"
//  source = "../region"
//  role = "${module.setup.role}"
//  file_version = "${var.version}"
//  // Pass your providers down into the module, so the module doesn't
//  // have to worry about custom provider configurations (role
//  // assumption, etc)
//  providers = {
//    aws = "aws.east"
//  }
//}

