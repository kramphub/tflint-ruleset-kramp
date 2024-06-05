package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-kramp/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "kramp",
			Version: "0.0.0-semantically-released",
			Rules: []tflint.Rule{
				rules.NewAuthoritativeIAMPolicyOnFolderLevelRule(),
				rules.NewAuthoritativeIAMPolicyOnProjectLevelRule(),
				rules.NewIAMPolicyOnFolderLevelRule(),
				rules.NewIAMPolicyOnProjectLevelRule(),
				rules.NewAuthoritativeIAMPolicyOnResourceLevelRule(),
				rules.NewCreatingKeyForServiceAccountRule(),
			},
		},
	})
}
