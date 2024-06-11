package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type CreatingKeyForServiceAccountRule struct {
	tflint.DefaultRule
}

func NewCreatingKeyForServiceAccountRule() *CreatingKeyForServiceAccountRule {
	return &CreatingKeyForServiceAccountRule{}
}

func (rule *CreatingKeyForServiceAccountRule) Name() string {
	return "creating_key_for_service_account"
}

func (rule *CreatingKeyForServiceAccountRule) Enabled() bool {
	return true
}

func (rule *CreatingKeyForServiceAccountRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *CreatingKeyForServiceAccountRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *CreatingKeyForServiceAccountRule) Check(runner tflint.Runner) error {
	resourceName := "google_service_account_key"
	if err := FindAndReportResource(runner, rule, resourceName, DefaultMessageResourceNotAllowed); err != nil {
		return err
	}
	return nil
}
