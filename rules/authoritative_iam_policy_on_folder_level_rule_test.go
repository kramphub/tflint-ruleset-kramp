package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AuthoritativeIAMPolicyOnFolderLevelRule(t *testing.T) {

	const hclForTest = `
resource "google_folder_iam_policy" "folder_admin_policy" {
  folder      = google_folder.department1.name
  policy_data = data.google_iam_policy.admin.policy_data
}

resource "google_folder_iam_binding" "admin" {
  folder = google_folder.department1.name
  role   = "roles/editor"

  members = [
    "user:alice@gmail.com",
  ]
}

resource "google_folder_iam_audit_config" "folder" {
  folder  = "folders/1234567"
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
					Rule:    NewAuthoritativeIAMPolicyOnFolderLevelRule(),
					Message: "`google_folder_iam_policy` resource found. The usage of this resource is not allowed (or strongly discouraged).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 58},
					},
				},
				{
					Rule:    NewAuthoritativeIAMPolicyOnFolderLevelRule(),
					Message: "`google_folder_iam_binding` resource found. The usage of this resource is not allowed (or strongly discouraged).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 45},
					},
				},
				{
					Rule:    NewAuthoritativeIAMPolicyOnFolderLevelRule(),
					Message: "`google_folder_iam_audit_config` resource found. The usage of this resource is not allowed (or strongly discouraged).",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 16, Column: 1},
						End:      hcl.Pos{Line: 16, Column: 51},
					},
				},
			},
		},
	}

	rule := NewAuthoritativeIAMPolicyOnFolderLevelRule()

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
