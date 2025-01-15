resource "google_project_iam_member" "secret-accessor-on-project-level" {
  project = "your-project-id"
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
}

# This role is allowed by the rule
resource "google_project_iam_member" "trace-agent-on-project-level" {
  project = "your-project-id"
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
}

# This role is allowed by the rule
resource "google_project_iam_member" "run-service-agent-on-project-level" {
  project = "your-project-id"
  role    = "roles/run.serviceAgent"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
}
