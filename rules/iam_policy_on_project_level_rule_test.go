package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_IAMPolicyOnProjectLevelRule(t *testing.T) {

	const hclForTest = `
resource "google_project_iam_member" "project" {
  project = "your-project-id"
  role    = "roles/editor"
  member  = "user:jane@example.com"
}

resource "google_project_iam_member" "project" {
  project = "your-project-id"
  role    = "roles/cloudtrace.agent"
  member  = "user:jane@example.com"
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
					Rule:    NewIAMPolicyOnProjectLevelRule(),
					Message: "`roles/editor` is not allowed for `google_project_iam_member`. Please grant permissions on resource level (principle of least privilege).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 47},
					},
				},
			},
		},
	}

	rule := NewIAMPolicyOnProjectLevelRule()

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
