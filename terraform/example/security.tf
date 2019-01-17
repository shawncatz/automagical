// simple security group

resource "aws_security_group" "automagical" {
  name   = "automagical"
  vpc_id = "${var.vpc_id}"

  tags {
    Name = "automagical"
  }
}

resource "aws_security_group_rule" "allow-ssh" {
  security_group_id = "${aws_security_group.automagical.id}"
  type              = "ingress"
  from_port         = "22"
  to_port           = "22"
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow-outgoing" {
  security_group_id = "${aws_security_group.automagical.id}"
  type              = "egress"
  from_port         = "0"
  to_port           = "0"
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
}
