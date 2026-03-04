# Production Deployment Guide

## Prerequisites

1. **AWS Account** with appropriate permissions
2. **Domain Name** registered in Route53
3. **GitHub Repository** with GitHub Actions enabled
4. **AWS Credentials** configured for GitHub Actions

## Step 1: Configure Secrets

Store sensitive values in GitHub Secrets:

```bash
# Navigate to repository settings > Secrets and variables > Actions
# Add the following secrets:

AWS_ACCOUNT_ID=your-aws-account-id
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-iam-user-access-key
AWS_SECRET_ACCESS_KEY=your-iam-user-secret-key
DOMAIN_NAME=revueexchange.com
```

## Step 2: Create Secrets in AWS Secrets Manager

```bash
# Database password
aws secretsmanager create-secret \
  --name revueexchange/prod/db-password \
  --secret-string 'your-secure-database-password'

# Stripe API key
aws secretsmanager create-secret \
  --name revueexchange/prod/stripe-key \
  --secret-string 'sk_live_your-stripe-secret-key'

# JWT Secret (generate a secure random string)
aws secretsmanager create-secret \
  --name revueexchange/prod/jwt-secret \
  --secret-string 'your-secure-jwt-secret-at-least-32-characters'
```

## Step 3: Initialize Terraform Backend

```bash
# Create S3 bucket for state (if not exists)
aws s3 mb s3://revueexchange-terraform-state --region us-east-1

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket revueexchange-terraform-state \
  --versioning-configuration Status=Enabled

# Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name revueexchange-terraform-locks \
  --attribute-definitions AttributeName=LockID,Type=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

## Step 4: Deploy Infrastructure

```bash
cd terraform/env/prod

# Initialize Terraform
terraform init

# Plan deployment
terraform plan -var="db_password=your-db-password" \
  -var="db_password_secret_arn=arn:aws:secretsmanager:us-east-1:ACCOUNT:secret:revueexchange/prod/db-password" \
  -var="stripe_secret_key_arn=arn:aws:secretsmanager:us-east-1:ACCOUNT:secret:revueexchange/prod/stripe-key"

# Apply (will create all resources)
terraform apply -var="db_password=your-db-password" \
  -var="db_password_secret_arn=arn:aws:secretsmanager:us-east-1:ACCOUNT:secret:revueexchange/prod/db-password" \
  -var="stripe_secret_key_arn=arn:aws:secretsmanager:us-east-1:ACCOUNT:secret:revueexchange/prod/stripe-key"
```

## Step 5: Configure DNS

After Terraform deployment, update your domain's nameservers:

```bash
# Get nameservers from Route53
terraform output nameservers
```

Add these nameservers to your domain registrar.

## Step 6: First Deployment

```bash
# Push to main to trigger CI/CD
git checkout main
git merge develop
git push origin main
```

The CI pipeline will:
1. Run tests for backend and frontend
2. Build Docker images
3. Push to ECR
4. Deploy to ECS

## Monitoring

### CloudWatch Logs
```bash
# View API logs
aws logs tail /ecs/revueexchange-prod-api --follow

# View frontend logs
aws logs tail /ecs/revueexchange-prod-frontend --follow
```

### ECS Service
```bash
# Check service status
aws ecs describe-services \
  --cluster revueexchange-prod \
  --services revueexchange-backend

# View running tasks
aws ecs list-tasks \
  --cluster revueexchange-prod \
  --service-name revueexchange-backend
```

## Rollback

```bash
# Previous revision
aws ecs update-service \
  --cluster revueexchange-prod \
  --service revueexchange-backend \
  --force-new-deployment

# Or revert to previous ECS task definition
aws ecs describe-task-definition \
  --task-definition revueexchange-backend:REVISION

aws ecs update-service \
  --cluster revueexchange-prod \
  --service revueexchange-backend \
  --task-definition revueexchange-backend:OLD-REVISION
```

## Troubleshooting

### Check logs
```bash
aws logs filter-log-events \
  --log-group-name /ecs/revueexchange-prod-api \
  --filter-pattern "ERROR"
```

### Restart service
```bash
aws ecs update-service \
  --cluster revueexchange-prod \
  --service revueexchange-backend \
  --force-new-deployment
```

### SSH into container (if needed)
```bash
# Get task ID
TASK_ID=$(aws ecs list-tasks --cluster revueexchange-prod --service-name revueexchange-backend --query 'taskArns[0]' --output text | cut -d'/' -f3)

# Execute command
aws ecs execute-command \
  --cluster revueexchange-prod \
  --task $TASK_ID \
  --container api \
  --interactive \
  --command "/bin/sh"
```
