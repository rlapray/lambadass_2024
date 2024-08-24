output "route" {
  value = "${var.method} ${var.gateway_invoke_url}/${var.path}"
}

output "method" {
  value = var.method
}

output "url" {
  value = "${var.gateway_invoke_url}/${var.path}"
}

output "logs" {
  value = "https://${var.region}.console.aws.amazon.com/cloudwatch/home?region=${var.region}#logStream:group=%252Faws%252Flambda%252F${local.function_name}"
}

output "console" {
  value = "https://${var.region}.console.aws.amazon.com/lambda/home?region=${var.region}#/functions/${local.function_name}?tab=configure"
}

output "log_group" {
  value = aws_cloudwatch_log_group.default.name
}
