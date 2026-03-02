# =============================================================================
# Variables
# =============================================================================

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "revueexchange"
}

variable "domain_name" {
  description = "Domain name"
  type        = string
  default     = ""
}

variable "enable_domain" {
  description = "Enable custom domain"
  type        = bool
  default     = false
}

# Database
variable "db_instance_class" {
  description = "Aurora instance class"
  type        = string
  default     = "db.serverless"
}

variable "db_engine_version" {
  description = "Aurora PostgreSQL version"
  type        = string
  default     = "15.3"
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

variable "db_allocated_capacity" {
  description = "Aurora Serverless v2 min capacity"
  type        = number
  default     = 0.5
}

variable "db_max_capacity" {
  description = "Aurora Serverless v2 max capacity"
  type        = number
  default     = 4
}

variable "db_backup_retention_days" {
  description = "Backup retention days"
  type        = number
  default     = 7
}

# Redis
variable "redis_node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.t3.micro"
}

variable "redis_num_cache_nodes" {
  description = "Number of cache nodes"
  type        = number
  default     = 1
}

# ECS
variable "ecs_cluster_name" {
  description = "ECS cluster name"
  type        = string
  default     = "revueexchange"
}

variable "ecs_task_memory" {
  description = "ECS task memory in MB"
  type        = number
  default     = 512
}

variable "ecs_task_cpu" {
  description = "ECS task CPU units"
  type        = number
  default     = 256
}

variable "ecs_desired_count" {
  description = "ECS service desired count"
  type        = number
  default     = 1
}

# API Gateway
variable "api_throttle_rate" {
  description = "API Gateway throttling rate"
  type        = number
  default     = 100
}

variable "api_throttle_burst" {
  description = "API Gateway throttling burst"
  type        = number
  default     = 200
}

# Cognito
variable "cognito_mfa" {
  description = "Cognito MFA setting"
  type        = string
  default     = "OFF"
}

# Monitoring
variable "enable_monitoring" {
  description = "Enable CloudWatch monitoring"
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "CloudWatch logs retention"
  type        = number
  default     = 7
}

variable "alert_email" {
  description = "Alert email address"
  type        = string
  default     = ""
}
