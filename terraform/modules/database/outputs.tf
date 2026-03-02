# =============================================================================
# Database Module Outputs
# =============================================================================

output "aurora_endpoint" {
  description = "Aurora endpoint"
  value       = aws_rds_cluster.main.endpoint
  sensitive   = true
}

output "aurora_reader_endpoint" {
  description = "Aurora reader endpoint"
  value       = aws_rds_cluster.main.reader_endpoint
  sensitive   = true
}

output "aurora_port" {
  description = "Aurora port"
  value       = aws_rds_cluster.main.port
}

output "database_secret_arn" {
  description = "Database secret ARN"
  value       = aws_secretsmanager_secret.database.arn
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
  sensitive   = true
}

output "redis_port" {
  description = "Redis port"
  value       = aws_elasticache_replication_group.main.port
}

output "redis_secret_arn" {
  description = "Redis secret ARN"
  value       = aws_secretsmanager_secret.redis.arn
}

output "badges_table_name" {
  description = "Badges table name"
  value       = aws_dynamodb_table.badges.name
}

output "leaderboard_table_name" {
  description = "Leaderboard table name"
  value       = aws_dynamodb_table.leaderboard.name
}

output "streaks_table_name" {
  description = "Streaks table name"
  value       = aws_dynamodb_table.streaks.name
}

output "rate_limits_table_name" {
  description = "Rate limits table name"
  value       = aws_dynamodb_table.rate_limits.name
}
