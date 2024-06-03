provider "google" {
  default_labels = { "github_repository" = "tflint-ruleset-kramp", "owner" = "team-cloud-platform", "terraform_managed" = "true" }
  project        = "some-project"
  region         = "europe-west1"
}
