# RevUExchange Terraform Infrastructure

This directory contains the Terraform infrastructure code for RevUExchange.

## Structure

```
terraform/
├── environments/
│   ├── dev/
│   │   └── dev.tfvars
│   └── prod/
│       └── prod.tfvars
├── modules/
│   ├── api/          # ECS Fargate, API Gateway, ALB
│   ├── cdn/         # CloudFront, S3 frontend
│   ├── core/        # VPC, IAM, Cognito, KMS, Secrets
│   ├── database/    # Aurora, ElastiCache, DynamoDB
│   ├── events/      # EventBridge, SQS
│   └── storage/     # S3 buckets
├── main.tf
├── variables.tf
├── outputs.tf
├── providers.tf
└── backend.tf
```

## Quick Start

### Prerequisites

- Terraform >= 1.0
- AWS CLI configured
- Docker (for local development)

### Development

```bash
# Plan dev environment
cd terraform
terraform plan -var-file=environments/dev/dev.tfvars

# Apply dev environment
terraform apply -var-file=environments/dev/dev.tfvars
```

### Production

```bash
# Plan prod environment
terraform plan -var-file=environments/prod/prod.tfvars

# Apply prod environment (with approval)
terraform apply -var-file=environments/prod/prod.tfvars
```

## Modules

### Core
VPC networking, IAM roles, Cognito User Pool, KMS encryption, Secrets Manager.

### Database
Aurora PostgreSQL (Serverless), ElastiCache Redis, DynamoDB tables (badges, leaderboard, streaks, rate-limits).

### Storage
S3 buckets for uploads, exports, and backups.

### Events
EventBridge event bus, SQS queues (email, webhooks, exports).

### API
ECS Fargate cluster, Application Load Balancer, API Gateway (HTTP API).

### CDN
CloudFront distribution, S3 bucket for frontend hosting.
