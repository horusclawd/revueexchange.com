# =============================================================================
# Core Module Outputs
# =============================================================================

output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "private_subnet_ids" {
  description = "Private subnet IDs"
  value       = [aws_subnet.private_1.id, aws_subnet.private_2.id]
}

output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = [aws_subnet.public_1.id, aws_subnet.public_2.id]
}

output "ecs_security_group_id" {
  description = "ECS security group ID"
  value       = aws_security_group.ecs.id
}

output "api_security_group_id" {
  description = "API security group ID"
  value       = aws_security_group.api.id
}

output "database_security_group_id" {
  description = "Database security group ID"
  value       = aws_security_group.database.id
}

output "redis_security_group_id" {
  description = "Redis security group ID"
  value       = aws_security_group.redis.id
}

output "kms_key_arn" {
  description = "KMS key ARN"
  value       = aws_kms_key.main.arn
}

output "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  value       = aws_cognito_user_pool.main.id
}

output "cognito_user_pool_arn" {
  description = "Cognito User Pool ARN"
  value       = aws_cognito_user_pool.main.arn
}

output "cognito_client_id" {
  description = "Cognito App Client ID"
  value       = aws_cognito_user_pool_client.main.id
}

output "cognito_user_pool_domain" {
  description = "Cognito User Pool Domain"
  value       = aws_cognito_user_pool_domain.main.domain
}

output "secrets_arn" {
  description = "Secrets Manager base ARN"
  value       = aws_secretsmanager_secret.database.arn
}

output "database_secret_arn" {
  description = "Database secret ARN"
  value       = aws_secretsmanager_secret.database.arn
}

output "jwt_secret_arn" {
  description = "JWT secret ARN"
  value       = aws_secretsmanager_secret.jwt.arn
}

output "stripe_secret_arn" {
  description = "Stripe secret ARN"
  value       = aws_secretsmanager_secret.stripe.arn
}

output "sendgrid_secret_arn" {
  description = "SendGrid secret ARN"
  value       = aws_secretsmanager_secret.sendgrid.arn
}

output "ecs_task_execution_role_arn" {
  description = "ECS task execution role ARN"
  value       = aws_iam_role.ecs_task_execution.arn
}

output "ecs_task_role_arn" {
  description = "ECS task role ARN"
  value       = aws_iam_role.ecs_task.arn
}

output "hosted_zone_id" {
  description = "Route53 hosted zone ID"
  value       = try(aws_route53_zone.main[0].zone_id, "")
}

output "acm_certificate_arn" {
  description = "ACM certificate ARN"
  value       = try(aws_acm_certificate.main[0].arn, "")
}
