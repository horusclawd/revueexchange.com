# =============================================================================
# Database Module - Aurora PostgreSQL, ElastiCache, DynamoDB
# =============================================================================

# -----------------------------------------------------------------------------
# Aurora PostgreSQL
# -----------------------------------------------------------------------------

resource "aws_rds_cluster" "main" {
  cluster_identifier     = "${var.project_name}-${var.environment}-aurora"
  engine                 = "aurora-postgresql"
  engine_version         = var.db_engine_version
  database_name          = var.db_name
  master_username        = var.db_username
  master_password        = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [var.security_group_id]
  serverlessv2_scaling_configuration {
    min_capacity = var.db_allocated_capacity
    max_capacity = var.db_max_capacity
  }
  storage_encrypted   = true
  skip_final_snapshot = var.environment != "prod"
  backup_retention_period = var.db_backup_retention_days

  tags = {
    Name = "${var.project_name}-${var.environment}-aurora"
  }
}

resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-${var.environment}-subnet-group"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project_name}-${var.environment}-subnet-group"
  }
}

resource "random_password" "db_password" {
  length  = 32
  special = true
  count   = var.db_password == "" ? 1 : 0
}

locals {
  resolved_db_password = var.db_password != "" ? var.db_password : random_password.db_password[0].result
}

resource "aws_secretsmanager_secret" "database" {
  name        = "${var.project_name}/${var.environment}/aurora"
  description = "Aurora credentials"
  kms_key_id  = var.kms_key_arn

  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "database" {
  secret_id = aws_secretsmanager_secret.database.arn
  secret_string = jsonencode({
    username = var.db_username
    password = local.resolved_db_password
    engine   = "aurora-postgresql"
    host     = aws_rds_cluster.main.endpoint
    port     = 5432
    dbname   = var.db_name
  })
}

# -----------------------------------------------------------------------------
# ElastiCache Redis
# -----------------------------------------------------------------------------

resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.project_name}-${var.environment}-redis-subnet"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project_name}-${var.environment}-redis-subnet"
  }
}

resource "aws_elasticache_replication_group" "main" {
  replication_group_id       = "${var.project_name}-${var.environment}-redis"
  engine                     = "redis"
  engine_version            = "7.0"
  node_type                = var.redis_node_type
  number_cache_clusters     = var.redis_num_cache_nodes
  port                      = 6379
  parameter_group_name      = "default.redis7"
  subnet_group_name         = aws_elasticache_subnet_group.main.name
  security_group_ids        = [var.redis_security_group_id]
  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  auth_token_enabled         = true

  tags = {
    Name = "${var.project_name}-${var.environment}-redis"
  }
}

resource "random_password" "redis_auth" {
  length  = 32
  special = false
}

resource "aws_secretsmanager_secret" "redis" {
  name        = "${var.project_name}/${var.environment}/redis"
  description = "Redis auth token"
  kms_key_id  = var.kms_key_arn

  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "redis" {
  secret_id = aws_secretsmanager_secret.redis.arn
  secret_string = jsonencode({
    auth_token = random_password.redis_auth.result
    endpoint   = aws_elasticache_replication_group.main.primary_endpoint_address
    port       = 6379
  })
}

# -----------------------------------------------------------------------------
# DynamoDB Tables
# -----------------------------------------------------------------------------

resource "aws_dynamodb_table" "badges" {
  name           = "${var.project_name}-${var.environment}-badges"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "userId"
    type = "S"
  }

  global_secondary_index {
    name            = "userId-index"
    hash_key       = "userId"
    projection_type = "ALL"
  }

  server_side_encryption {
    enabled = true
  }
}

resource "aws_dynamodb_table" "leaderboard" {
  name           = "${var.project_name}-${var.environment}-leaderboard"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "rank"
    type = "N"
  }

  global_secondary_index {
    name            = "rank-index"
    hash_key       = "PK"
    range_key      = "rank"
    projection_type = "ALL"
  }

  server_side_encryption {
    enabled = true
  }
}

resource "aws_dynamodb_table" "streaks" {
  name           = "${var.project_name}-${var.environment}-streaks"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  server_side_encryption {
    enabled = true
  }
}

resource "aws_dynamodb_table" "rate_limits" {
  name           = "${var.project_name}-${var.environment}-rate-limits"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  ttl {
    attribute_name = "expiresAt"
    enabled        = true
  }

  server_side_encryption {
    enabled = true
  }
}

# -----------------------------------------------------------------------------
# CloudWatch Alarms
# -----------------------------------------------------------------------------

resource "aws_cloudwatch_metric_alarm" "aurora_cpu" {
  count = var.enable_monitoring ? 1 : 0
  alarm_name          = "${var.project_name}-${var.environment}-aurora-cpu"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "300"
  statistic           = "Average"
  threshold           = "80"
  alarm_actions       = var.alert_email != "" ? [aws_sns_topic.alerts[0].arn] : []

  dimensions = {
    DBClusterIdentifier = aws_rds_cluster.main.id
  }
}

resource "aws_sns_topic" "alerts" {
  count = var.alert_email != "" ? 1 : 0
  name = "${var.project_name}-${var.environment}-alerts"
}

resource "aws_sns_topic_subscription" "alerts_email" {
  count     = var.alert_email != "" ? 1 : 0
  topic_arn = aws_sns_topic.alerts[0].arn
  protocol  = "email"
  endpoint  = var.alert_email
}
