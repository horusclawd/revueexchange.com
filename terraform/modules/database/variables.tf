# =============================================================================
# Database Module Variables
# =============================================================================

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment"
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

variable "security_group_id" {
  description = "Security group ID"
  type        = string
}

variable "redis_security_group_id" {
  description = "Redis security group ID"
  type        = string
}

variable "kms_key_arn" {
  description = "KMS key ARN"
  type        = string
}

variable "db_instance_class" {
  description = "DB instance class"
  type        = string
  default     = "db.serverless"
}

variable "db_engine_version" {
  description = "DB engine version"
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

variable "db_password" {
  description = "Database password"
  type        = string
  default     = ""
}

variable "db_allocated_capacity" {
  description = "DB allocated capacity"
  type        = number
  default     = 0.5
}

variable "db_max_capacity" {
  description = "DB max capacity"
  type        = number
  default     = 4
}

variable "db_backup_retention_days" {
  description = "Backup retention days"
  type        = number
  default     = 7
}

variable "redis_node_type" {
  description = "Redis node type"
  type        = string
  default     = "cache.t3.micro"
}

variable "redis_num_cache_nodes" {
  description = "Redis cache nodes"
  type        = number
  default     = 1
}

variable "enable_monitoring" {
  description = "Enable monitoring"
  type        = bool
  default     = true
}

variable "alert_email" {
  description = "Alert email"
  type        = string
  default     = ""
}
