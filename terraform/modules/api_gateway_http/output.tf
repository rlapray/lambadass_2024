output "id" {
  value = aws_apigatewayv2_api.default.id
}

output "execution_arn" {
  value = aws_apigatewayv2_api.default.execution_arn
}

output "invoke_url" {
  value = aws_apigatewayv2_stage.default.invoke_url
}
