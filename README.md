# xspf

```
Usage: xspf (-a | -v | -e=<ext>...)... [-s] [-o=<output file>] DIRS...

print an xspf (VLC) playlist to standard output

Examples:
# find audio files in the current directory
  xspf -a . > playlist.xspf

# find and shuffle video files in dir1 and dir2
  xspf -v -s dir1 dir2 > playlist.xspf

# find *.ext1 and *.ext2 files in dir1 and dir2
  xspf -e ext1 -e ext2 dir1 dir2 > playlist.xspf

Arguments:
  DIRS            directories to search

Options:
  -a, --audio     match a bunch of audio extensions
  -v, --video     match a bunch of video extensions
  -e, --ext       extension to match; can be used multiple times
  -s, --shuffle   shuffle the output
  -o, --output    output to file instead stdout
```

## Installation

Once you've installed the [go compiler](https://www.golang.org),

```
go install github.com/7fffffff/xpsf@latest
```

## License

Blue Oak Model License 1.0.0:
<https://blueoakcouncil.org/license/1.0.0>