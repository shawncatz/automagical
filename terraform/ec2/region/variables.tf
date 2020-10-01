variable "file_version" {
  description = "automagical release version"
}

variable "role" {
  description = "the role arn for the lambda"
}

variable "environment" {
  type        = "map"
  description = "pass enviornment variables to the function"
  default     = {}
}
