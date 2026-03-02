# =============================================================================
# Events Module Outputs
# =============================================================================

output "event_bus_name" {
  description = "EventBridge bus name"
  value       = aws_cloudwatch_event_bus.main.name
}

output "event_bus_arn" {
  description = "EventBridge bus ARN"
  value       = aws_cloudwatch_event_bus.main.arn
}

output "email_queue_url" {
  description = "Email queue URL"
  value       = aws_sqs_queue.email.url
}

output "email_queue_arn" {
  description = "Email queue ARN"
  value       = aws_sqs_queue.email.arn
}

output "webhook_queue_url" {
  description = "Webhook queue URL"
  value       = aws_sqs_queue.webhook.url
}

output "webhook_queue_arn" {
  description = "Webhook queue ARN"
  value       = aws_sqs_queue.webhook.arn
}

output "export_queue_url" {
  description = "Export queue URL"
  value       = aws_sqs_queue.export.url
}

output "export_queue_arn" {
  description = "Export queue ARN"
  value       = aws_sqs_queue.export.arn
}
