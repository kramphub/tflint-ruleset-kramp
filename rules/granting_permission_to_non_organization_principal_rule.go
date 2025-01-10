package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type GrantingPermissionToNonOrganizationPrincipalRule struct {
	tflint.DefaultRule
}

func NewGrantingPermissionToNonOrganizationPrincipalRule() *GrantingPermissionToNonOrganizationPrincipalRule {
	return &GrantingPermissionToNonOrganizationPrincipalRule{}
}

func (rule *GrantingPermissionToNonOrganizationPrincipalRule) Name() string {
	return "granting_permission_to_non_organization_principal"
}

func (rule *GrantingPermissionToNonOrganizationPrincipalRule) Enabled() bool {
	return true
}

func (rule *GrantingPermissionToNonOrganizationPrincipalRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *GrantingPermissionToNonOrganizationPrincipalRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *GrantingPermissionToNonOrganizationPrincipalRule) Check(runner tflint.Runner) error {
	body, err := FindResourcesThatGrantPermissionToPrincipal(runner)
	if err != nil {
		return err
	}
	attributesForPermissionAssignment := GetAttributesUsedForPermissionAssignment()
	for _, resource := range body.Blocks {
		attributes := GetAttributesForResource(resource, attributesForPermissionAssignment)
		for _, attribute := range attributes {
			err := EvaluateAttributeStringValue(attribute, runner, func(attributeValue string) error {
				if !isOrganizationPrincipal(attributeValue) {
					msg := fmt.Sprintf("Granting permission to non-organization principal '%s' in `%s` `%s`", attributeValue,
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

func isOrganizationPrincipal(attributeValue string) bool {
	// This rule focuses on principals that represent a person or group of persons
	if strings.HasPrefix(attributeValue, "serviceAccount:") {
		return true
	}
	return strings.HasSuffix(attributeValue, "@kramp.com") ||
		strings.HasSuffix(attributeValue, "@kramphub.com") ||
		strings.HasSuffix(attributeValue, "@example.com") // For testing purposes
}
