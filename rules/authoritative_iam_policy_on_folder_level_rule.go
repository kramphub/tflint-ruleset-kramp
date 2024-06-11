package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AuthoritativeIAMPolicyOnFolderLevelRule struct {
	tflint.DefaultRule
}

func NewAuthoritativeIAMPolicyOnFolderLevelRule() *AuthoritativeIAMPolicyOnFolderLevelRule {
	return &AuthoritativeIAMPolicyOnFolderLevelRule{}
}

func (rule *AuthoritativeIAMPolicyOnFolderLevelRule) Name() string {
	return "authoritative_iam_policy_on_folder_level"
}

func (rule *AuthoritativeIAMPolicyOnFolderLevelRule) Enabled() bool {
	return true
}

func (rule *AuthoritativeIAMPolicyOnFolderLevelRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *AuthoritativeIAMPolicyOnFolderLevelRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *AuthoritativeIAMPolicyOnFolderLevelRule) Check(runner tflint.Runner) error {
	resourceNames := []string{
		"google_folder_iam_policy",
		"google_folder_iam_binding",
		"google_folder_iam_audit_config",
	}
	if err := FindAndReportResources(runner, rule, resourceNames, DefaultMessageResourceNotAllowed); err != nil {
		return err
	}
	return nil
}
