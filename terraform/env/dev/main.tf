# =============================================================================
# Dev Environment Terraform Configuration
# =============================================================================

terraform {
  backend "s3" {
    bucket         = "revueexchange-terraform-state"
    key            = "revueexchange/dev/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "revueexchange-terraform-locks"
  }
}

# Get current account ID
data "aws_caller_identity" "current" {}

# Networking Module
module "networking" {
  source = "../../modules/networking"

  aws_region    = var.aws_region
  environment   = var.environment
  project_name  = var.project_name
}

# Database Module
module "database" {
  source = "../../modules/database"

  environment                   = var.environment
  project_name                  = var.project_name
  private_subnet_ids            = module.networking.private_subnet_ids
  security_group_id             = module.networking.rds_security_group_id
  elasticache_security_group_id = module.networking.elasticache_security_group_id

  db_instance_class    = var.db_instance_class
  db_allocated_storage = var.db_allocated_storage
  db_name             = var.db_name
  db_username         = var.db_username
  db_password         = var.db_password
  node_type           = var.redis_node_type
  num_cache_nodes     = var.redis_num_cache_nodes
}

# Route53 Module (Certificate and DNS)
module "route53" {
  source = "../../modules/route53"

  domain_name           = var.domain_name
  environment          = var.environment
  project_name         = var.project_name
  cloudfront_domain_name = module.s3.cloudfront_domain_name
  alb_dns_name         = module.ecs.alb_dns_name
  alb_zone_id          = module.ecs.alb_zone_id
}

# S3 Module
module "s3" {
  source = "../../modules/s3"

  environment           = var.environment
  project_name          = var.project_name
  aws_region            = var.aws_region
  account_id           = data.aws_caller_identity.current.account_id

  documents_bucket_name = var.documents_bucket_name
  logs_bucket_name     = var.logs_bucket_name
  ui_bucket_name       = var.ui_bucket_name
  enable_versioning    = var.enable_versioning
  enable_cloudfront    = var.enable_cloudfront
  acm_certificate_arn = module.route53.certificate_arn
  zone_id             = module.route53.hosted_zone_id
  ui_domain           = "${var.environment}.${var.domain_name}"
}

# ECS Module
module "ecs" {
  source = "../../modules/ecs"

  environment        = var.environment
  project_name       = var.project_name
  aws_region        = var.aws_region
  account_id        = data.aws_caller_identity.current.account_id

  vpc_id            = module.networking.vpc_id
  public_subnet_ids = module.networking.public_subnet_ids
  private_subnet_ids = module.networking.private_subnet_ids

  alb_security_group_id        = module.networking.alb_security_group_id
  ecs_tasks_security_group_id = module.networking.ecs_tasks_security_group_id

  api_container_name   = var.api_container_name
  api_container_image  = var.api_container_image
  api_container_port   = var.api_container_port
  api_cpu            = var.ecs_api_cpu
  api_memory         = var.ecs_api_memory
  api_desired_count  = var.ecs_api_desired_count

  db_host               = module.database.rds_endpoint
  db_port               = module.database.rds_port
  db_name              = var.db_name
  db_password          = var.db_password
  db_password_secret_arn = var.db_password_secret_arn
  redis_host           = module.database.redis_endpoint

  stripe_secret_key_arn = var.stripe_secret_key_arn

  acm_certificate_arn = module.route53.certificate_arn
  zone_id            = module.route53.hosted_zone_id
  api_domain         = "${var.environment}-api.${var.domain_name}"
  ui_url             = "https://${var.environment}.${var.domain_name}"
}

# Outputs
output "alb_dns_name" {
  value = module.ecs.alb_dns_name
}

output "api_url" {
  value = "https://${var.environment}-api.${var.domain_name}"
}

output "database_endpoint" {
  value = module.database.rds_endpoint
}

output "redis_endpoint" {
  value = module.database.redis_endpoint
}

output "vpc_id" {
  value = module.networking.vpc_id
}

output "ui_bucket_name" {
  value = module.s3.ui_bucket_name
}

output "cloudfront_domain_name" {
  value = module.s3.cloudfront_domain_name
}

output "nameservers" {
  value = module.route53.nameservers
}

output "certificate_arn" {
  value = module.route53.certificate_arn
}
