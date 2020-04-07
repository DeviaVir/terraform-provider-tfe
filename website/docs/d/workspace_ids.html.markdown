---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_workspace_ids"
sidebar_current: "docs-datasource-tfe-workspace-ids"
description: |-
  Get information on (external) workspace IDs.
---

# Data Source: tfe_workspace_ids

Use this data source to get a map of (external) workspace IDs.

## Example Usage

```hcl
data "tfe_workspace_ids" "app-frontend" {
  names        = ["app-frontend-prod", "app-frontend-dev1", "app-frontend-staging"]
  organization = "my-org-name"
}

data "tfe_workspace_ids" "all" {
  names        = ["*"]
  organization = "my-org-name"
}
```

## Argument Reference

The following arguments are supported:

* `names` - (Required) A list of workspace names to search for. Names that don't
  match a real workspace will be omitted from the results, but are not an error.

    To select _all_ workspaces for an organization, provide a list with a single
    asterisk, like `["*"]`. No other use of wildcards is supported.
* `organization` - (Required) Name of the organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

~> **NOTE** In versions < 0.15.1, `ids` returned a map of workspaces names 
and their human-readable IDs, which looked like `<ORGANIZATION>/<WORKSPACE>`.

* `ids` - A map of workspace names and their opaque external IDs, which
  look like `ws-<RANDOM STRING>`.
* `external_ids` - **Deprecated** Use `ids` instead. A map of workspace names and their opaque external IDs, which
  look like `ws-<RANDOM STRING>`.
