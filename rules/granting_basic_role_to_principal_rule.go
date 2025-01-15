package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type GrantingBasicRoleToPrincipalRule struct {
	tflint.DefaultRule
}

func NewGrantingBasicRoleToPrincipalRule() *GrantingBasicRoleToPrincipalRule {
	return &GrantingBasicRoleToPrincipalRule{}
}

func (rule *GrantingBasicRoleToPrincipalRule) Name() string {
	return "granting_basic_role_to_principal"
}

func (rule *GrantingBasicRoleToPrincipalRule) Enabled() bool {
	return true
}

func (rule *GrantingBasicRoleToPrincipalRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *GrantingBasicRoleToPrincipalRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *GrantingBasicRoleToPrincipalRule) Check(runner tflint.Runner) error {
	// Any `google_` resource that has a `role` attribute

	patternResourceName, err := regexp.Compile("google_.*")
	attributesAndValuePatterns := map[string]*regexp.Regexp{
		"role": regexp.MustCompile("roles/(viewer|editor|owner|admin)"),
	}

	message := "Granting basic role to principal. It's strongly discouraged to grant basic roles. Instead, grant the most limited predefined roles that meets your needs, unless there is no alternative."
	if err != nil {
		return err
	}
	if err := FindAndReportResourcesWithAttributeHavingValue(runner, rule, *patternResourceName, attributesAndValuePatterns, message); err != nil {
		return err
	}
	return nil
}
