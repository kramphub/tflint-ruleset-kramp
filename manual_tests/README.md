Some manual 'smoke tests'.
These can be useful to verify that the rules are working as expected when providing a more real-world configuration.

Please find the unit tests in the [rules](../rules) directory (ending with `_test.go`).

**Good to know:** normally `tflint --init` will download the ruleset on the fly, from this repository.
Which is different than what happens in [cloudbuild-test.yaml](../cloudbuild-test.yaml), since it uses the binary that was build just before.

---

A simple example of how to run the tests locally:

_Prerequisites: `go`, `tflint` and `jq`._

```shell
(cd .. && make test && make install)
TFLINT_LOG=debug tflint --enable-plugin=kramp --chdir ./verify
tflint --enable-plugin=kramp --chdir ./verify --no-color --format=json | jq ".issues"
tflint --enable-plugin=kramp --chdir ./verify --no-color --format=json | jq ".issues | length"
tflint --enable-plugin=kramp --chdir ./verify --no-color --format=json | jq ".issues | length == 10"
```

```shell
(cd .. && make test && make install)
TFLINT_LOG=debug tflint --enable-plugin=kramp --minimum-failure-severity=notice --chdir ./verify_ignores
tflint --enable-plugin=kramp --chdir ./verify_ignores --no-color --format=json | jq ".issues"
tflint --enable-plugin=kramp --chdir ./verify_ignores --no-color --format=json | jq ".issues | length == 0"
```
