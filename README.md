# Cyral Terraform Provider

The Cyral Terraform Provider contains resources that can be used to interact with the Cyral API through Terraform code. It allows customers to maintain a history of changes in Cyral environment by transforming configurations into code and use Terraform to control state changes.

Our provider uses the same naming conventions and organization as stated in Terraform guidelines for [writing custom providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).

## Documentation

Full and comprehensive documentation for this provider with detailed description of its **resources**, **data sources** and **usage guides** are available in the [user documentation index](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs).

## Version history and compatibility

Please refer to our [Change Log](CHANGELOG.md) to learn about our version history, its features, bug fixes and Control Plane compatibility.

## Building, Documenting and Testing this Project

### Build Instructions

In order to build this repository, follow the steps below:

1.  Clone [terraform-provider-cyral](https://github.com/cyralinc/terraform-provider-cyral) repo from GitHub;

2.  Go to the root directory of the cloned repo using Linux shell and execute `make`. The build process will create binaries in directory `out` for both `darwin` and `linux` 64 bits. These binaries will be copied automatically to the local Terraform registry to be used by Terraform 13 and later.

Alternatively, you can use the dockerfile to build the image using `make docker-compose/build`

To use the local provider, the module must be configured to use the local provider path as follows:

```hcl
terraform {
  required_providers {
    cyral = {
      source = "local/terraform/cyral"
    }
  }
}
```

### Updating the Documentation

This project uses [`terraform-plugin-docs`](https://github.com/hashicorp/terraform-plugin-docs).
Add the attribute `Description` to the `Schema` in order to allow the documentation to be
created automatically and following Terraform good practices. See any of the resources in folder
`cyral` for guidance on how to document the `Schema`. See also folders `examples` and `templates`
for more information on where and how to store examples and define templates for documentation
artifacts.

To create the documentation automatically, run the commands:

```bash
# Creates the documentation files from the source code
make docker-compose/docs
# Runs the pre-commit linter
pre-commit run --show-diff-on-failure --color=always --all-files
```

> `pre-commit` can sometimes fail because your user is not the owner of the files in the `/docs` directory.
> To solve this problem, run the following command and re-run the `pre-commit run...` tried in the previous step:

```bash
find docs -exec sudo chown <your_username> {} \;
```

> The `make docker-compose/docs` command can sometimes fail. If this is your case, you can use the `tfplugindocs generate` command, which will do the same as `make docker-compose/docs`. You can get the binary from [this link](https://github.com/hashicorp/terraform-plugin-docs)

### Test Instructions

The test framework requires basic configuration before it can be executed as follows:

1. Set the configuration environment variables:

```bash
# Set the control plane DNS name and port:
export CYRAL_TF_CONTROL_PLANE=tenant.app.cyral.com

# Set client and secret ID:
export CYRAL_TF_CLIENT_ID=?
export CYRAL_TF_CLIENT_SECRET=?

# Initialize Terraform acceptance tests variable
export TF_ACC=true
```

2. Run `make`

#### Sweeper

(Feature still under implementation) To sweep leaked resources in the control
plane, run `make sweep`. The environment variables to access the control plane must be set as instructed
above.

### Commit instructions

This project uses [pre-commit](https://pre-commit.com/) to automatically lint changes during the commit process.

Before committing a change, you will need to install [`pre-commit`](https://pre-commit.com/#install) and then install
the hooks by running the following command in the root of the repository:

```shell
pre-commit install
```

### Running Project Built Locally

#### Terraform 0.12

Copy the desired binary file created in directory `out` (see [Build Instructions](#build-instructions)) to the root folder containing those `.tf` files that will be used to handle Cyral Terraform provider resources.

Run `terraform init` and proceed with `terraform apply` normally to execute your Terraform scripts.

#### Terraform 0.13+

Build the project using steps in [Build Instructions](#build-instructions), then proceed normally with `terraform init` and `terraform apply` commands.
