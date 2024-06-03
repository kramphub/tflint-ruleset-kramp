package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const TerraformGuidelinesConfluenceLink = "https://kramphub.atlassian.net/wiki/spaces/CPT/pages/6634111461/Terraform+custom+checks"

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

// Using two regexes to match and exclude resources because Go Regex doesn't support negative lookarounds (which would be a cleaner solution; only requiring one pattern)
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
		resourceName := resource.Labels[0]

		if &excludedResourcesPattern != nil {
			logger.Debug(fmt.Sprintf("Exclusion pattern provided '%s'", excludedResourcesPattern.String()))
			if excludedResourcesPattern.MatchString(resourceName) {
				logger.Debug(fmt.Sprintf("Resource '%s' matches regex '%s'. It's excluded and will be skipped.", resourceName, excludedResourcesPattern.String()))
				continue
			} else {
				logger.Debug(fmt.Sprintf("Resource '%s' doesn't match regex '%s'. Therefore it's not excluded and not be skipped.", resourceName, excludedResourcesPattern.String()))
			}
		}

		if includedResourcesPattern.MatchString(resourceName) {
			logger.Debug(fmt.Sprintf("Resource '%s' matches regex '%s'. Reporting problem.", resourceName, includedResourcesPattern.String()))
			if err := runner.EmitIssue(rule, createMessage(resourceName, message), resource.DefRange); err != nil {
				return err
			}
		} else {
			logger.Debug(fmt.Sprintf("Resource '%s' doesn't match regex '%s'. Not reporting problem.", resourceName, includedResourcesPattern.String()))
		}
	}
	return nil
}

func createMessage(resourceName string, message string) string {
	msg := fmt.Sprintf("`%s` resource found. %s", resourceName, message)
	return msg
}
