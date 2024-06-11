package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AuthoritativeIAMPolicyOnProjectLevelRule struct {
	tflint.DefaultRule
}

func NewAuthoritativeIAMPolicyOnProjectLevelRule() *AuthoritativeIAMPolicyOnProjectLevelRule {
	return &AuthoritativeIAMPolicyOnProjectLevelRule{}
}

func (rule *AuthoritativeIAMPolicyOnProjectLevelRule) Name() string {
	return "authoritative_iam_policy_on_project_level"
}

func (rule *AuthoritativeIAMPolicyOnProjectLevelRule) Enabled() bool {
	return true
}

func (rule *AuthoritativeIAMPolicyOnProjectLevelRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *AuthoritativeIAMPolicyOnProjectLevelRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *AuthoritativeIAMPolicyOnProjectLevelRule) Check(runner tflint.Runner) error {
	resourceNames := []string{
		"google_project_iam_policy",
		"google_project_iam_binding",
		"google_project_iam_audit_config",
	}
	if err := FindAndReportResources(runner, rule, resourceNames, DefaultMessageResourceNotAllowed); err != nil {
		return err
	}
	return nil
}
