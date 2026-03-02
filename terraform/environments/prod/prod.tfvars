# =============================================================================
# Prod Environment Variables
# =============================================================================

aws_region     = "us-east-1"
environment   = "prod"
project_name  = "revueexchange"
domain_name   = "revueexchange.com"
enable_domain = true

# Database
db_instance_class       = "db.serverless"
db_engine_version       = "15.3"
db_name                 = "revueexchange"
db_username             = "revueadmin"
db_allocated_capacity   = 1
db_max_capacity         = 8
db_backup_retention_days = 30

# Redis
redis_node_type       = "cache.t3.medium"
redis_num_cache_nodes = 2

# ECS
ecs_cluster_name   = "revueexchange-prod"
ecs_task_memory   = 1024
ecs_task_cpu      = 512
ecs_desired_count = 2

# API Gateway
api_throttle_rate = 1000
api_throttle_burst = 2000

# Cognito
cognito_mfa = "OPTIONAL"

# Monitoring
enable_monitoring  = true
log_retention_days = 30
alert_email       = ""
