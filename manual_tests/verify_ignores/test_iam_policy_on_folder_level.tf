#tflint-ignore: iam_policy_on_folder_level,
resource "google_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}

#tflint-ignore: iam_policy_on_folder_level, granting_basic_role_to_principal
resource "google_folder_iam_member" "admin" {
  folder = google_folder.department1.name
  role   = "roles/editor"
  member = "user:alice@example.com"
}
