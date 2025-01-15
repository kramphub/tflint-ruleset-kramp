resource "google_project_iam_binding" "authoritative_iam_policy_on_project_level" {
  project = "your-project-id"
  role    = "roles/container.admin"
  members = [
    "group:somegroup@example.com"
  ]
}
