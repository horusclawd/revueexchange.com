variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "revueexchange"
}

variable "domain_name" {
  description = "Domain name (e.g., revueexchange.com)"
  type        = string
  default     = "revueexchange.com"
}

# Database
variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.serverless"
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

variable "db_allocated_storage" {
  description = "RDS allocated storage (for non-serverless)"
  type        = number
  default     = 20
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
variable "ecs_api_cpu" {
  description = "ECS API CPU units"
  type        = number
  default     = 256
}

variable "ecs_api_memory" {
  description = "ECS API memory in MB"
  type        = number
  default     = 512
}

variable "ecs_api_desired_count" {
  description = "ECS API desired count"
  type        = number
  default     = 1
}

# S3 Bucket Names
variable "documents_bucket_name" {
  description = "S3 bucket for documents"
  type        = string
  default     = ""
}

variable "logs_bucket_name" {
  description = "S3 bucket for logs"
  type        = string
  default     = ""
}

variable "ui_bucket_name" {
  description = "S3 bucket for UI"
  type        = string
  default     = ""
}

variable "enable_versioning" {
  description = "Enable S3 versioning"
  type        = bool
  default     = false
}

variable "enable_cloudfront" {
  description = "Enable CloudFront for UI"
  type        = bool
  default     = true
}

# Secrets
variable "db_password_secret_arn" {
  description = "Secret ARN for database password"
  type        = string
  default     = ""
}

variable "stripe_secret_key_arn" {
  description = "Secret ARN for Stripe"
  type        = string
  default     = ""
}

# Container Images
variable "api_container_image" {
  description = "API container image"
  type        = string
  default     = ""
}

variable "api_container_name" {
  description = "API container name"
  type        = string
  default     = "api"
}

variable "api_container_port" {
  description = "API container port"
  type        = number
  default     = 8080
}
