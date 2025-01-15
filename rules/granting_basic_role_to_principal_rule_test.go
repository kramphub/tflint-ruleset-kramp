package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_GrantingBasicRoleToPrincipalRule(t *testing.T) {

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
					Rule:    NewGrantingBasicRoleToPrincipalRule(),
					Message: "Issue with `google_project_iam_member` resource. Granting basic role to principal. It's strongly discouraged to grant basic roles. Instead, grant the most limited predefined roles that meets your needs, unless there is no alternative.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 45},
					},
				},
			},
		},
	}

	rule := NewGrantingBasicRoleToPrincipalRule()

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
