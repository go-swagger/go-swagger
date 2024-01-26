---
title: swagger diff
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 10
---
# Inspecting differences between swagger specs

The toolkit has a command to display differences between two swagger specifications.

### Usage

To diff specifications:

```
Usage:
  swagger [OPTIONS] diff [diff-OPTIONS] {original spec} {spec}

diff specs showing which changes will break existing clients

Application Options:
  -q, --quiet                    silence logs
      --log-output=LOG-FILE      redirect logs to file

Help Options:
  -h, --help                     Show this help message

[diff command options]
      -b, --break                When present, only shows incompatible changes
      -f, --format=[txt|json]    When present, writes output as json (default: txt)
      -i, --ignore=              Exception file of diffs to ignore (copy output from json diff format) (default: none specified)
      -d, --dest=                Output destination file or stdout (default: stdout)
```
