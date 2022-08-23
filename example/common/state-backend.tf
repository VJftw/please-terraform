terraform {
  backend "local" {
    path = "$NAME-terraform.tfstate"
  }
}
