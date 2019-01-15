provider "aws" {
  region = "${var.region}"
}

resource "aws_lambda_function" "automagical" {
  filename          = "${var.file}"
  function_name     = "automagical"
  role              = "${var.role}"
  handler           = "automagical"
  runtime           = "go1.x"
  timeout           = 60
  description       = "see github.com/shawncatz/automagical"
}

resource "aws_lambda_permission" "automagical" {
  provider      = "aws.west"
  statement_id  = "LambdaDnsAllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.automagical.arn}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.cloudwatch.arn}"
  //  qualifier values don't update correctly, don't use them
}

resource "aws_cloudwatch_event_rule" "cloudwatch" {
  provider      = "aws.west"
  name          = "lambda-dns-rule-west"
  description   = "Capture EC2 instance events"
  event_pattern = "${file("${path.module}/rule.json")}"
}

resource "aws_cloudwatch_event_target" "cloudwatch" {
  provider  = "aws.west"
  rule      = "${aws_cloudwatch_event_rule.cloudwatch.name}"
  target_id = "lambda-dns-cloudwatch-event"
  arn       = "${aws_lambda_function.automagical.arn}"
}
