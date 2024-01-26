---
title: server
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Untyped API server

The toolkit supports serving a swagger spec with untyped data. This means that it uses mostly interface{} as params to each operation and as result type. It does allow you to serve a spec up quickly. This is one of the building blocks required to serve up stub APIs and to generate a test server with predictable responses.

<!--more-->

In this document we'll use the todo list application to build a server that doesn't use any code generation. So we'll load a predefined swagger specification document and serve it up, while writing a minimal amount of code and try to avoid requiring a type cast.  

Let's start with just 
