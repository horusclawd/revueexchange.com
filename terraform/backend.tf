terraform {
  backend "s3" {
    bucket = "revueexchange-terraform-state"
    key    = "terraform.tfstate"
    region = "us-east-1"
  }
}
