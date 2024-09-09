package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type ServiceAccountWithBasicRoleRule struct {
	tflint.DefaultRule
}

func NewServiceAccountWithBasicRoleRule() *ServiceAccountWithBasicRoleRule {
	return &ServiceAccountWithBasicRoleRule{}
}

func (rule *ServiceAccountWithBasicRoleRule) Name() string {
	return "service_account_with_basic_role"
}

func (rule *ServiceAccountWithBasicRoleRule) Enabled() bool {
	return true
}

func (rule *ServiceAccountWithBasicRoleRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (rule *ServiceAccountWithBasicRoleRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *ServiceAccountWithBasicRoleRule) Check(runner tflint.Runner) error {
	patternResourceName, err := regexp.Compile("google_.*_iam_*")
	attributesAndValuePatterns := map[string]*regexp.Regexp{
		"role":   regexp.MustCompile("roles/(viewer|editor|owner|admin)"),
		"member": regexp.MustCompile("serviceAccount:.*"),
	}
	message := "Basic role for service account found. Normally service accounts don't require a basic role. Since they grant too many permissions. Please use a specific role that matches the intended activity of the application (identified by this service account) as close as possible."
	if err != nil {
		return err
	}
	if err := FindAndReportResourcesWithAttributeHavingValue(runner, rule, *patternResourceName, attributesAndValuePatterns, message); err != nil {
		return err
	}
	return nil
}
