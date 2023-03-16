terraform {
  required_providers {
    googleads = {
      source = "github.com/danielfrg/googleads"
    }
  }
}

provider "googleads" {}

# data "googleads_coffees" "example" {}
