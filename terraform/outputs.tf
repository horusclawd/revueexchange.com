# =============================================================================
# Outputs
# =============================================================================

output "project_name" {
  description = "Project name"
  value       = var.project_name
}

output "environment" {
  description = "Environment"
  value       = var.environment
}

output "aws_region" {
  description = "AWS region"
  value       = var.aws_region
}

# VPC
output "vpc_id" {
  description = "VPC ID"
  value       = module.core.vpc_id
}

output "private_subnet_ids" {
  description = "Private subnet IDs"
  value       = module.core.private_subnet_ids
}

# Database
output "database_endpoint" {
  description = "Aurora endpoint"
  value       = module.database.aurora_endpoint
  sensitive   = true
}

output "database_name" {
  description = "Database name"
  value       = var.db_name
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = module.database.redis_endpoint
  sensitive   = true
}

# ECS
output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.api.ecs_cluster_name
}

output "ecs_cluster_arn" {
  description = "ECS cluster ARN"
  value       = module.api.ecs_cluster_arn
}

# API Gateway
output "api_gateway_url" {
  description = "API Gateway URL"
  value       = module.api.api_gateway_url
}

output "api_gateway_arn" {
  description = "API Gateway ARN"
  value       = module.api.api_gateway_arn
}

# Cognito
output "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  value       = module.core.cognito_user_pool_id
}

output "cognito_client_id" {
  description = "Cognito App Client ID"
  value       = module.core.cognito_client_id
}

# CloudFront
output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID"
  value       = module.cdn.cloudfront_distribution_id
}

output "cloudfront_domain_name" {
  description = "CloudFront domain name"
  value       = module.cdn.cloudfront_domain_name
}

output "frontend_url" {
  description = "Frontend URL"
  value       = module.cdn.cloudfront_domain_name
}

# Storage
output "uploads_bucket_name" {
  description = "Uploads bucket name"
  value       = module.storage.uploads_bucket_name
}

output "exports_bucket_name" {
  description = "Exports bucket name"
  value       = module.storage.exports_bucket_name
}

# Events
output "event_bus_name" {
  description = "EventBridge bus name"
  value       = module.events.event_bus_name
}

output "email_queue_url" {
  description = "Email queue URL"
  value       = module.events.email_queue_url
}
