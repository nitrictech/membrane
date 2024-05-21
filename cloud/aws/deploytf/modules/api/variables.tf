variable "name" {
  description = "The name of the API Gateway"
  type = string
}

variable "stack_id" {
  description = "The ID of the stack"
  type = string
}

variable "spec" {
  description = "Open API spec"
  type = string
}

variable "target_lambda_arns" {
  description = "The ARNs of the target lambda functions"
  type = list(string)
}

variable "domains" {
  description = "The domains to associate with the API Gateway"
  type = list(string)
}