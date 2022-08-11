## Notes on testing

Please make sure to use the function `accTestName` for all control plane
resource names, such as repository or sidecar names. This way, we can have
consistent prefixes for all resource names, which enables Terraform test
sweepers, avoids name clashes, and facilitates developing and maintaining test
code. If `accTestName` is generating a name that is invalid for a particular
resource type, please mention it in the code explicitly through comments, and
adjust the test sweeper. We can also create other functions, such as
`accTestNameUnderscore`, in the future.

If you wish to test what Terraform calls an _in-place update_, be sure that the
resouce is not being recreated in between the test steps. For example, if the
resource name is modified, or a `ForceNew` argument is modified, the resource
will be recreated.
