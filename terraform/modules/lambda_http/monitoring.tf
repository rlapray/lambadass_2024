/******************************************************************************
********** Log insight ********************************************************
******************************************************************************/
resource "aws_cloudwatch_query_definition" "default" {
  name = local.function_name

  log_group_names = [
    "${aws_cloudwatch_log_group.default.name}"
  ]

  query_string = <<EOF
fields @timestamp, log.sequence, awsRequestId, level, category, action, msg, @message
| sort @timestamp desc, log.sequence desc, awsRequestId desc
| limit 1000
EOF
}
