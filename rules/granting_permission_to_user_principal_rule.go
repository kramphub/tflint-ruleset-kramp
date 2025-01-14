package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type GrantingPermissionToUserPrincipalRule struct {
	tflint.DefaultRule
}

func NewGrantingPermissionToUserPrincipalRule() *GrantingPermissionToUserPrincipalRule {
	return &GrantingPermissionToUserPrincipalRule{}
}

func (rule *GrantingPermissionToUserPrincipalRule) Name() string {
	return "granting_permission_to_user_principal"
}

func (rule *GrantingPermissionToUserPrincipalRule) Enabled() bool {
	return true
}

func (rule *GrantingPermissionToUserPrincipalRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (rule *GrantingPermissionToUserPrincipalRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *GrantingPermissionToUserPrincipalRule) Check(runner tflint.Runner) error {
	body, err := FindResourcesThatGrantPermissionToPrincipal(runner)
	if err != nil {
		return err
	}
	attributesForPermissionAssignment := GetAttributesUsedForPermissionAssignment()
	for _, resource := range body.Blocks {
		attributes := GetAttributesForResource(resource, attributesForPermissionAssignment)
		for _, attribute := range attributes {
			err := EvaluateAttributeStringValue(attribute, runner, func(attributeValue string) error {
				if isUserPrincipal(attribute.Name, attributeValue) {
					msg := fmt.Sprintf("Granting permission to user principal '%s' in `%s` `%s`. It's strongly advised to assign permissions on group level.", attributeValue,
						GetResourceBlockType(resource), GetResourceBlockName(resource))
					return runner.EmitIssue(rule, msg, resource.DefRange)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isUserPrincipal(attributeName string, attributeValue string) bool {
	userPrincipal := strings.HasPrefix(attributeName, "user_by_") || // Check the name of the attribute
		strings.HasPrefix(attributeValue, "user:") // Check the value of the attribute
	testAddress := strings.HasSuffix(attributeValue, "@example.com") // For testing purposes
	return userPrincipal && !testAddress
}
