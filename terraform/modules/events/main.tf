# =============================================================================
# Events Module - EventBridge and SQS
# =============================================================================

# EventBridge Event Bus
resource "aws_cloudwatch_event_bus" "main" {
  name = "${var.project_name}-${var.environment}-events"

  tags = {
    Name = "${var.project_name}-${var.environment}-event-bus"
  }
}

# Dead Letter Queue
resource "aws_sqs_queue" "event_dlq" {
  name = "${var.project_name}-${var.environment}-event-dlq"

  message_retention_seconds = 1209600

  tags = {
    Name = "${var.project_name}-${var.environment}-event-dlq"
  }
}

# Email Queue
resource "aws_sqs_queue" "email" {
  name = "${var.project_name}-${var.environment}-email-queue"

  message_retention_seconds = 86400
  receive_wait_time_seconds = 20
  visibility_timeout_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.event_dlq.arn
    maxReceiveCount     = 5
  })

  tags = {
    Name = "${var.project_name}-${var.environment}-email-queue"
  }
}

# Webhook Queue
resource "aws_sqs_queue" "webhook" {
  name = "${var.project_name}-${var.environment}-webhook-queue"

  message_retention_seconds = 86400
  receive_wait_time_seconds = 20
  visibility_timeout_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.event_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${var.project_name}-${var.environment}-webhook-queue"
  }
}

# Export Queue
resource "aws_sqs_queue" "export" {
  name = "${var.project_name}-${var.environment}-export-queue"

  message_retention_seconds = 604800
  receive_wait_time_seconds = 20
  visibility_timeout_seconds = 600

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.event_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Name = "${var.project_name}-${var.environment}-export-queue"
  }
}
