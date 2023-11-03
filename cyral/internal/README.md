## Notes on testing

Please make sure to use the function `utils.AccTestName` for all control plane
resource names, such as repository or sidecar names. This way, we can have
consistent prefixes for all resource names, which enables Terraform test
sweepers, avoids name clashes, and facilitates developing and maintaining test
code. If `utils.AccTestName` is generating a name that is invalid for a particular
resource type, please mention it in the code explicitly through comments, and
adjust the test sweeper. We can also create other functions, such as
`utils.AccTestNameUnderscore`, in the future, if necessary.

If you wish to test what Terraform calls an _in-place update_, be sure that the
resource is not being recreated in between the test steps. For example, if the
resource name is modified, or a `ForceNew` argument is modified, the resource
will be recreated.
