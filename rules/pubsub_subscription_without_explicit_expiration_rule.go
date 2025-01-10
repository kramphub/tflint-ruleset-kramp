package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type PubsubSubscriptionWithoutExplicitExpirationRule struct {
	tflint.DefaultRule
}

func NewPubsubSubscriptionWithoutExplicitExpirationRule() *PubsubSubscriptionWithoutExplicitExpirationRule {
	return &PubsubSubscriptionWithoutExplicitExpirationRule{}
}

func (rule *PubsubSubscriptionWithoutExplicitExpirationRule) Name() string {
	return "pubsub_subscription_without_explicit_expiration"
}

func (rule *PubsubSubscriptionWithoutExplicitExpirationRule) Enabled() bool {
	return true
}

func (rule *PubsubSubscriptionWithoutExplicitExpirationRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (rule *PubsubSubscriptionWithoutExplicitExpirationRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *PubsubSubscriptionWithoutExplicitExpirationRule) Check(runner tflint.Runner) error {
	resourceName := "google_pubsub_subscription"
	expirationPolicyBlockName := "expiration_policy"
	ttlAttributeName := "ttl"

	schema := &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{
				Name:     "topic",
				Required: false,
			},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: expirationPolicyBlockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name:     ttlAttributeName,
							Required: false,
						},
					},
				},
			},
		},
	}

	resources, err := runner.GetResourceContent(resourceName, schema, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		logger.Debug(fmt.Sprintf("`%s` block found, attributes: %#v", resourceName, resource.Body.Attributes))

		foundTtl := false
		for _, block := range resource.Body.Blocks {
			if block.Type == expirationPolicyBlockName {
				_, exists := block.Body.Attributes[ttlAttributeName]
				if exists {
					foundTtl = true
					break
				}
			}
		}
		if !foundTtl {
			resourceId := GetResourceBlockName(resource)
			err := runner.EmitIssue(rule, fmt.Sprintf("`%s` `%s` doesn't have an explicit `%s.%s`. Please be aware that the subscription will be deleted after 31 days of inactivity.", resourceName, resourceId, expirationPolicyBlockName, ttlAttributeName), resource.DefRange)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
