# tflint-ruleset-kramp

This is a custom ruleset for [tflint](https://github.com/terraform-linters/tflint).

We use this publicly available ruleset to enforce our best practices and conventions in our Terraform code.

## Usage

It is enabled in `.tflint.hcl` files, which are present in the repositories for `terraform-infra` and `terraform-module`.
For example like so:

```hcl
plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

plugin "google" {
  enabled = true
  version = "0.29.0"
  source  = "github.com/terraform-linters/tflint-ruleset-google"
}

plugin "kramp" {
  enabled = true
  version = "1.1.0"
  source  = "github.com/kramphub/tflint-ruleset-kramp"
}

config {
  call_module_type = "all"
}
```

`tflint --init` will then download the ruleset from GitHub and make it available for use.

## Releasing a new version

Using [goreleaser](https://goreleaser.com/) to build the binaries and deploy then to GitHub.
Installation instructions can be found [here](https://goreleaser.com/install/).

_An example for MacOS:_
```shell
brew install goreleaser
```

To try it out, simply only building the binaries without actual releasing:
```shell
goreleaser build --clean --snapshot
```

To do a 'dry-run' release (thus not actually releasing):
```shell
goreleaser release --clean --skip=publish --skip=validate
```

The actual release is done by this CloudBuild script: [cloudbuild-goreleaser.yaml](cloudbuild-goreleases.yaml).

## Development

This repository is based on this template: [tflint-ruleset-template](https://github.com/terraform-linters/tflint-ruleset-template).
There you can also find some example [rules](https://github.com/terraform-linters/tflint-ruleset-template/tree/main/rules).

And for reference, you can also find all the standard rules here:
- [tflint-ruleset-google](https://github.com/terraform-linters/tflint-ruleset-google/tree/master/rules)
- [tflint-ruleset-terraform](https://github.com/terraform-linters/tflint-ruleset-terraform/tree/main/rules)

The rules are written in Go, and are located in the [rules](rules) directory.
It is also possible to write rules in [OPA/Rego](https://github.com/terraform-linters/tflint-ruleset-opa) (as we tried during our efforts to use `trivy`). 
But Go gives more flexibility and allows for easy unit-testing, thus being more powerful.

#### Install the plugin

```shell
make test && make install
```

_It will be copied to the tflint plugins directory (`~/.tflint.d/plugins`)._

#### Enable debug logging

You can use this environment variable: `TFLINT_LOG=debug`

#### Trying it out on a terraform directory

```shell
cd manual_tests/verify
TFLINT_LOG=debug tflint --enable-plugin=kramp
```

#### If you don't have Go installed

Then you could use the GoLang Docker image to run the tests etc:

```shell
docker run -it --rm \
  --name "tflint-ruleset-kramp" \
  --volume "$(pwd):/workspace" \
  --workdir="/workspace" \
  --env "USERID=$(id -u):$(id -g)" \
  --entrypoint=bash \
  golang:1.22.3 \
  -c "go mod tidy && make build && make test"
```

If you want to build a binary for Apple Silicon you could do something like this:

```shell
docker run -it --rm \
  --name "tflint-ruleset-kramp" \
  --volume "$(pwd):/workspace" \
  --workdir="/workspace" \
  --env "USERID=$(id -u):$(id -g)" \
  --env "GOOS=darwin" \
  --env "GOARCH=arm64" \
  --entrypoint=bash \
  arm64v8/golang:1.22.3 \
  -c "go mod tidy && go build"
```

## Hints/Notes

### Ignoring a rule

Please do this only as a last resort! 

If you think a rule is wrong, please contact Team Cloud-Platform, so we can modify it accordingly.

To ignore a rule _(the white space after the colon is important)_:

```hcl
# !!! Add here a clear description of the reason why you disable this rule !!!
# tflint-ignore: authoritative_iam_policy_on_folder_level
resource "google_folder_iam_policy" "folder_admin_policy" {
  folder      = google_folder.department1.name
  policy_data = data.google_iam_policy.admin.policy_data
}
```

#### Regarding ranges (indicating where the problem is found)

It is possible, when you check on attribute level, to report as range the location of that attribute.
For 'range' you can thus use `attribute.Expr.Range()`. 
But then an 'ignore' comment on block level is not respected (it should then be on attribute level).
See [iam_policy_on_project_level_rule.go](rules/iam_policy_on_project_level_rule.go) for an example.
