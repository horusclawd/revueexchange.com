# =============================================================================
# Storage Module Outputs
# =============================================================================

output "uploads_bucket_name" {
  description = "Uploads bucket name"
  value       = aws_s3_bucket.uploads.id
}

output "uploads_bucket_arn" {
  description = "Uploads bucket ARN"
  value       = aws_s3_bucket.uploads.arn
}

output "exports_bucket_name" {
  description = "Exports bucket name"
  value       = aws_s3_bucket.exports.id
}

output "exports_bucket_arn" {
  description = "Exports bucket ARN"
  value       = aws_s3_bucket.exports.arn
}

output "backups_bucket_name" {
  description = "Backups bucket name"
  value       = aws_s3_bucket.backups.id
}

output "backups_bucket_arn" {
  description = "Backups bucket ARN"
  value       = aws_s3_bucket.backups.arn
}
