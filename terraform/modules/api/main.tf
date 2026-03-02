# =============================================================================
# API Module - ECS Fargate and API Gateway
# =============================================================================

# -----------------------------------------------------------------------------
# ECS Cluster
# -----------------------------------------------------------------------------

resource "aws_ecs_cluster" "main" {
  name = var.ecs_cluster_name

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-cluster"
  }
}

# ECS Task Definition (placeholder - actual task definition in later sprints)
resource "aws_ecs_task_definition" "api" {
  family                   = "${var.project_name}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = tostring(var.ecs_task_cpu)
  memory                  = tostring(var.ecs_task_memory)
  execution_role_arn       = var.ecs_execution_role_arn
  task_role_arn           = var.ecs_task_role_arn

  container_definitions = jsonencode([
    {
      name      = "api"
      image     = "${var.project_name}/api:${var.environment}"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]
      environment = [
        { name = "ENVIRONMENT", value = var.environment },
        { name = "LOG_LEVEL", value = "info" }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${var.project_name}-api"
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  tags = {
    Name = "${var.project_name}-${var.environment}-api-task"
  }
}

# ECS Service
resource "aws_ecs_service" "api" {
  name            = "${var.project_name}-api"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.ecs_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups   = var.security_group_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.api]

  tags = {
    Name = "${var.project_name}-${var.environment}-api-service"
  }
}

# -----------------------------------------------------------------------------
# Application Load Balancer
# -----------------------------------------------------------------------------

resource "aws_lb" "api" {
  name               = "${var.project_name}-${var.environment}-api"
  internal           = false
  load_balancer_type = "application"
  security_groups    = var.security_group_ids
  subnets           = var.public_subnet_ids

  enable_deletion_protection = var.environment == "prod"

  tags = {
    Name = "${var.project_name}-${var.environment}-api-alb"
  }
}

resource "aws_lb_target_group" "api" {
  name     = "${var.project_name}-${var.environment}-api"
  port     = 80
  protocol = "HTTP"
  vpc_id   = var.vpc_id

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path               = "/health"
    port               = "traffic-port"
    protocol           = "HTTP"
    timeout            = 5
    unhealthy_threshold = 2
  }
}

resource "aws_lb_listener" "api" {
  load_balancer_arn = aws_lb.api.arn
  port             = "443"
  protocol         = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = var.alb_certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

resource "aws_lb_listener" "api_http" {
  load_balancer_arn = aws_lb.api.arn
  port             = "80"
  protocol         = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# -----------------------------------------------------------------------------
# API Gateway
# -----------------------------------------------------------------------------

resource "aws_api_gatewayv2_api" "main" {
  name        = "${var.project_name}-${var.environment}-api"
  protocol_type = "HTTP"
  description = "RevUExchange API"

  cors_configuration {
    allow_origins = ["http://localhost:5173", "https://revueexchange.com"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["Authorization", "Content-Type"]
    allow_credentials = true
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-api-gateway"
  }
}

resource "aws_api_gatewayv2_stage" "main" {
  api_id = aws_api_gatewayv2_api.main.id
  name   = var.environment

  auto_deploy = true

  access_log_settings {
    destination_arn = aws_api_gatewayv2_stage_log.api.arn
    format = "$context.requestId: $context.httpMethod $context.path $context.status $context.responseLatency"
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-stage"
  }
}

resource "aws_api_gatewayv2_api_mapping" "main" {
  api_id      = aws_api_gatewayv2_api.main.id
  domain_name = aws_api_gatewayv2_domain_name.main.domain_name
  stage       = aws_api_gatewayv2_stage.main.name
}

# Custom domain for API Gateway
resource "aws_api_gatewayv2_domain_name" "main" {
  domain_name = "api.${var.domain_name}"

  domain_name_configuration {
    certificate_arn = var.acm_certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-api-domain"
  }
}

resource "aws_route53_record" "api" {
  zone_id = var.hosted_zone_id
  name    = "api.${var.domain_name}"
  type    = "A"

  alias {
    name                   = aws_api_gatewayv2_domain_name.main.domain_name_configuration[0].target_domain_name
    zone_id                = aws_api_gatewayv2_domain_name.main.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

# API Gateway integration with ALB
resource "aws_api_gatewayv2_integration" "alb" {
  api_id           = aws_api_gatewayv2_api.main.id
  integration_type = "HTTP_PROXY"
  integration_uri  = aws_lb.api.arn
  connection_type  = "VPC_LINK"
  connection_id    = aws_api_gatewayv2_vpc_link.main.id

  timeout {
    plain_request  = 30
    plain_response = 30
  }
}

resource "aws_api_gatewayv2_route" "main" {
  api_id    = aws_api_gatewayv2_api.main.id
  route_key = "$default"
  target    = "integrations/${aws_api_gatewayv2_integration.alb.id}"
}

# VPC Link for API Gateway
resource "aws_api_gatewayv2_vpc_link" "main" {
  name               = "${var.project_name}-${var.environment}-vpc-link"
  security_group_ids = var.security_group_ids
  subnet_ids        = var.private_subnet_ids

  tags = {
    Name = "${var.project_name}-${var.environment}-vpc-link"
  }
}

# CloudWatch Logs for API Gateway
resource "aws_api_gatewayv2_stage_log" "api" {
  api_id      = aws_api_gatewayv2_api.main.id
  stage_name  = aws_api_gatewayv2_stage.main.name
  log_group_name = "/aws/apigateway/${var.project_name}-${var.environment}"
}

resource "aws_cloudwatch_log_group" "api_gateway" {
  name              = "/aws/apigateway/${var.project_name}-${var.environment}"
  retention_in_days = var.log_retention_days
}
