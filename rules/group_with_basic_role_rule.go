package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type GroupWithBasicRoleRule struct {
	tflint.DefaultRule
}

func NewGroupWithBasicRoleRule() *GroupWithBasicRoleRule {
	return &GroupWithBasicRoleRule{}
}

func (rule *GroupWithBasicRoleRule) Name() string {
	return "group_account_with_basic_role"
}

func (rule *GroupWithBasicRoleRule) Enabled() bool {
	return true
}

func (rule *GroupWithBasicRoleRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *GroupWithBasicRoleRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *GroupWithBasicRoleRule) Check(runner tflint.Runner) error {
	patternResourceName, err := regexp.Compile("google_.*_iam_*")
	attributesAndValuePatterns := map[string]*regexp.Regexp{
		"role":   regexp.MustCompile("roles/(viewer|editor|owner|admin)"),
		"member": regexp.MustCompile("group:.*"),
	}
	message := "Basic role for group found. Normally groups don't require a basic role. Since they grant too many permissions. Please use a specific role that matches what the group needs to do as close as possible."
	if err != nil {
		return err
	}
	if err := FindAndReportResourcesWithAttributeHavingValue(runner, rule, *patternResourceName, attributesAndValuePatterns, message); err != nil {
		return err
	}
	return nil
}
