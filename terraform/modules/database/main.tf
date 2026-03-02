# =============================================================================
# Database Module - RDS PostgreSQL and ElastiCache Redis
# =============================================================================

variable "environment" {
  description = "Environment"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs"
  type        = list(string)
}

variable "security_group_id" {
  description = "RDS security group ID"
  type        = string
}

variable "elasticache_security_group_id" {
  description = "ElastiCache security group ID"
  type        = string
}

variable "db_instance_class" {
  description = "DB instance class"
  type        = string
  default     = "db.serverless"
}

variable "db_allocated_storage" {
  description = "DB allocated storage"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "revueexchange"
}

variable "db_username" {
  description = "Database username"
  type        = string
  default     = "revueadmin"
}

variable "db_password" {
  description = "Database password"
  type        = string
  default     = ""
}

variable "node_type" {
  description = "Redis node type"
  type        = string
  default     = "cache.t3.micro"
}

variable "num_cache_nodes" {
  description = "Number of cache nodes"
  type        = number
  default     = 1
}

variable "redis_at_rest_encryption_enabled" {
  description = "Enable Redis at-rest encryption"
  type        = bool
  default     = true
}

variable "redis_transit_encryption_enabled" {
  description = "Enable Redis transit encryption"
  type        = bool
  default     = true
}

# RDS Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-${var.environment}-subnet-group"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project_name}-${var.environment}-subnet-group"
  }
}

# RDS Instance
resource "aws_db_instance" "main" {
  identifier     = "${var.project_name}-${var.environment}-rds"
  engine         = "postgres"
  engine_version = "15.3"
  instance_class = var.db_instance_class

  db_name  = var.db_name
  username = var.db_username
  password = local.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [var.security_group_id]

  allocated_storage     = var.db_instance_class == "db.serverless" ? null : var.dballocated_storage
  storage_encrypted     = true
  storage_type          = var.db_instance_class == "db.serverless" ? null : "gp3"

  # Serverless v2 configuration
  dynamic "serverlessv2_scaling_configuration" {
    for_each = var.db_instance_class == "db.serverless" ? [1] : []
    content {
      min_capacity = 0.5
      max_capacity = 4
    }
  }

  backup_retention_period = 7
  skip_final_snapshot     = true
  deletion_protection     = false

  tags = {
    Name = "${var.project_name}-${var.environment}-rds"
  }
}

# Random password if not provided
resource "random_password" "db_password" {
  length  = 32
  special = true
  count   = var.db_password == "" ? 1 : 0
}

locals {
  db_password = var.db_password != "" ? var.db_password : random_password.db_password[0].result
}

# ElastiCache Subnet Group
resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.project_name}-${var.environment}-redis-subnet"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project_name}-${var.environment}-redis-subnet"
  }
}

# ElastiCache Redis
resource "aws_elasticache_replication_group" "main" {
  replication_group_id       = "${var.project_name}-${var.environment}-redis"
  engine                     = "redis"
  engine_version            = "7.0"
  node_type                = var.node_type
  number_cache_clusters     = var.num_cache_nodes
  port                      = 6379
  parameter_group_name      = "default.redis7"
  subnet_group_name         = aws_elasticache_subnet_group.main.name
  security_group_ids        = [var.elasticache_security_group_id]
  at_rest_encryption_enabled = var.redis_at_rest_encryption_enabled
  transit_encryption_enabled = var.redis_transit_encryption_enabled
  auth_token_enabled         = var.redis_transit_encryption_enabled
  auto_minor_version_upgrade = true

  tags = {
    Name = "${var.project_name}-${var.environment}-redis"
  }
}

# Outputs
output "rds_endpoint" {
  value = aws_db_instance.main.endpoint
}

output "rds_port" {
  value = aws_db_instance.main.port
}

output "redis_endpoint" {
  value = aws_elasticache_replication_group.main.primary_endpoint_address
}

output "redis_port" {
  value = aws_elasticache_replication_group.main.port
}
