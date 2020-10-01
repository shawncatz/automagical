variable "table_name" {
  default     = "automagical_ec2"
  description = "dynamodb table name"
}

variable "read_capacity" {
  default     = 10
  description = "dynamodb read capacity (simultaneous reads)"
}

variable "write_capacity" {
  default     = 10
  description = "dynamodb write capacity (simultaneous writes)"
}
