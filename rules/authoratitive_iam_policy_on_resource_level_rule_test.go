package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AuthoritativeIAMPolicyOnResourceLevelRule(t *testing.T) {

	const hclForTest = `
resource "google_cloud_run_v2_service_iam_policy" "policy" {
  project = google_cloud_run_v2_service.default.project
  location = google_cloud_run_v2_service.default.location
  name = google_cloud_run_v2_service.default.name
  policy_data = data.google_iam_policy.admin.policy_data
}

resource "google_cloud_run_v2_service_iam_binding" "binding" {
  project = google_cloud_run_v2_service.default.project
  location = google_cloud_run_v2_service.default.location
  name = google_cloud_run_v2_service.default.name
  role = "roles/viewer"
  members = [
    "user:jane@example.com",
  ]
}

resource "google_folder_iam_member" "admin" {
  folder = google_folder.department1.name
  role   = "roles/editor"
  member = "user:alice@gmail.com"
}

resource "google_project_iam_binding" "project" {
  project = "your-project-id"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
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
					Rule:    NewAuthoritativeIAMPolicyOnResourceLevelRule(),
					Message: "`google_cloud_run_v2_service_iam_policy` resource found. The usage of this resource is discouraged, `_member` may be a better choice in most cases.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 59},
					},
				},
				{
					Rule:    NewAuthoritativeIAMPolicyOnResourceLevelRule(),
					Message: "`google_cloud_run_v2_service_iam_binding` resource found. The usage of this resource is discouraged, `_member` may be a better choice in most cases.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 1},
						End:      hcl.Pos{Line: 9, Column: 61},
					},
				},
			},
		},
	}

	rule := NewAuthoritativeIAMPolicyOnResourceLevelRule()

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
