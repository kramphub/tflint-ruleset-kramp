resource "google_cloud_run_v2_service_iam_binding" "authoritative_iam_policy_on_cloud_run_resource_level" {
  project  = "your-project-id"
  location = "europe-west1"
  name     = "foo"
  role     = "roles/run.invoker"
  members = [
    "group:somegroup@example.com"
  ]
}

data "google_iam_policy" "cloud_function_invoker_policy" {
  binding {
    role = "cloudfunctions.invoker"
    members = [
      "group:somegroup@example.com"
    ]
  }
}

resource "google_cloudfunctions_function_iam_policy" "authoritative_iam_policy_on_cloudfunction_resource_level" {
  project        = "your-project-id"
  region         = "europe-west1"
  cloud_function = "your-cloud-function"
  policy_data    = data.google_iam_policy.cloud_function_invoker_policy.policy_data
}
