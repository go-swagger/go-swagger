---
title: Navigate the docs
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 5
---
# Summary

## General
  - [README](_index.md)
  - [About this project](about.md)
  - [Installing](install.md)
  - [Features list](features.md)
  - [Guidelines for maintainers](guidelines/_index.md)

## Working with the go-swagger tool
  - [Options and commands](usage/swagger.md)
  - [Shell autocompletion](cli_helpers.md)

{{< details title="Use-cases" open=false >}}
  - [Validate](usage/validate.md)
  - [Diff](usage/diff.md)

  - [Transform spec](reference/transform.md)
    - [Expand](usage/expand.md)
    - [Flatten](usage/flatten.md)
    - [Mixin](usage/mixin.md)

{{< tabs "usecase" >}}
{{< tab "Generate code from spec" >}}
  - [Generate code from spec](generate/_index.md)
    - API Client
      - [Generate an API Client](generate/client.md)
      - [Generate a Command Line client](generate/cli.md)

    - API Server
      - [Server Usage](generate/server.md)
      - [How to use the server](reference/server.md)

    - API Model
      - [Model Usage](generate/model.md)
      - [How to use models](reference/models/_index.md)

    - [Dependencies & Requirements](generate/requirements.md)
{{< /tab >}}

{{< tab "Generate spec from source code" >}}
  - [Generate spec from source code](generate-spec/_index.md)
      - [Spec Usage](generate/spec.md)
      - [Spec generation rules](generate-spec/spec/annotations.md)
        - [swagger:meta](generate-spec/spec/annotations/meta.md)
        - [swagger:route](generate-spec/spec/annotations/route.md)
        - [swagger:params](generate-spec/spec/annotations/params.md)
        - [swagger:operation](generate-spec/spec/annotations/operation.md)
        - [swagger:response](generate-spec/spec/annotations/response.md)
        - [swagger:model](generate-spec/spec/annotations/model.md)
        - [swagger:allOf](generate-spec/spec/annotations/allOf.md)
        - [swagger:strfmt](generate-spec/spec/annotations/strfmt.md)
        - [swagger:discriminated](generate-spec/spec/annotations/discriminated.md)
        - [swagger:ignore](generate-spec/spec/annotations/ignore.md)
{{< /tab >}}

{{< tab "Document your API" >}}
  - [Serve UI](usage/serve_ui.md)
  - [Generate markdown](usage/markdown.md)
{{< /tab >}}

{{< /tabs >}}
{{< /details >}}

## [Tutorials](tutorial/_index.md)

{{< details title="Servers" open=false >}}
  - [Simple Server](tutorial/todo-list.md)
  - [Custom Server](tutorial/custom-server.md)
  - [Dynamic Server](tutorial/dynamic.md)
{{< /details >}}

{{< details title="Authentication" open=false >}}
  - [Authentication](tutorial/authentication.md)
  - [OAuth2](tutorial/oauth2.md)
  - [Composition](tutorial/composed-auth.md)
{{< /details >}}
    
## Advanced use-cases

- [Using middleware](reference/middleware.md)
- [Custom generation](reference/templates/template_layout.md)
  - [Custom templates](reference/templates/templates.md)

## [FAQ](faq/_index.md)
  - [Installation,setup and environment](faq/faq_setup.md)
  - [Model generation](faq/faq_model.md)
  - [Server generation & customization](faq/faq_server.md)
  - [Client generation](faq/faq_client.md)
  - [Spec generation from source](faq/faq_spec.md)
  - [API testing](faq/faq_testing.md)
  - [Documenting your API](faq/faq_documenting.md)
  - [Swagger specification](faq/faq_swagger.md)
