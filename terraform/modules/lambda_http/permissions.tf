/*****************************************************************************/
/** Role *********************************************************************/
/*****************************************************************************/
resource "aws_iam_role" "default" {
  name        = "${local.function_name}-role"
  description = "Default role for ${local.function_name} lambda"

  assume_role_policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Action : "sts:AssumeRole",
        Principal : {
          Service : "lambda.amazonaws.com"
        },
        Effect : "Allow",
        Sid : ""
      }
    ]
  })
}

resource "aws_iam_policy" "default" {
  name        = "${local.function_name}-default-policy"
  path        = "/"
  description = "Default policy for ${local.function_name} lambda"

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Sid : "Network",
        Effect : "Allow",
        Action : [
          "ec2:CreateNetworkInterface",
          "ec2:DescribeNetworkInterfaces",
          "ec2:DeleteNetworkInterface"
        ],
        Resource : "*"
      },
      {
        Sid : "Logs",
        Effect : "Allow",
        Action : [
          "logs:CreateLogStream",
          "logs:CreateLogGroup",
          "logs:PutLogEvents"
        ],
        Resource : "${aws_cloudwatch_log_group.default.arn}:*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "default" {
  role       = aws_iam_role.default.name
  policy_arn = aws_iam_policy.default.arn
}
