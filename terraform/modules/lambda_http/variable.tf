variable "gateway_id" {
  type = string
}

variable "gateway_execution_arn" {
  type = string
}

variable "gateway_invoke_url" {
  type = string
}

variable "path" {
  type = string
}

variable "method" {
  type = string
}

variable "timeout" {
  type    = number
  default = 5
}

variable "memory_size" {
  type    = number
  default = 128
}

variable "region" {
  type = string
}

variable "environment_variables" {
  type = map(string)
}
