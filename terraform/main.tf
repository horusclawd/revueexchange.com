# =============================================================================
# RevUExchange - AWS Infrastructure
# =============================================================================

locals {
  name_prefix = "${var.project_name}-${var.environment}"
  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

# -----------------------------------------------------------------------------
# Core Module (VPC, IAM, Cognito, KMS, Secrets)
# -----------------------------------------------------------------------------

module "core" {
  source = "./modules/core"

  project_name  = var.project_name
  environment   = var.environment
  aws_region    = var.aws_region

  cognito_mfa   = var.cognito_mfa
  enable_domain = var.enable_domain
  domain_name   = var.domain_name

  enable_monitoring     = var.enable_monitoring
  log_retention_days   = var.log_retention_days
  alert_email          = var.alert_email
}

# -----------------------------------------------------------------------------
# Database Module (Aurora, ElastiCache, DynamoDB)
# -----------------------------------------------------------------------------

module "database" {
  source = "./modules/database"

  project_name = var.project_name
  environment = var.environment

  vpc_id            = module.core.vpc_id
  private_subnet_ids = module.core.private_subnet_ids
  security_group_id = module.core.database_security_group_id
  redis_security_group_id = module.core.redis_security_group_id
  kms_key_arn       = module.core.kms_key_arn

  db_instance_class       = var.db_instance_class
  db_engine_version       = var.db_engine_version
  db_name                 = var.db_name
  db_username             = var.db_username
  db_allocated_capacity   = var.db_allocated_capacity
  db_max_capacity         = var.db_max_capacity
  db_backup_retention_days = var.db_backup_retention_days

  redis_node_type       = var.redis_node_type
  redis_num_cache_nodes = var.redis_num_cache_nodes

  enable_monitoring = var.enable_monitoring
  alert_email      = var.alert_email
}

# -----------------------------------------------------------------------------
# Storage Module (S3)
# -----------------------------------------------------------------------------

module "storage" {
  source = "./modules/storage"

  project_name = var.project_name
  environment  = var.environment

  kms_key_arn = module.core.kms_key_arn
}

# -----------------------------------------------------------------------------
# Events Module (EventBridge, SQS)
# -----------------------------------------------------------------------------

module "events" {
  source = "./modules/events"

  project_name = var.project_name
  environment = var.environment
}

# -----------------------------------------------------------------------------
# API Module (ECS Fargate, API Gateway)
# -----------------------------------------------------------------------------

module "api" {
  source = "./modules/api"

  project_name    = var.project_name
  environment     = var.environment
  aws_region     = var.aws_region

  vpc_id             = module.core.vpc_id
  private_subnet_ids = module.core.private_subnet_ids
  public_subnet_ids  = module.core.public_subnet_ids
  security_group_ids = [module.core.api_security_group_id]

  kms_key_arn        = module.core.kms_key_arn
  secrets_arn        = module.core.secrets_arn
  jwt_secret_arn     = module.core.jwt_secret_arn
  stripe_secret_arn  = module.core.stripe_secret_arn

  database_secret_arn = module.database.database_secret_arn
  database_endpoint   = module.database.aurora_endpoint
  database_name       = var.db_name
  database_port       = module.database.aurora_port

  redis_endpoint = module.database.redis_endpoint

  cognito_user_pool_id = module.core.cognito_user_pool_id

  event_bus_arn = module.events.event_bus_arn
  event_bus_name = module.events.event_bus_name

  email_queue_url   = module.events.email_queue_url
  webhook_queue_url = module.events.webhook_queue_url
  export_queue_url  = module.events.export_queue_url

  uploads_bucket_name = module.storage.uploads_bucket_name

  ecs_cluster_name = var.ecs_cluster_name
  ecs_task_memory  = var.ecs_task_memory
  ecs_task_cpu     = var.ecs_task_cpu
  ecs_desired_count = var.ecs_desired_count

  ecs_execution_role_arn = module.core.ecs_task_execution_role_arn
  ecs_task_role_arn     = module.core.ecs_task_role_arn

  api_throttle_rate  = var.api_throttle_rate
  api_throttle_burst = var.api_throttle_burst

  enable_monitoring   = var.enable_monitoring
  log_retention_days = var.log_retention_days

  domain_name         = var.enable_domain ? var.domain_name : ""
  acm_certificate_arn = var.enable_domain ? module.core.acm_certificate_arn : ""
  hosted_zone_id      = var.enable_domain ? module.core.hosted_zone_id : ""
  alb_certificate_arn = var.enable_domain ? module.core.acm_certificate_arn : ""
}

# -----------------------------------------------------------------------------
# CDN Module (CloudFront, S3)
# -----------------------------------------------------------------------------

module "cdn" {
  source = "./modules/cdn"

  project_name  = var.project_name
  environment   = var.environment

  domain_name   = var.domain_name
  enable_domain = var.enable_domain

  acm_certificate_arn = module.core.acm_certificate_arn
  hosted_zone_id      = module.core.hosted_zone_id
}
