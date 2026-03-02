# =============================================================================
# Core Module Variables
# =============================================================================

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "cognito_mfa" {
  description = "Cognito MFA setting"
  type        = string
  default     = "OFF"
}

variable "enable_domain" {
  description = "Enable custom domain"
  type        = bool
  default     = false
}

variable "domain_name" {
  description = "Domain name"
  type        = string
  default     = ""
}

variable "enable_monitoring" {
  description = "Enable monitoring"
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "Log retention days"
  type        = number
  default     = 7
}

variable "alert_email" {
  description = "Alert email"
  type        = string
  default     = ""
}
