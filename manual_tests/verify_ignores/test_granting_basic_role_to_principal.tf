#tflint-ignore: granting_basic_role_to_principal, iam_policy_on_project_level
resource "google_project_iam_member" "owner" {
  project = "your-project-id"
  role    = "roles/owner"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
}
