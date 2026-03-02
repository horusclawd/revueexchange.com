# RevUExchange Terraform Infrastructure

Terraform infrastructure for RevUExchange.

## Structure

```
terraform/
├── env/
│   ├── dev/           # Dev environment
│   │   ├── main.tf
│   │   └── variables.tf
│   └── prod/         # Prod environment
│       ├── main.tf
│       └── variables.tf
├── modules/
│   ├── networking/    # VPC, subnets, security groups
│   ├── database/     # RDS PostgreSQL, ElastiCache Redis
│   ├── ecs/          # ECS Fargate, ALB
│   ├── s3/           # S3 buckets, CloudFront
│   └── route53/      # Route53, ACM certificates
├── provider.tf
├── variables.tf
└── README.md
```

## Prerequisites

- Terraform >= 1.0
- AWS CLI configured
- Existing Route53 hosted zone for `revueexchange.com`

## Setup

### 1. Create S3 bucket for Terraform state

```bash
aws s3 mb s3://revueexchange-terraform-state --region us-east-1
aws s3api put-bucket-encryption --bucket revueexchange-terraform-state --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
```

### 2. Create DynamoDB table for state locking

```bash
aws dynamodb create-table \
  --table-name revueexchange-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

### 3. Deploy

**Dev Environment:**
```bash
cd terraform/env/dev
terraform init
terraform plan
terraform apply
```

**Prod Environment:**
```bash
cd terraform/env/prod
terraform init
terraform plan
terraform apply
```

## Resource Naming

All resources follow the pattern: `{project}-{environment}-{resource}`

Examples:
- S3: `revueexchange-dev-uploads`, `revueexchange-prod-frontend`
- ECS: `revueexchange-dev-cluster`
- RDS: `revueexchange-dev-rds`
- Redis: `revueexchange-dev-redis`

## DNS

- Dev UI: `https://dev.revueexchange.com`
- Dev API: `https://dev-api.revueexchange.com`
- Prod UI: `https://revueexchange.com`
- Prod API: `https://api.revueexchange.com`

## Route53

The hosted zone is assumed to exist externally. This Terraform only:
- Looks up the existing hosted zone
- Creates ACM certificate
- Creates DNS records

The hosted zone itself should be created manually in AWS console and survive `terraform destroy`.
