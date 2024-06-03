#!/usr/bin/env sh

if ! tflint --version
then
    echo "tflint is not installed!"
    exit 1
fi

mkdir -p ~/.tflint.d/plugins/
if ! cp ../tflint-ruleset-kramp ~/.tflint.d/plugins/
then
    echo "Failed to copy plugin to ~/.tflint.d/plugins/. Please build it first."
    exit 1
fi

# TODO: use the JSON output format and parse it to check the results, instead of only looking at the exit code?

echo "*** Start verifying that the rules are applied"
cd verify || exit 1
if tflint --enable-plugin=kramp --no-color
then
    echo "tflint succeeded, but it is expected to fail. Was the plugin loaded?"
    exit 1
fi
cd ..
echo
echo "*** Done verifying that the rules are applied"

echo "*** Start verifying that the ignores are respected"
cd verify_ignores || exit 1
if ! tflint --enable-plugin=kramp --no-color
then
    echo "tflint failed, but it is expected to succeed. The ignores should have been respected."
    exit 1
fi
cd ..
echo "*** Done verifying that the ignores are respected"
echo
