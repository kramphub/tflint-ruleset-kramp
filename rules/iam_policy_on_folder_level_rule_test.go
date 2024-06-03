package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_IAMPolicyOnFolderLevelRule(t *testing.T) {

	const hclForTest = `
resource "google_folder_iam_member" "admin" {
  folder = google_folder.department1.name
  role   = "roles/editor"
  member = "user:alice@gmail.com"
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
					Rule:    NewIAMPolicyOnFolderLevelRule(),
					Message: "`google_folder_iam_member` resource found. The usage of this resource is not allowed (or strongly discouraged).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 44},
					},
				},
			},
		},
	}

	rule := NewIAMPolicyOnFolderLevelRule()

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
