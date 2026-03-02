# =============================================================================
# CDN Module Variables
# =============================================================================

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "domain_name" {
  description = "Domain name"
  type        = string
}

variable "enable_domain" {
  description = "Enable custom domain"
  type        = bool
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN"
  type        = string
}

variable "hosted_zone_id" {
  description = "Route53 hosted zone ID"
  type        = string
}
