package rules

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const terraformGuidelinesConfluenceLink = "https://kramphub.atlassian.net/wiki/spaces/CPT/pages/6634111461/Terraform+custom+checks"

const DefaultMessageResourceNotAllowed = "The usage of this resource is not allowed (or strongly discouraged)."

func FindAndReportResources(runner tflint.Runner, rule tflint.Rule, resourceNames []string, message string) error {
	for _, resourceName := range resourceNames {
		if err := FindAndReportResource(runner, rule, resourceName, message); err != nil {
			return err
		}
	}
	return nil
}

func FindAndReportResource(runner tflint.Runner, rule tflint.Rule, resourceName string, message string) error {
	resources, err := runner.GetResourceContent(resourceName, nil, nil)
	if err != nil {
		return err
	}
	for _, resource := range resources.Blocks {
		if err := runner.EmitIssue(rule, createMessage(resourceName, message), resource.DefRange); err != nil {
			return err
		}
	}
	return nil
}

// FindAndReportResourcesForPattern Using two regexes to match and exclude resources because Go Regex doesn't support negative lookarounds (which would be a cleaner solution; only requiring one pattern)
func FindAndReportResourcesForPattern(runner tflint.Runner, rule tflint.Rule, includedResourcesPattern regexp.Regexp, excludedResourcesPattern regexp.Regexp, message string) error {
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: nil},
		},
	}, nil)
	if err != nil {
		return err
	}
	for _, resource := range body.Blocks {
		resourceType := GetResourceBlockType(resource)

		if &excludedResourcesPattern != nil {
			logger.Debug(fmt.Sprintf("Exclusion pattern provided '%s'", excludedResourcesPattern.String()))
			if excludedResourcesPattern.MatchString(resourceType) {
				logger.Debug(fmt.Sprintf("Resource '%s' matches regex '%s'. It's excluded and will be skipped.", resourceType, excludedResourcesPattern.String()))
				continue
			} else {
				logger.Debug(fmt.Sprintf("Resource '%s' doesn't match regex '%s'. Therefore it's not excluded and not be skipped.", resourceType, excludedResourcesPattern.String()))
			}
		}

		if includedResourcesPattern.MatchString(resourceType) {
			logger.Debug(fmt.Sprintf("Resource '%s' matches regex '%s'. Reporting problem.", resourceType, includedResourcesPattern.String()))
			if err := runner.EmitIssue(rule, createMessage(resourceType, message), resource.DefRange); err != nil {
				return err
			}
		} else {
			logger.Debug(fmt.Sprintf("Resource '%s' doesn't match regex '%s'. Not reporting problem.", resourceType, includedResourcesPattern.String()))
		}
	}
	return nil
}

// FindAndReportResourcesWithAttributeHavingValue
// 1. Finds resources that have a name matching the given pattern
// 2. Then checks for those matching resources if they contain an attribute matching the given name (key in map)
// 3. And in turn checks for those attributes if they contain a value matching the given pattern (value in map)
// 4. Reports an issue if all given attribute patterns match
func FindAndReportResourcesWithAttributeHavingValue(runner tflint.Runner, rule tflint.Rule, resourcesPattern regexp.Regexp, attributeAndValuePatterns map[string]*regexp.Regexp, message string) error {

	// Filter for resources having the attributes that exist in the given map
	var attributeSchemas []hclext.AttributeSchema
	for attributeNameKey := range attributeAndValuePatterns {
		attributeSchemas = append(attributeSchemas, hclext.AttributeSchema{Name: attributeNameKey})
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: attributeSchemas,
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range body.Blocks {
		resourceType := GetResourceBlockType(resource)

		if !resourcesPattern.MatchString(resourceType) {
			logger.Debug(fmt.Sprintf("Resource '%s' doesn't match regex '%s'. Not checking its attributes.", resourceType, resourcesPattern.String()))
			continue
		}

		logger.Debug(fmt.Sprintf("Resource '%s' matches regex '%s'. Checking its attributes: %v", resourceType, resourcesPattern.String(), resource.Body.Attributes))

		matchingAttributes := 0

		for attributeName, attribute := range resource.Body.Attributes {
			attributeValuePattern, exists := attributeAndValuePatterns[attributeName]
			if !exists {
				logger.Debug(fmt.Sprintf("Attribute '%s' doesn't have a value pattern. Not checking its value.", attributeName))
				continue
			}

			logger.Debug(fmt.Sprintf("Attribute '%s' has a value pattern. Checking its value.", attributeName))

			err := runner.EvaluateExpr(attribute.Expr, func(attributeValue string) error {
				if attributeValuePattern.MatchString(attributeValue) {
					logger.Debug(fmt.Sprintf("Resource '%s' has attribute '%s' with value '%s' that matches regex '%s'", resourceType, attributeName, attributeValue, attributeValuePattern.String()))
					matchingAttributes++
				}
				return nil
			}, nil)

			if err != nil {
				return err
			}
		}

		if matchingAttributes == len(attributeAndValuePatterns) {
			if err := runner.EmitIssue(rule, fmt.Sprintf("Issue with `%s` resource. %s", resourceType, message), resource.DefRange); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetAttributesUsedForPermissionAssignment() []string {
	return []string{"member", "members", "user_by_email", "group_by_email"}
}

func FindResourcesThatGrantPermissionToPrincipal(runner tflint.Runner) (*hclext.BodyContent, error) {
	attrSchema := []hclext.AttributeSchema{}
	for _, attr := range GetAttributesUsedForPermissionAssignment() {
		attrSchema = append(attrSchema, hclext.AttributeSchema{Name: attr})
	}
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{
				Attributes: attrSchema,
			}},
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// EvaluateAttributeStringValue evaluates the value of an attribute as a string, respecting the different types of expressions (like array / tuple)
func EvaluateAttributeStringValue(attribute *hclext.Attribute, runner tflint.Runner, evaluateFn func(attributeValue string) error) error {
	switch attribute.Expr.(type) {
	case *hclsyntax.TemplateExpr:
		templateExpr := attribute.Expr.(*hclsyntax.TemplateExpr)
		return runner.EvaluateExpr(templateExpr, evaluateFn, nil)
	case *hclsyntax.TupleConsExpr:
		tupleExpr := attribute.Expr.(*hclsyntax.TupleConsExpr)
		for _, exprInTuple := range tupleExpr.Exprs {
			err := runner.EvaluateExpr(exprInTuple, evaluateFn, nil)
			if err != nil {
				return err
			}
		}
	default:
		logger.Warn(fmt.Sprintf("Unknown expression type '%s'", reflect.TypeOf(attribute.Expr)))
	}
	return nil
}

func GetAttributesForResource(resource *hclext.Block, attributeNames []string) []*hclext.Attribute {
	var applicableAttributes []*hclext.Attribute
	for _, attribute := range resource.Body.Attributes {
		for _, attrName := range attributeNames {
			if attribute.Name == attrName {
				applicableAttributes = append(applicableAttributes, attribute)
			}
		}
	}
	return applicableAttributes
}

func createMessage(resourceName string, message string) string {
	msg := fmt.Sprintf("`%s` resource found. %s", resourceName, message)
	return msg
}

func GetLinkForRule(ruleName string) string {
	if len(ruleName) == 0 {
		return terraformGuidelinesConfluenceLink
	}
	return fmt.Sprintf("%s#%s", terraformGuidelinesConfluenceLink, ruleName)
}

func GetResourceBlockType(resource *hclext.Block) string {
	if len(resource.Labels) < 1 {
		return "?"
	}
	return resource.Labels[0]
}

func GetResourceBlockName(resource *hclext.Block) string {
	if len(resource.Labels) < 2 {
		return "?"
	}
	return resource.Labels[1]
}
