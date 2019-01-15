resource "aws_iam_role" "lambda-dns-role" {
  name               = "lambda-dns-role"
  assume_role_policy = "${data.aws_iam_policy_document.role.json}"
}

resource "aws_iam_policy" "lambda-dns-policy" {
  name   = "lambda-dns-policy"
  policy = "${data.aws_iam_policy_document.policy.json}"
}

resource "aws_iam_policy_attachment" "lambda-dns-attach" {
  name       = "lambda-dns-attach"
  roles      = ["${aws_iam_role.lambda-dns-role.name}"]
  policy_arn = "${aws_iam_policy.lambda-dns-policy.arn}"
}

data "aws_iam_policy_document" "policy" {
  statement {
    effect = "Allow"
    actions = [
      "ec2:Describe*",
      "ec2:CreateNetworkInterface",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DeleteNetworkInterface",
      "route53:*",
      "dynamodb:*",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
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
