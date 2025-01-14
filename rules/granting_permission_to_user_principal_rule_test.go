package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_GrantingPermissionToUserPrincipalRule(t *testing.T) {

	const hclForTest = `
resource "google_project_iam_member" "project_iam_member_user" {
  project = "your-project-id"
  role    = "roles/secretmanager.secretAccessor"
  member  = "user:some-body@acme.com"
}
resource "google_project_iam_member" "project_iam_member_group" {
  project = "your-project-id"
  role    = "roles/secretmanager.secretAccessor"
  member  = "group:some-group@acme.com"
}
resource "google_project_iam_binding" "project_iam_members" {
  project = "your-project-id"
  role    = "roles/secretmanager.secretAccessor"
  members = [
    "user:some-body@acme.com",
    "group:some-group@acme.com",
    "user:some-body@example.com",
    "group:some-group@example.com"
  ]
}
resource "google_bigquery_dataset_access" "bq_dataset_access_user_by_email" {
  dataset_id    = google_bigquery_dataset.dataset.dataset_id
  role          = "OWNER"
  user_by_email = "some-body@acme.com"
}
resource "google_bigquery_dataset_iam_binding" "bq_dataset_iam_binding_members" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "roles/bigquery.dataViewer"
  members = [
    "user:some-body@acme.com",
    "group:some-group@acme.com",
    "user:some-body@example.com",
    "group:some-group@example.com"
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
					Rule:    NewGrantingPermissionToUserPrincipalRule(),
					Message: "Granting permission to user principal 'user:some-body@acme.com' in `google_project_iam_member` `project_iam_member_user`. It's strongly advised to assign permissions on group level.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 63},
					},
				},
				{
					Rule:    NewGrantingPermissionToUserPrincipalRule(),
					Message: "Granting permission to user principal 'user:some-body@acme.com' in `google_project_iam_binding` `project_iam_members`. It's strongly advised to assign permissions on group level.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 1},
						End:      hcl.Pos{Line: 12, Column: 60},
					},
				},
				{
					Rule:    NewGrantingPermissionToUserPrincipalRule(),
					Message: "Granting permission to user principal 'some-body@acme.com' in `google_bigquery_dataset_access` `bq_dataset_access_user_by_email`. It's strongly advised to assign permissions on group level.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 22, Column: 1},
						End:      hcl.Pos{Line: 22, Column: 76},
					},
				},
				{
					Rule:    NewGrantingPermissionToUserPrincipalRule(),
					Message: "Granting permission to user principal 'user:some-body@acme.com' in `google_bigquery_dataset_iam_binding` `bq_dataset_iam_binding_members`. It's strongly advised to assign permissions on group level.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 27, Column: 1},
						End:      hcl.Pos{Line: 27, Column: 80},
					},
				},
			},
		},
	}

	rule := NewGrantingPermissionToUserPrincipalRule()

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
