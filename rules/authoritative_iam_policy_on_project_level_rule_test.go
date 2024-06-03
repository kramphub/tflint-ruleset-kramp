package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AuthoritativeIAMPolicyOnProjectLevelRule(t *testing.T) {

	const hclForTest = `
resource "google_project_iam_policy" "project" {
  project     = "your-project-id"
  policy_data = data.google_iam_policy.admin.policy_data
}

resource "google_project_iam_binding" "project" {
  project = "your-project-id"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}

resource "google_project_iam_audit_config" "project" {
  project = "your-project-id"
  service = "allServices"
  audit_log_config {
    log_type = "ADMIN_READ"
  }
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:joebloggs@hashicorp.com",
    ]
  }
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
					Rule:    NewAuthoritativeIAMPolicyOnProjectLevelRule(),
					Message: "`google_project_iam_policy` resource found. The usage of this resource is not allowed (or strongly discouraged).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 47},
					},
				},
				{
					Rule:    NewAuthoritativeIAMPolicyOnProjectLevelRule(),
					Message: "`google_project_iam_binding` resource found. The usage of this resource is not allowed (or strongly discouraged).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 48},
					},
				},
				{
					Rule:    NewAuthoritativeIAMPolicyOnProjectLevelRule(),
					Message: "`google_project_iam_audit_config` resource found. The usage of this resource is not allowed (or strongly discouraged).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 16, Column: 1},
						End:      hcl.Pos{Line: 16, Column: 53},
					},
				},
			},
		},
	}

	rule := NewAuthoritativeIAMPolicyOnProjectLevelRule()

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
