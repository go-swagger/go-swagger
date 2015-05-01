# Bootie Docs

**Bootie Docs** is a simple [hugo](http://gohugo.io/) theme for documentation.  
The name "bootie" comes from [Bootstrap](http://getbootstrap.com/) CSS.

![Bootie Docs screenshot](https://raw.githubusercontent.com/key-amb/hugo-theme-bootie-docs/master/images/tn.png)

You can see demo and full documentation at http://key-amb.github.io/bootie-docs-demo/ .

## CONTENTS

* [QUICKSTART](#quickstart)
* [OPTIONS](#options)
* [LIMITATION](#limitation)
* [DEPENDENCIES](#dependencies)
* [LICENSE](#license)

## QUICKSTART

1. `hugo new _index.md`
1. Edit `content/_index.md`

Then the content appears on top page.

## OPTIONS

You can customize the menu items in the header navigation bar by configuring `params.mainMenu` in your _config.toml_ (or _config.yaml_).

```
# example of config.toml
[params]
  mainMenu = ["about", "usage"]
```

All other options and usages are described at the documentation site -- http://key-amb.github.io/bootie-docs-demo/ .

## LIMITATION

Because _Bootie Docs_ is developed for documentation, it lacks many blog-type facilities such as RSS feeds, pagination of posts and so on.

## DEPENDENCIES

**Bootie Docs** includes following libraries:

* [Bootstrap](http://getbootstrap.com/) v3.3.4 ... Well-known CSS framework.
* [jQuery](https://jquery.com/) v1.11.2 ... Requried by _Bootstrap_.
* [highlight.js](https://highlightjs.org/) v8.5 ... For syntax highlighting.

## LICENSE

Copyright (C) 2015 YASUTAKE Kiyoshi.  
Released under the MIT License.
