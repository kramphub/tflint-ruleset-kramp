plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

plugin "google" {
  enabled = true
  version = "0.29.0"
  source  = "github.com/terraform-linters/tflint-ruleset-google"
}

# https://kramphub.atlassian.net/wiki/spaces/CPT/pages/6634111461/Terraform+custom+checks
plugin "kramp" {
  enabled = true
  version = "1.0.0"
  source  = "github.com/kramphub/tflint-ruleset-kramp"
}

config {
  call_module_type = "all"
}
