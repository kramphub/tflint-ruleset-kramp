#tflint-ignore: service_account_with_admin_role
resource "google_pubsub_topic_iam_member" "topic_admin" {
  project = "your-project-id"
  topic   = "some-topic"
  role    = "roles/pubsub.admin"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
}
