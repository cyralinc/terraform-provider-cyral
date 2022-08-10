## Notes on testing

Please make sure to use the function `accTestName`, whose definition can be
found in `testutils.go`, for all control plane resource names. This way, we can
have consistent prefixes for all resource names, which makes Terraform test
sweepers possible, and avoids name clashes. If `accTestName` is generating a
name that is invalid for a particular resource type, please mention it in the
code explicitly through comments, and adjust the test sweeper. We can also
create other functions, such as `accTestNameUnderscore`, in the future.
