![Build Status](https://github.com/ricardobranco777/dedup/actions/workflows/ci.yml/badge.svg)

# dedup
Deduplicate files in directory

## Usage

```
Usage: ./dedup [OPTIONS] DIRECTORY...
  -n, --dry-run           dry run
  -x, --one-file-system   do not cross filesystems
  -q, --quiet             be quiet about it
      --version           print version and exit
```

## Deprecation notice

The [util-linux](https://github.com/util-linux/util-linux) package, available in other Unixes, has a [hardlink](https://man7.org/linux/man-pages/man1/hardlink.1.html) command that does the same and also supports reflinks.
