package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ServiceAccountWithBasicRoleRule(t *testing.T) {

	const hclForTest = `
resource "google_project_iam_member" "owner" {
  project = "your-project-id"
  role    = "roles/owner"
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
					Rule:    NewServiceAccountWithBasicRoleRule(),
					Message: "Issue with `google_project_iam_member` resource. Basic role for service account found. Normally service accounts don't require a basic role. Since they grant too many permissions. Please use a specific role that matches the intended activity of the application (identified by this service account) as close as possible.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 45},
					},
				},
			},
		},
	}

	rule := NewServiceAccountWithBasicRoleRule()

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
