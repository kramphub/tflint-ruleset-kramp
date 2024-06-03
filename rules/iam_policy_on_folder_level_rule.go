package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type IAMPolicyOnFolderLevelRule struct {
	tflint.DefaultRule
}

func NewIAMPolicyOnFolderLevelRule() *IAMPolicyOnFolderLevelRule {
	return &IAMPolicyOnFolderLevelRule{}
}

func (rule *IAMPolicyOnFolderLevelRule) Name() string {
	return "iam_policy_on_folder_level"
}

func (rule *IAMPolicyOnFolderLevelRule) Enabled() bool {
	return true
}

func (rule *IAMPolicyOnFolderLevelRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *IAMPolicyOnFolderLevelRule) Link() string {
	return TerraformGuidelinesConfluenceLink
}

func (rule *IAMPolicyOnFolderLevelRule) Check(runner tflint.Runner) error {
	resourceName := "google_folder_iam_member"
	if err := FindAndReportResource(runner, rule, resourceName, DefaultMessageResourceNotAllowed); err != nil {
		return err
	}
	return nil
}
