# This should pass without error... if the 'ignores' are respected by tflint (that's what being verified with this manual_tests)

# --- authoritative_iam_policy_on_folder_level ---

# tflint-ignore: authoritative_iam_policy_on_folder_level #tfsec:ignore:custom-custom-authoritative_iam_policy_on_folder_level
resource "google_folder_iam_policy" "folder" {
  folder      = "folders/1234567"
  policy_data = data.google_iam_policy.admin.policy_data
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/compute.admin"

    members = [
      "user:jane@example.com",
    ]

    condition {
      title       = "expires_after_2019_12_31"
      description = "Expiring at midnight of 2019-12-31"
      expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
    }
  }
}

# --- authoritative_iam_policy_on_project_level ---

# tflint-ignore: authoritative_iam_policy_on_project_level #tfsec:ignore:custom-custom-authoritative_iam_policy_on_project_level
resource "google_project_iam_binding" "authoritative_iam_policy_on_project_level" {
  project = "your-project-id"
  role    = "roles/container.admin"
  members = [
    "group:somegroup@example.com"
  ]
}

# --- iam_policy_on_folder_level ---

resource "google_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}

# tflint-ignore: iam_policy_on_folder_level #tfsec:ignore:custom-custom-iam_policy_on_folder_level
resource "google_folder_iam_member" "admin" {
  folder = google_folder.department1.name
  role   = "roles/editor"
  member = "user:alice@gmail.com"
}

# --- iam_policy_on_project_level ---

# tflint-ignore: iam_policy_on_project_level #tfsec:ignore:custom-custom-iam_policy_on_project_level
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

# --- authoritative_iam_policy_on_resource_level ---

# tflint-ignore: authoritative_iam_policy_on_resource_level #tfsec:ignore:custom-custom-authoritative_iam_policy_on_resource_level
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

# tflint-ignore: authoritative_iam_policy_on_resource_level #tfsec:ignore:custom-custom-authoritative_iam_policy_on_resource_level
resource "google_cloudfunctions_function_iam_policy" "authoritative_iam_policy_on_cloudfunction_resource_level" {
  project        = "your-project-id"
  region         = "europe-west1"
  cloud_function = "your-cloud-function"
  policy_data    = data.google_iam_policy.cloud_function_invoker_policy.policy_data
}

# --- creating_key_for_service_account ---

resource "google_service_account" "myaccount" {
  account_id   = "myaccount"
  display_name = "My Service Account"
}

# tflint-ignore: creating_key_for_service_account #tfsec:ignore:custom-custom-creating_key_for_service_account
resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}
