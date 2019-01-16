resource "aws_dynamodb_table" "lambda-dns-table" {
  name           = "${var.table_name}"
  read_capacity  = "${var.read_capacity}"
  write_capacity = "${var.write_capacity}"
  hash_key       = "InstanceID"

  attribute {
    name = "InstanceID"
    type = "S"
  }
}
