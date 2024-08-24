locals {
  function_name = "${replace(var.path, "/", "-")}-${var.method}"
}

resource "aws_lambda_function" "default" {
  filename         = "../build/zip/functions/${local.function_name}/package.zip"
  function_name    = local.function_name
  role             = aws_iam_role.default.arn
  handler          = "main"
  source_code_hash = filebase64sha256("../build/zip/functions/${local.function_name}/package.zip")

  architectures = ["arm64"]
  runtime       = "provided.al2023"
  memory_size   = var.memory_size
  timeout       = var.timeout

  timeouts {
    create = "3m"
  }

  environment {
    variables = var.environment_variables
  }

  logging_config {
    log_format            = "JSON"
    system_log_level      = "DEBUG"
    application_log_level = "TRACE"
  }
}

resource "aws_cloudwatch_log_group" "default" {
  name              = "/aws/lambda/${aws_lambda_function.default.function_name}"
  retention_in_days = 180
}

resource "aws_apigatewayv2_integration" "default" {
  api_id           = var.gateway_id
  integration_type = "AWS_PROXY"

  connection_type      = "INTERNET"
  description          = "Lambda example"
  integration_method   = "POST"
  integration_uri      = aws_lambda_function.default.invoke_arn
  passthrough_behavior = "WHEN_NO_MATCH"

  lifecycle {
    ignore_changes = [
      passthrough_behavior
    ]
  }
}

resource "aws_apigatewayv2_route" "default" {
  api_id    = var.gateway_id
  route_key = "${var.method} /${var.path}"
  target    = "integrations/${aws_apigatewayv2_integration.default.id}"
}

// See https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/terraform-sam-metadata.html
resource "null_resource" "sam_metadata_test" {
  triggers = {
    "resource_name"        = "aws_lambda_function.default"
    "resource_type"        = "ZIP_LAMBDA_FUNCTION"
    "original_source_code" = "../build/bin/functions/${local.function_name}"
    "built_output_path"    = "../build/zip/functions/${local.function_name}/package.zip"
  }
}

resource "aws_lambda_permission" "lambda_permission" {
  statement_id  = "Allow${aws_lambda_function.default.function_name}Invoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.default.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${var.gateway_execution_arn}/*/*/*"
}
