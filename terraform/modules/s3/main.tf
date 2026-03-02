# =============================================================================
# S3 Module - UI Bucket, Documents Bucket, Logs Bucket
# =============================================================================

variable "environment" {
  description = "Environment"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "account_id" {
  description = "AWS account ID"
  type        = string
  default     = ""
}

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
  default     = false
}

variable "enable_cloudfront" {
  description = "Enable CloudFront"
  type        = bool
  default     = true
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN"
  type        = string
  default     = ""
}

variable "zone_id" {
  description = "Route53 zone ID"
  type        = string
  default     = ""
}

variable "ui_domain" {
  description = "UI domain"
  type        = string
  default     = ""
}

variable "alb_dns_name" {
  description = "ALB DNS name for API distribution"
  type        = string
  default     = ""
}

variable "alb_zone_id" {
  description = "ALB zone ID"
  type        = string
  default     = ""
}

locals {
  documents_bucket = var.documents_bucket_name != "" ? var.documents_bucket_name : "${var.project_name}-${var.environment}-documents"
  logs_bucket     = var.logs_bucket_name != "" ? var.logs_bucket_name : "${var.project_name}-${var.environment}-logs"
  ui_bucket       = var.ui_bucket_name != "" ? var.ui_bucket_name : "${var.project_name}-${var.environment}-ui"
}

# Documents Bucket
resource "aws_s3_bucket" "documents" {
  bucket = local.documents_bucket

  tags = {
    Name = "${var.project_name}-${var.environment}-documents"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "documents" {
  bucket = aws_s3_bucket.documents.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "documents" {
  bucket = aws_s3_bucket.documents.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Logs Bucket
resource "aws_s3_bucket" "logs" {
  bucket = local.logs_bucket

  tags = {
    Name = "${var.project_name}-${var.environment}-logs"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "logs" {
  bucket = aws_s3_bucket.logs.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "logs" {
  bucket = aws_s3_bucket.logs.id

  rule {
    id     = "expire-old-logs"
    status = "Enabled"

    expiration {
      days = 90
    }
  }
}

# UI Bucket
resource "aws_s3_bucket" "ui" {
  bucket = local.ui_bucket

  tags = {
    Name = "${var.project_name}-${var.environment}-ui"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "ui" {
  bucket = aws_s3_bucket.ui.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_website_configuration" "ui" {
  bucket = aws_s3_bucket.ui.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

resource "aws_s3_bucket_public_access_block" "ui" {
  bucket = aws_s3_bucket.ui.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "ui" {
  bucket = aws_s3_bucket.ui.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.ui.arn}/*"
      }
    ]
  })
}

# Versioning
resource "aws_s3_bucket_versioning" "documents" {
  count  = var.enable_versioning ? 1 : 0
  bucket = aws_s3_bucket.documents.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_versioning" "ui" {
  count  = var.enable_versioning ? 1 : 0
  bucket = aws_s3_bucket.ui.id

  versioning_configuration {
    status = "Enabled"
  }
}

# CloudFront Distribution
resource "aws_cloudfront_origin_access_control" "ui" {
  count = var.enable_cloudfront ? 1 : 0
  name                              = "${var.project_name}-${var.environment}-ui-oac"
  origin_access_control_origin_type = "s3"
  signing_behavior                 = "always"
  signing_protocol                 = "sigv4"
}

resource "aws_cloudfront_distribution" "ui" {
  count = var.enable_cloudfront ? 1 : 0
  enabled             = true
  comment             = "${var.project_name} UI"
  default_root_object = "index.html"

  aliases = var.ui_domain != "" ? [var.ui_domain] : []

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "ui-origin"
    viewer_protocol_policy = "redirect-to-https"
    compress               = true

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  price_class = "PriceClass_100"

  origin {
    domain_name              = aws_s3_bucket.ui.bucket_regional_domain_name
    origin_id                = "ui-origin"
    origin_access_control_id = aws_cloudfront_origin_access_control.ui[0].id
  }

  viewer_certificate {
    acm_certificate_arn      = var.acm_certificate_arn != "" ? var.acm_certificate_arn : null
    ssl_support_method       = var.acm_certificate_arn != "" ? "sni-only" : null
    minimum_protocol_version = var.acm_certificate_arn != "" ? "TLSv1.2_2021" : "TLSv1"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-cloudfront"
  }
}

# Route53 DNS Records
resource "aws_route53_record" "ui" {
  count   = var.enable_cloudfront && var.ui_domain != "" ? 1 : 0
  zone_id = var.zone_id
  name    = var.ui_domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.ui[0].domain_name
    zone_id               = aws_cloudfront_distribution.ui[0].hosted_zone_id
    evaluate_target_health = false
  }
}

# Outputs
output "documents_bucket_name" {
  value = aws_s3_bucket.documents.id
}

output "documents_bucket_arn" {
  value = aws_s3_bucket.documents.arn
}

output "logs_bucket_name" {
  value = aws_s3_bucket.logs.id
}

output "ui_bucket_name" {
  value = aws_s3_bucket.ui.id
}

output "ui_bucket_arn" {
  value = aws_s3_bucket.ui.arn
}

output "cloudfront_domain_name" {
  value = var.enable_cloudfront ? aws_cloudfront_distribution.ui[0].domain_name : ""
}

output "cloudfront_distribution_id" {
  value = var.enable_cloudfront ? aws_cloudfront_distribution.ui[0].id : ""
}
