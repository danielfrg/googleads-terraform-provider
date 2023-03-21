terraform {
  required_providers {
    googleads = {
      source = "github.com/danielfrg/googleads"
    }
  }
}

variable "customer_id" {
  type = string
}

provider "googleads" {
  customer_id = var.customer_id
}

data "googleads_coffees" "example" {}
