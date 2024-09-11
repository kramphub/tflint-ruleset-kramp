package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type ServiceAccountWithAdminRoleRule struct {
	tflint.DefaultRule
}

func NewServiceAccountWithAdminRoleRule() *ServiceAccountWithAdminRoleRule {
	return &ServiceAccountWithAdminRoleRule{}
}

func (rule *ServiceAccountWithAdminRoleRule) Name() string {
	return "service_account_with_admin_role"
}

func (rule *ServiceAccountWithAdminRoleRule) Enabled() bool {
	return true
}

func (rule *ServiceAccountWithAdminRoleRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *ServiceAccountWithAdminRoleRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *ServiceAccountWithAdminRoleRule) Check(runner tflint.Runner) error {
	patternResourceName, err := regexp.Compile("google_.*_iam_*")
	attributesAndValuePatterns := map[string]*regexp.Regexp{
		"role":   regexp.MustCompile("roles/.*\\.admin"),
		"member": regexp.MustCompile("serviceAccount:.*"),
	}
	message := "Admin role for service account found. Normally service accounts don't require an admin role. Since administration of resources is done via Terraform and not by the application itself. Please use a more specific role that matches the intended activity of the application (identified by this service account) as close as possible (e.g. `roles/pubsub.publisher` instead of `roles/pubsub.admin` for a topic)."
	if err != nil {
		return err
	}
	if err := FindAndReportResourcesWithAttributeHavingValue(runner, rule, *patternResourceName, attributesAndValuePatterns, message); err != nil {
		return err
	}
	return nil
}
