# =============================================================================
# API Module Outputs
# =============================================================================

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "ecs_cluster_arn" {
  description = "ECS cluster ARN"
  value       = aws_ecs_cluster.main.arn
}

output "ecs_task_definition_arn" {
  description = "ECS task definition ARN"
  value       = aws_ecs_task_definition.api.arn
}

output "ecs_service_name" {
  description = "ECS service name"
  value       = aws_ecs_service.api.name
}

output "api_gateway_url" {
  description = "API Gateway URL"
  value       = aws_api_gatewayv2_stage.main.invoke_url
}

output "api_gateway_arn" {
  description = "API Gateway ARN"
  value       = aws_api_gatewayv2_api.main.arn
}

output "api_gateway_id" {
  description = "API Gateway ID"
  value       = aws_api_gatewayv2_api.main.id
}

output "alb_dns_name" {
  description = "ALB DNS name"
  value       = aws_lb.api.dns_name
}

output "alb_arn" {
  description = "ALB ARN"
  value       = aws_lb.api.arn
}
