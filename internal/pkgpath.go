package internal

// Copied from https://github.com/mailru/easyjson/pull/185 with some changes.
// Thanks to @zelenin for original realization.

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func GetPkgPath(fname string, isDir bool) (string, error) {
	if !filepath.IsAbs(fname) {
		pwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		fname = filepath.Join(pwd, fname)
	}
	goModPath, _ := goModPath(fname, isDir)
	if strings.Contains(goModPath, "go.mod") {
		pkgPath, err := getPkgPathFromGoMod(fname, isDir, goModPath)
		if err != nil {
			return "", err
		}
		return pkgPath, nil
	}
	return getPkgPathFromGOPATH(fname, isDir)
}

var (
	goModPathCache = make(map[string]string)
)

// empty if no go.mod, GO111MODULE=off or go without go modules support
func goModPath(fname string, isDir bool) (string, error) {
	root := fname
	if !isDir {
		root = filepath.Dir(fname)
	}
	goModPath, ok := goModPathCache[root]
	if ok {
		return goModPath, nil
	}
	defer func() {
		goModPathCache[root] = goModPath
	}()
	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Dir = root
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	goModPath = string(bytes.TrimSpace(stdout))
	return goModPath, nil
}

func getPkgPathFromGoMod(fname string, isDir bool, goModPath string) (string, error) {
	modulePath := getModulePath(goModPath)
	if modulePath == "" {
		return "", fmt.Errorf("cannot determine module path from %s", goModPath)
	}
	rel := path.Join(modulePath, filePathToPackagePath(strings.TrimPrefix(fname, filepath.Dir(goModPath))))
	if !isDir {
		return path.Dir(rel), nil
	}
	return path.Clean(rel), nil
}

var (
	modulePrefix          = []byte("\nmodule ")
	pkgPathFromGoModCache = make(map[string]string)
	gopathCache           = ""
)

func getModulePath(goModPath string) string {
	pkgPath, ok := pkgPathFromGoModCache[goModPath]
	if ok {
		return pkgPath
	}
	defer func() {
		pkgPathFromGoModCache[goModPath] = pkgPath
	}()
	data, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	var i int
	if bytes.HasPrefix(data, modulePrefix[1:]) {
		i = 0
	} else {
		i = bytes.Index(data, modulePrefix)
		if i < 0 {
			return ""
		}
		i++
	}
	line := data[i:]
	// Cut line at \n, drop trailing \r if present.
	if j := bytes.IndexByte(line, '\n'); j >= 0 {
		line = line[:j]
	}
	if line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	line = line[len("module "):]
	// If quoted, unquote.
	pkgPath = strings.TrimSpace(string(line))
	if pkgPath != "" && pkgPath[0] == '"' {
		s, err := strconv.Unquote(pkgPath)
		if err != nil {
			return ""
		}
		pkgPath = s
	}
	return pkgPath
}

func GetRelatedFilePath(pkg string) (string, error) {
	if gopathCache == "" {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			var err error
			gopath, err = getDefaultGoPath()
			if err != nil {
				return "", fmt.Errorf("cannot determine GOPATH: %s", err)
			}
		}
		gopathCache = gopath
	}
	paths := allPaths(filepath.SplitList(gopathCache))
	for _, p := range paths {
		checkingPath := filepath.Join(p, pkg)
		if info, err := os.Stat(checkingPath); err == nil && info.IsDir() {
			return checkingPath, nil
		}
	}
	return "", fmt.Errorf("file '%v' is not in GOROOT or GOPATH. Checked paths:\n%s", pkg, strings.Join(paths, "\n"))
}

func allPaths(gopaths []string) []string {
	const _2 = 2
	res := make([]string, len(gopaths)+_2)
	res[0] = filepath.Join(build.Default.GOROOT, "src")
	res[1] = "vendor"
	for i := range res[_2:] {
		res[i+_2] = filepath.Join(gopaths[i], "src")
	}
	return res
}

func getPkgPathFromGOPATH(fname string, isDir bool) (string, error) {
	if gopathCache == "" {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			var err error
			gopath, err = getDefaultGoPath()
			if err != nil {
				return "", fmt.Errorf("cannot determine GOPATH: %s", err)
			}
		}
		gopathCache = gopath
	}
	for _, p := range filepath.SplitList(gopathCache) {
		prefix := filepath.Join(p, "src") + string(filepath.Separator)
		if rel := strings.TrimPrefix(fname, prefix); rel != fname {
			if !isDir {
				return path.Dir(filePathToPackagePath(rel)), nil
			} else {
				return path.Clean(filePathToPackagePath(rel)), nil
			}
		}
	}
	return "", fmt.Errorf("file '%v' is not in GOPATH", fname)
}

func filePathToPackagePath(path string) string {
	return filepath.ToSlash(path)
}

func getDefaultGoPath() (string, error) {
	if build.Default.GOPATH != "" {
		return build.Default.GOPATH, nil
	}
	output, err := exec.Command("go", "env", "GOPATH").Output()
	return string(bytes.TrimSpace(output)), err
}