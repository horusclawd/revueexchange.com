# =============================================================================
# API Module Variables
# =============================================================================

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs"
  type        = list(string)
}

variable "public_subnet_ids" {
  description = "Public subnet IDs"
  type        = list(string)
}

variable "security_group_ids" {
  description = "Security group IDs"
  type        = list(string)
}

variable "kms_key_arn" {
  description = "KMS key ARN"
  type        = string
}

variable "secrets_arn" {
  description = "Secrets Manager ARN"
  type        = string
}

variable "jwt_secret_arn" {
  description = "JWT secret ARN"
  type        = string
}

variable "stripe_secret_arn" {
  description = "Stripe secret ARN"
  type        = string
}

variable "database_secret_arn" {
  description = "Database secret ARN"
  type        = string
}

variable "database_endpoint" {
  description = "Database endpoint"
  type        = string
}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "revueexchange"
}

variable "database_port" {
  description = "Database port"
  type        = number
  default     = 5432
}

variable "redis_endpoint" {
  description = "Redis endpoint"
  type        = string
}

variable "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  type        = string
}

variable "event_bus_arn" {
  description = "EventBridge bus ARN"
  type        = string
}

variable "event_bus_name" {
  description = "EventBridge bus name"
  type        = string
}

variable "email_queue_url" {
  description = "Email queue URL"
  type        = string
}

variable "webhook_queue_url" {
  description = "Webhook queue URL"
  type        = string
}

variable "export_queue_url" {
  description = "Export queue URL"
  type        = string
}

variable "uploads_bucket_name" {
  description = "Uploads bucket name"
  type        = string
}

variable "ecs_cluster_name" {
  description = "ECS cluster name"
  type        = string
}

variable "ecs_task_memory" {
  description = "ECS task memory"
  type        = number
}

variable "ecs_task_cpu" {
  description = "ECS task CPU"
  type        = number
}

variable "ecs_desired_count" {
  description = "ECS desired count"
  type        = number
}

variable "ecs_execution_role_arn" {
  description = "ECS execution role ARN"
  type        = string
}

variable "ecs_task_role_arn" {
  description = "ECS task role ARN"
  type        = string
}

variable "api_throttle_rate" {
  description = "API throttling rate"
  type        = number
}

variable "api_throttle_burst" {
  description = "API throttling burst"
  type        = number
}

variable "enable_monitoring" {
  description = "Enable monitoring"
  type        = bool
}

variable "log_retention_days" {
  description = "Log retention days"
  type        = number
}

variable "domain_name" {
  description = "Domain name"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN"
  type        = string
}

variable "hosted_zone_id" {
  description = "Route53 hosted zone ID"
  type        = string
}

variable "alb_certificate_arn" {
  description = "ALB certificate ARN"
  type        = string
}
