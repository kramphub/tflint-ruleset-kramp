package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_GroupWithBasicRoleRule(t *testing.T) {

	const hclForTest = `
resource "google_project_iam_member" "owner" {
  project = "your-project-id"
  role    = "roles/owner"
  member  = "group:somegroup@example.com"
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
					Rule:    NewGroupWithBasicRoleRule(),
					Message: "Issue with `google_project_iam_member` resource. Basic role for group found. Normally groups don't require a basic role. Since they grant too many permissions. Please use a specific role that matches what the group needs to do as close as possible.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 45},
					},
				},
			},
		},
	}

	rule := NewGroupWithBasicRoleRule()

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
