package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/vcollector/internal/conf"
	"github.com/gvcgo/vcollector/internal/gh"
)

const (
	ShaFileName            string = "sdklist_sha256.json"
	VersionFileNamePattern string = "%s.version.json"
)

type Sha256List map[string]string // fileName -> sha256

/*
1. Check sha256
2. Upload file to remote repo
3. Delete file from remote repo
*/
type Uploader struct {
	ShaFile    string
	VersionDir string
	Github     *gh.Github
	Sha256List Sha256List
}

func NewUploader() (u *Uploader) {
	u = &Uploader{
		ShaFile:    filepath.Join(conf.GetWorkDir(), ShaFileName),
		VersionDir: conf.GetVersionDir(),
		Github:     gh.NewGithub(),
		Sha256List: make(Sha256List),
	}
	u.loadSha256Info()
	return
}

func (u *Uploader) loadSha256Info() {
	if ok, _ := gutils.PathIsExist(u.ShaFile); ok {
		content, _ := os.ReadFile(u.ShaFile)
		json.Unmarshal(content, &u.Sha256List)
	}
}

func (u *Uploader) saveSha256Info() {
	content, _ := json.Marshal(u.Sha256List)
	os.WriteFile(u.ShaFile, content, os.ModePerm)
}

func (u *Uploader) getVersionFilePath(sdkName string) string {
	fName := fmt.Sprintf(VersionFileNamePattern, sdkName)
	return filepath.Join(u.VersionDir, fName)
}

func (u *Uploader) saveVersionFile(sdkName string, content []byte) {
	os.WriteFile(u.getVersionFilePath(sdkName), content, os.ModePerm)
}

func (u *Uploader) checkSha256(sdkName string, content []byte) (ok bool) {
	h := sha256.New()
	h.Write(content)
	shaStr := fmt.Sprintf("%x", h.Sum(nil))

	if len(u.Sha256List) == 0 {
		u.loadSha256Info()
	}

	if ss, ok1 := u.Sha256List[sdkName]; !ok1 {
		u.Sha256List[sdkName] = shaStr
		u.saveSha256Info()
		u.saveVersionFile(sdkName, content)
		return true
	} else {
		if ss == shaStr {
			return false
		} else {
			u.Sha256List[sdkName] = shaStr
			u.saveSha256Info()
			u.saveVersionFile(sdkName, content)
			return true
		}
	}
}

func (u *Uploader) Upload(sdkName string, content []byte) {
	if u.checkSha256(sdkName, content) {
		localFilePath := u.getVersionFilePath(sdkName)
		remoteFilePath := filepath.Base(localFilePath)
		u.Github.UploadFile(remoteFilePath, localFilePath)
	}
}