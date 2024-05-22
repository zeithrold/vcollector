package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type Fd struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewFd() (f *Fd) {
	f = &Fd{
		SDKName:  "fd",
		RepoName: "sharkdp/fd",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (f *Fd) GetSDKName() string {
	return f.SDKName
}

func (f *Fd) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (f *Fd) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".deb") {
		return false
	}
	if strings.HasSuffix(a.Name, "-linux-musl.tar.gz") {
		return false
	}
	if strings.HasSuffix(a.Name, "-windows-msvc.zip") {
		return false
	}
	return true
}

func (f *Fd) osParser(fName string) (osStr string) {
	if strings.Contains(fName, "darwin") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "windows") {
		return "windows"
	}
	return
}

func (f *Fd) archParser(fName string) (archStr string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-aarch64") {
		return "arm64"
	}
	return
}

func (f *Fd) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (f *Fd) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (f *Fd) Start() {
	f.GhSearcher.Search(
		f.RepoName,
		f.tagFilter,
		f.fileFilter,
		f.vParser,
		f.archParser,
		f.osParser,
		f.insParser,
		nil,
	)
}

func TestFd() {
	nn := NewFd()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}