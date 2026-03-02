# =============================================================================
# Dev Environment Variables
# =============================================================================

aws_region     = "us-east-1"
environment   = "dev"
project_name  = "revueexchange"
domain_name   = "dev.revueexchange.com"
enable_domain = false

# Database
db_instance_class       = "db.serverless"
db_engine_version       = "15.3"
db_name                 = "revueexchange"
db_username             = "revueadmin"
db_allocated_capacity   = 0.5
db_max_capacity         = 2
db_backup_retention_days = 1

# Redis
redis_node_type       = "cache.t3.micro"
redis_num_cache_nodes = 1

# ECS
ecs_cluster_name   = "revueexchange-dev"
ecs_task_memory   = 512
ecs_task_cpu      = 256
ecs_desired_count = 1

# API Gateway
api_throttle_rate = 100
api_throttle_burst = 200

# Cognito
cognito_mfa = "OFF"

# Monitoring
enable_monitoring  = false
log_retention_days = 1
