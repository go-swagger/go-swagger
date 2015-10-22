# Material-Design

Material-Design is a simple material design theme for [Hugo](http://gohugo.io/).

![](https://github.com/pdevty/material-design/blob/master/images/tn.png)

demo : [http://pdevty.github.io/blog/](http://pdevty.github.io/blog/)

## Features

- Simple Material Design by [Materialize](http://materializecss.com/)
- Google Analytics (optional)
- Pagination
- Disqus (optional)
- Twitter, Facebook, GitHub, Google+, LinkedIn links (optional)
- Tags
- Categories
- Cover image (optional)
- Highlighting source code

## Installation

```shell
$ mkdir themes
$ cd themes
$ git clone https://github.com/pdevty/material-design
```

## Usage

```shell
$ hugo server -t material-design -w -D
```

## Configuration

config.toml

```toml
theme="material-design"
baseurl = "Your Site URL"
languageCode = "en-us"
title = "Your Site Title"
MetaDataFormat = "toml"
paginate = 9 # To specify a multiple of 3
disqusShortname = "Your Disqus Name" # optional
copyright = "© 2015 Copyright Text"

[params]
  description = "Your Site Description" # optional
  twitter = "Your Twitter Name" # optional
  github = "Your Github Name" # optional
  facebook = "Your facebook Name" # optional
  gplus = "Your Google+ profile name" # optional
  linkedin = "Your LinkedIn Name" # optional
  headerCover = "images/headerCover.png" # optional
  footerCover = "images/footerCover.png" # optional
  googleAnalyticsUserID = "Your Analytics User Id" # optional

[permalinks]
  post = "/:year/:month/:day/:title/" # optional
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
