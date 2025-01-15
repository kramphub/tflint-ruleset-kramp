#tflint-ignore: granting_permission_to_non_organization_principal, iam_policy_on_project_level, granting_permission_to_user_principal
resource "google_project_iam_member" "project_iam_member_user" {
  project = "your-project-id"
  role    = "roles/secretmanager.secretAccessor"
  member  = "user:some-body@acme.com"
}

#tflint-ignore: granting_permission_to_non_organization_principal, iam_policy_on_project_level
resource "google_project_iam_member" "project_iam_member_group" {
  project = "your-project-id"
  role    = "roles/secretmanager.secretAccessor"
  member  = "group:some-group@acme.com"
}

#tflint-ignore: granting_permission_to_non_organization_principal,authoritative_iam_policy_on_project_level, granting_permission_to_user_principal
resource "google_project_iam_binding" "project_iam_members" {
  project = "your-project-id"
  role    = "roles/secretmanager.secretAccessor"
  members = [
    "user:some-body@acme.com",
    "group:some-group@acme.com",
    "user:some-body@example.com",
    "group:some-group@example.com"
  ]
}

#tflint-ignore: granting_permission_to_non_organization_principal, granting_permission_to_user_principal
resource "google_bigquery_dataset_access" "bq_dataset_access_user_by_email" {
  dataset_id    = "your-project-id:your-dataset-id"
  role          = "OWNER"
  user_by_email = "some-body@acme.com"
}

#tflint-ignore: granting_permission_to_non_organization_principal, authoritative_iam_policy_on_resource_level, granting_permission_to_user_principal
resource "google_bigquery_dataset_iam_binding" "bq_dataset_iam_binding_members" {
  dataset_id = "your-project-id:your-dataset-id"
  role       = "roles/bigquery.dataViewer"
  members = [
    "user:some-body@acme.com",
    "group:some-group@acme.com",
    "user:some-body@example.com",
    "group:some-group@example.com"
  ]
}
