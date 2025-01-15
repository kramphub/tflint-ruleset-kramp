# Added this ignore to avoid other rule(s) to be triggered as well
#tflint-ignore: iam_policy_on_project_level
resource "google_project_iam_member" "owner" {
  project = "your-project-id"
  role    = "roles/owner"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
}
