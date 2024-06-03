package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AuthoritativeIAMPolicyOnResourceLevelRule struct {
	tflint.DefaultRule
}

func NewAuthoritativeIAMPolicyOnResourceLevelRule() *AuthoritativeIAMPolicyOnResourceLevelRule {
	return &AuthoritativeIAMPolicyOnResourceLevelRule{}
}

func (rule *AuthoritativeIAMPolicyOnResourceLevelRule) Name() string {
	return "authoritative_iam_policy_on_resource_level"
}

func (rule *AuthoritativeIAMPolicyOnResourceLevelRule) Enabled() bool {
	return true
}

func (rule *AuthoritativeIAMPolicyOnResourceLevelRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (rule *AuthoritativeIAMPolicyOnResourceLevelRule) Link() string {
	return TerraformGuidelinesConfluenceLink
}

func (rule *AuthoritativeIAMPolicyOnResourceLevelRule) Check(runner tflint.Runner) error {
	// Not including folder and project resources, to prevent clashing with other rules
	patternIncluded, err := regexp.Compile("google_.*_iam_(policy|binding|audit_config)")
	patternExcluded, err := regexp.Compile(".*(folder|project).*")
	message := "The usage of this resource is discouraged, `_member` may be a better choice in most cases."
	if err != nil {
		return err
	}
	if err := FindAndReportResourcesForPattern(runner, rule, *patternIncluded, *patternExcluded, message); err != nil {
		return err
	}
	return nil
}
