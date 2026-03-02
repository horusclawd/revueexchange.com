# =============================================================================
# Route53 Module - Look up existing hosted zone and create DNS records
# =============================================================================

variable "domain_name" {
  description = "Domain name (e.g., revueexchange.com)"
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "cloudfront_domain_name" {
  description = "CloudFront domain name"
  type        = string
  default     = ""
}

variable "alb_dns_name" {
  description = "ALB DNS name"
  type        = string
  default     = ""
}

variable "alb_zone_id" {
  description = "ALB zone ID"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags"
  type        = map(string)
  default     = {}
}

# Look up existing hosted zone
data "aws_route53_zone" "main" {
  name = var.domain_name
}

# ACM Certificate
resource "aws_acm_certificate" "main" {
  domain_name       = var.domain_name
  validation_method = "DNS"

  subject_alternative_names = [
    "*.${var.domain_name}"
  ]

  lifecycle {
    create_before_destroy = true
  }

  tags = merge(var.tags, {
    Name = "${var.domain_name}-certificate"
  })
}

# ACM Certificate validation DNS records
resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.main.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id         = data.aws_route53_zone.main.zone_id
  name            = each.value.name
  type            = each.value.type
  ttl             = 60
  allow_overwrite = true

  records = [each.value.record]
}

resource "aws_acm_certificate_validation" "main" {
  certificate_arn = aws_acm_certificate.main.arn

  validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}

# DNS record for UI subdomain to CloudFront
resource "aws_route53_record" "dev_ui" {
  count = var.cloudfront_domain_name != "" ? 1 : 0
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "${var.environment}.${var.domain_name}"
  type    = "A"

  alias {
    name                   = var.cloudfront_domain_name
    zone_id               = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}

# DNS record for API subdomain to ALB
resource "aws_route53_record" "dev_api" {
  count = var.alb_dns_name != "" ? 1 : 0
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "${var.environment}-api.${var.domain_name}"
  type    = "A"

  alias {
    name                   = var.alb_dns_name
    zone_id               = var.alb_zone_id
    evaluate_target_health = false
  }
}

# Outputs
output "nameservers" {
  value = data.aws_route53_zone.main.name_servers
}

output "hosted_zone_id" {
  value = data.aws_route53_zone.main.zone_id
}

output "certificate_arn" {
  value = aws_acm_certificate.main.arn
}
