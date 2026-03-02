variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment"
  type        = string
  default     = "prod"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "revueexchange"
}

variable "domain_name" {
  description = "Domain name"
  type        = string
  default     = "revueexchange.com"
}

# Database - larger for prod
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
  description = "RDS allocated storage"
  type        = number
  default     = 50
}

# Redis - larger cluster for prod
variable "redis_node_type" {
  description = "Redis node type"
  type        = string
  default     = "cache.t3.medium"
}

variable "redis_num_cache_nodes" {
  description = "Number of cache nodes"
  type        = number
  default     = 2
}

# ECS - larger for prod
variable "ecs_api_cpu" {
  description = "ECS API CPU units"
  type        = number
  default     = 512
}

variable "ecs_api_memory" {
  description = "ECS API memory in MB"
  type        = number
  default     = 1024
}

variable "ecs_api_desired_count" {
  description = "ECS API desired count"
  type        = number
  default     = 2
}

# S3
variable "documents_bucket_name" {
  description = "Documents bucket name"
  type        = string
  default     = ""
}

variable "logs_bucket_name" {
  description = "Logs bucket name"
  type        = string
  default     = ""
}

variable "ui_bucket_name" {
  description = "UI bucket name"
  type        = string
  default     = ""
}

variable "enable_versioning" {
  description = "Enable S3 versioning"
  type        = bool
  default     = true
}

variable "enable_cloudfront" {
  description = "Enable CloudFront"
  type        = bool
  default     = true
}

# Container
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
