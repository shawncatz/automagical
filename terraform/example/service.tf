data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

data "aws_availability_zones" "available" {}

resource "aws_launch_configuration" "automagical" {
  name_prefix     = "automagical"
  image_id        = "${data.aws_ami.ubuntu.image_id}"
  instance_type   = "${var.instance_type}"
  key_name        = "${var.key_name}"
  security_groups = ["${aws_security_group.automagical.id}"]
  user_data       = "${data.template_file.userdata.rendered}"

  root_block_device {
    volume_type           = "gp2"
    volume_size           = 20
    delete_on_termination = true
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_eip" "eip" {
  count = "${var.count}"
  vpc   = true

  tags {
    Name                  = "automagical-${count.index}"
    "automagical:address" = "automagical-address-${count.index}"
  }
}

resource "aws_ebs_volume" "vol" {
  count             = "${var.count}"
  size              = 10
  availability_zone = "${data.aws_availability_zones.available.names[count.index]}"

  tags = {
    Name                 = "automagical-${count.index}"
    "automagical:volume" = "automagical-volume-${count.index}"
  }
}

data "aws_route53_zone" "zone" {
  zone_id = "${var.zone_id}"
}

data "template_file" "userdata" {
  template = "${file("${path.module}/userdata.sh.tpl")}"

  vars {
    hostname = "automagical.${data.aws_route53_zone.zone.name}"
  }
}

resource "aws_autoscaling_group" "automagical" {
  count                = "${var.count}"
  name_prefix          = "automagical"
  vpc_zone_identifier  = ["${split(",", var.private_subnets)}"]
  launch_configuration = "${aws_launch_configuration.automagical.name}"
  desired_capacity     = 1
  max_size             = 1
  min_size             = 1

  tags = [
    {
      key                 = "Name"
      value               = "automagical-${count.index}"
      propagate_at_launch = true
    },
    {
      key                 = "automagical"
      value               = "true"
      propagate_at_launch = true
    },
    {
      key                 = "automagical:address"
      value               = "automagical-address-${count.index}"
      propagate_at_launch = true
    },
    {
      key                 = "automagical:volume"
      value               = "automagical-volume-${count.index}"
      propagate_at_launch = true
    },
    {
      key                 = "automagical:record"
      value               = "automagical.${data.aws_route53_zone.zone.name}"
      propagate_at_launch = true
    },
  ]
}
