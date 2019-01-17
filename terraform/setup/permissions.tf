resource "aws_iam_role" "role" {
  name               = "automagical-ec2-role"
  assume_role_policy = "${data.aws_iam_policy_document.role.json}"
}

resource "aws_iam_policy" "policy" {
  name   = "automagical-ec2-policy"
  policy = "${data.aws_iam_policy_document.policy.json}"
}

resource "aws_iam_policy_attachment" "attach" {
  name       = "automagical-ec2-attach"
  roles      = ["${aws_iam_role.role.name}"]
  policy_arn = "${aws_iam_policy.policy.arn}"
}

data "aws_iam_policy_document" "policy" {
  statement {
    effect = "Allow"

    actions = [
      "ec2:*",
      "route53:*",
      "dynamodb:*",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["*"]
  }
}

data "aws_iam_policy_document" "role" {
  statement {
    effect = "Allow"

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }

    actions = ["sts:AssumeRole"]
  }
}
