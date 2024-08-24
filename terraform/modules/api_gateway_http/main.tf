resource "aws_apigatewayv2_api" "default" {
  name          = var.project
  protocol_type = "HTTP"

  /*
  cors_configuration {
    allow_credentials = true
    allow_headers = ["Content-Type", "Accept", "Location", "Authorization", "Cache-Control"]
    allow_methods = ["GET", "POST", "DELETE", "PATCH"]
    allow_origins = concat(["https://${var.domain.address}"], var.additional_origins)
    max_age = 60
  }
  */
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.default.id
  name        = "default"
  auto_deploy = true
  //stage_variables = var.environment_variables

  lifecycle {
    ignore_changes = [deployment_id]
  }
}
