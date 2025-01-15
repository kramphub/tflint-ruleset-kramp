#tflint-ignore: pubsub_subscription_without_explicit_expiration
resource "google_pubsub_subscription" "with_expiration" {
  name  = "example-subscription"
  topic = "example-topic"

  labels = {
    foo = "bar"
  }

  # 20 minutes
  message_retention_duration = "1200s"
  retain_acked_messages      = true

  ack_deadline_seconds = 20

  expiration_policy {
    ttl = "300000.5s"
  }
  retry_policy {
    minimum_backoff = "10s"
  }

  enable_message_ordering    = false
}

#tflint-ignore: pubsub_subscription_without_explicit_expiration
resource "google_pubsub_subscription" "without_expiration" {
  name  = "example-subscription"
  topic = "example-topic"

  labels = {
    foo = "bar"
  }

  # 20 minutes
  message_retention_duration = "1200s"
  retain_acked_messages      = true

  ack_deadline_seconds = 20

  retry_policy {
    minimum_backoff = "10s"
  }

  enable_message_ordering    = false
}

#tflint-ignore: pubsub_subscription_without_explicit_expiration
resource "google_pubsub_subscription" "without_expiration_ttl" {
  name  = "example-subscription"
  topic = "example-topic"

  labels = {
    foo = "bar"
  }

  # 20 minutes
  message_retention_duration = "1200s"
  retain_acked_messages      = true

  ack_deadline_seconds = 20

  expiration_policy {
  }
  retry_policy {
    minimum_backoff = "10s"
  }

  enable_message_ordering    = false
}
