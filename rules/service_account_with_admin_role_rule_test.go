package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ServiceAccountWithAdminRoleRule(t *testing.T) {

	const hclForTest = `
resource "google_pubsub_topic_iam_member" "topic_admin" {
  project = "your-project-id"
  topic   = "some-topic"
  role    = "roles/pubsub.admin"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
}

# Should not trigger an issue
resource "google_pubsub_topic_iam_member" "topic_publisher" {
  project = "your-project-id"
  topic   = "some-topic"
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:your-service@your-project-id.iam.gserviceaccount.com"
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
					Rule:    NewServiceAccountWithAdminRoleRule(),
					Message: "Issue with `google_pubsub_topic_iam_member` resource. Admin role for service account found. Normally service accounts don't require an admin role. Since administration of resources is done via Terraform and not by the application itself. Please use a more specific role that matches the intended activity of the application (identified by this service account) as close as possible (e.g. `roles/pubsub.publisher` instead of `roles/pubsub.admin` for a topic).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 56},
					},
				},
			},
		},
	}

	rule := NewServiceAccountWithAdminRoleRule()

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
