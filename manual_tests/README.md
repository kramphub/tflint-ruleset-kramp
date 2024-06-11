Some manual 'smoke tests'.
These can be useful to verify that the rules are working as expected when providing a more real-world configuration.

Please find the unit tests in the [rules](../rules) directory (ending with `_test.go`).

**Good to know:** normally `tflint --init` will download the ruleset on the fly, from this repository.
Which is different than what happens in [cloudbuild-test.yaml](../cloudbuild-test.yaml), since it uses the binary that was build just before.

---

For example:

```shell
cd ..
make test && make install
cd manual_tests/verify
TFLINT_LOG=debug tflint --enable-plugin=kramp
cd ..
```

```shell
cd ..
make test && make install
cd manual_tests/verify_ignores
TFLINT_LOG=debug tflint --enable-plugin=kramp --minimum-failure-severity=notice
cd ..
```

Note that for the command below you need to build the plugin for the CPU architecture used by the Docker image.
See second example to see how to do that.

```shell
PARENT_DIR=$(dirname "$(pwd)")
docker run -it --rm \
  --name "tflint-ruleset-kramp-gcp" \
  --volume "${PARENT_DIR}/:/workspace/" \
  --workdir="/workspace/manual_tests/" \
  --env "USERID=$(id -u):$(id -g)" \
  --entrypoint="./run_in_cloudbuild.sh" \
  ghcr.io/terraform-linters/tflint:v0.51.1
```

This shows how to build a compatible binary for the `tflint` Docker image:

```shell
cd ..
docker run -it --rm \
  --name "tflint-ruleset-kramp" \
  --volume "$(pwd):/workspace" \
  --workdir="/workspace" \
  --env "USERID=$(id -u):$(id -g)" \
  --entrypoint=sh \
  golang:1.22-alpine3.19 \
  -c "go mod tidy && go build && go test -v ./..."
cd manual_tests
```
