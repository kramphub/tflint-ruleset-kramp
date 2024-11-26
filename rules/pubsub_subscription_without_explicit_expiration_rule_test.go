package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_PubsubSubscriptionWithoutExplicitExpirationRule(t *testing.T) {

	const hclForTest = `
resource "google_pubsub_subscription" "with_expiration" {
  name  = "example-subscription"
  topic = google_pubsub_topic.example.id

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

resource "google_pubsub_subscription" "without_expiration" {
  name  = "example-subscription"
  topic = google_pubsub_topic.example.id

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

resource "google_pubsub_subscription" "without_expiration_ttl" {
  name  = "example-subscription"
  topic = google_pubsub_topic.example.id

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
}`

	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:    "issue found",
			Content: hclForTest,
			Expected: helper.Issues{
				{
					Rule:    NewPubsubSubscriptionWithoutExplicitExpirationRule(),
					Message: "`google_pubsub_subscription` `without_expiration` doesn't have an explicit `expiration_policy.ttl`. Please be aware that the subscription will be deleted after 31 days of inactivity.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 26, Column: 1},
						End:      hcl.Pos{Line: 26, Column: 59},
					},
				},
				{
					Rule:    NewPubsubSubscriptionWithoutExplicitExpirationRule(),
					Message: "`google_pubsub_subscription` `without_expiration_ttl` doesn't have an explicit `expiration_policy.ttl`. Please be aware that the subscription will be deleted after 31 days of inactivity.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 47, Column: 1},
						End:      hcl.Pos{Line: 47, Column: 63},
					},
				},
			},
		},
	}

	rule := NewPubsubSubscriptionWithoutExplicitExpirationRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
