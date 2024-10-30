package rules

import (
	"fmt"
	"slices"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type IAMPolicyOnProjectLevelRule struct {
	tflint.DefaultRule
}

func NewIAMPolicyOnProjectLevelRule() *IAMPolicyOnProjectLevelRule {
	return &IAMPolicyOnProjectLevelRule{}
}

func (rule *IAMPolicyOnProjectLevelRule) Name() string {
	return "iam_policy_on_project_level"
}

func (rule *IAMPolicyOnProjectLevelRule) Enabled() bool {
	return true
}

func (rule *IAMPolicyOnProjectLevelRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *IAMPolicyOnProjectLevelRule) Link() string {
	return GetLinkForRule(rule.Name())
}

func (rule *IAMPolicyOnProjectLevelRule) Check(runner tflint.Runner) error {
	resourceName := "google_project_iam_member"
	allowedRoles := []string{
		"roles/cloudtrace.agent",                                       // https://cloud.google.com/iam/docs/understanding-roles#cloudtrace.agent
		"roles/errorreporting.writer",                                  // https://cloud.google.com/iam/docs/understanding-roles#errorreporting.writer
		"roles/logging.logWriter",                                      // https://cloud.google.com/iam/docs/understanding-roles#logging.logWriter
		"roles/monitoring.metricWriter",                                // https://cloud.google.com/iam/docs/understanding-roles#monitoring.metricWriter
		"roles/run.serviceAgent",                                       // https://cloud.google.com/iam/docs/understanding-roles#run.serviceAgent
		"roles/cloudsql.client",                                        // https://cloud.google.com/iam/docs/understanding-roles#cloudsql.client
		"roles/bigquery.jobUser",                                       // https://cloud.google.com/iam/docs/understanding-roles#bigquery.jobUser
		"roles/datastore.user",                                         // https://cloud.google.com/iam/docs/understanding-roles#datastore.user
		"roles/cloudprofiler.agent",                                    // https://cloud.google.com/iam/docs/understanding-roles#cloudprofiler.agent
		"roles/pubsub.viewer",                                          // https://cloud.google.com/iam/docs/understanding-roles#pubsub.viewer
		"organizations/344471582607/roles/observability.metricsWriter", // https://github.com/kramphub/kramphub-gcp-iam-tf/pull/338
		"roles/apigee.environmentAdmin",                                // Needed for Apigee CICD pipeline SA to deploy API proxies
		"roles/apigee.apiAdminV2",                                      // Needed for Apigee CICD pipeline SA to deploy API proxies
	}

	schema := &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{
				Name: "role",
			},
		},
	}

	resources, err := runner.GetResourceContent(resourceName, schema, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		logger.Debug(fmt.Sprintf("`%s` block found, attributes: %#v", resourceName, resource.Body.Attributes))
		attribute, exists := resource.Body.Attributes["role"]
		if !exists {
			logger.Debug("attribute 'role' not found")
			continue
		}
		logger.Debug(fmt.Sprintf("attribute 'role' found: %#v", attribute))
		err := runner.EvaluateExpr(attribute.Expr, func(role string) error {
			if slices.Contains(allowedRoles, role) {
				return nil
			}
			// For 'range' it's also possible to use `attribute.Expr.Range()`. But then an 'ignore' comment on block level is not respected (it should then be on attribute level).
			return runner.EmitIssue(rule, fmt.Sprintf("`%s` is not allowed for `%s`. Please grant permissions on resource level (principle of least privilege).", role, resourceName), resource.DefRange)
		}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
