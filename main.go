package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/7fffffff/jsf"
	cli "github.com/jawher/mow.cli"
)

func main() {
	useAudioPresetExts := false
	useVideoPresetExts := true
	exts := []string{}
	outputPath := ""
	roots := []string{}
	cfg := &config{
		rng: rand.New(jsf.New(time.Now().UnixNano())),
	}
	app := cli.App("xspf", "print an xspf (VLC) playlist to standard output")
	app.LongDesc = `print an xspf (VLC) playlist to standard output

Examples:
# find audio files in the current directory
  xspf -a . > playlist.xspf

# find and shuffle video files in dir1 and dir2
  xspf -v -s dir1 dir2 > playlist.xspf

# find *.ext1 and *.ext2 files in dir1 and dir2
  xspf -e ext1 -e ext2 dir1 dir2 > playlist.xspf`
	app.Spec = "(-a | -v | -e=<ext>...)... [-o=<output file>] [OPTIONS] DIRS..."
	app.StringsArgPtr(&roots, "DIRS", []string{}, "directories to search")
	app.BoolOptPtr(&useAudioPresetExts, "a audio", false, "match a bunch of audio extensions")
	app.BoolOptPtr(&useVideoPresetExts, "v video", false, "match a bunch of video extensions")
	app.StringsOptPtr(&exts, "e ext", []string{}, "extension to match; can be used multiple times")
	app.BoolOptPtr(&cfg.mpcpl, "mpcpl", false, "output as mpcpl (Media Player Classic)")
	app.BoolOptPtr(&cfg.shuffle, "s shuffle", false, "shuffle the output")
	app.StringOptPtr(&outputPath, "o output", "", "output to file instead stdout")
	_ = app.BoolOpt("h help", false, "show help with examples")
	app.Action = func() {
		err := func() error {
			if useAudioPresetExts {
				exts = append(exts, presetAudioExts()...)
			}
			if useVideoPresetExts {
				exts = append(exts, presetVideoExts()...)
			}
			cfg.exts = newExtMatcher(exts)
			if cfg.exts.Empty() {
				app.PrintHelp()
				cli.Exit(1)
			}
			var output io.Writer = os.Stdout
			if outputPath != "" {
				file, err := os.Create(outputPath)
				if err != nil {
					return err
				}
				defer file.Close()
				output = file
			}
			return cfg.WriteAll(output, roots)
		}()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			cli.Exit(1)
		}
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

type config struct {
	exts    *extMatcher
	rng     *rand.Rand
	shuffle bool
	mpcpl   bool
}

func (cfg *config) WriteAll(wr io.Writer, roots []string) error {
	allFiles := []string{}
	for _, root := range roots {
		root = filepath.Clean(root)
		err := filepath.Walk(root, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
					fmt.Fprintln(os.Stderr, err)
					return nil
				}
				return err
			}
			if info.IsDir() {
				return nil
			}
			if info.Size() < 1 {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(p))
			if ext == "" {
				return nil
			}
			if cfg.exts.Match(ext) {
				abs, err := filepath.Abs(p)
				if err != nil {
					return err
				}
				if cfg.mpcpl {
					allFiles = append(allFiles, abs)
				} else {
					// xspf paths need to look like
					// /C:/A/B/C.mp3
					allFiles = append(allFiles, path.Clean("/"+filepath.ToSlash(abs)))
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	if cfg.shuffle && cfg.rng != nil {
		cfg.rng.Shuffle(len(allFiles), func(i, j int) {
			allFiles[i], allFiles[j] = allFiles[j], allFiles[i]
		})
	}
	if cfg.mpcpl {
		mpc := &mpcplWriter{
			wr: wr,
		}
		err := mpc.Begin()
		if err != nil {
			return err
		}
		for _, f := range allFiles {
			err = mpc.WriteFile(f)
			if err != nil {
				return err
			}
		}
		return nil
	}
	xspf := newXSPFWriter(wr)
	err := xspf.Begin()
	if err != nil {
		return err
	}
	for _, f := range allFiles {
		err = xspf.WriteTrack(f)
		if err != nil {
			return err
		}
	}
	return xspf.End()
}
