package main

import (
	"sort"
	"strings"
)

func presetAudioExts() []string {
	return []string{
		".aac",
		".ac3",
		".ape",
		".au",
		".dts",
		".flac",
		".m4a",
		".mka",
		".mp1",
		".mp2",
		".mp3",
		".ogg",
		".oma",
		".ra",
		".spx",
		".wav",
		".wma",
	}
}
func presetVideoExts() []string {
	return []string{
		".3g2",
		".3gp",
		".a52",
		".asf",
		".avi",
		".divx",
		".dv",
		".flv",
		".f4v",
		".gifv",
		".ifo",
		".iso",
		".m1v",
		".m2ts",
		".m2v",
		".m4v",
		".mkv",
		".mov",
		".mp4",
		".mpeg",
		".mpeg1",
		".mpeg2",
		".mpeg4",
		".mpg",
		".mts",
		".ogm",
		".rm",
		".rmvb",
		".ts",
		".vob",
		".webm",
		".wmv",
	}
}

type extMatcher struct {
	exts []string
}

func newExtMatcher(exts []string) *extMatcher {
	normalized := make([]string, 0, len(exts))
	for _, e := range exts {
		if e == "" {
			continue
		}
		if e[0] != '.' {
			e = "." + e
		}
		normalized = append(normalized, strings.ToLower(e))
	}
	sort.Strings(normalized)
	return &extMatcher{
		exts: normalized,
	}
}

func (m *extMatcher) Empty() bool {
	return len(m.exts) == 0
}

func (m *extMatcher) Match(ext string) bool {
	if i := sort.SearchStrings(m.exts, ext); i < len(m.exts) && m.exts[i] == ext {
		return true
	}
	return false
}
