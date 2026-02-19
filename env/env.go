package env

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/s-hammon/p"
)

var (
	homeDir           = os.Getenv("GRAPH_HOME_DIR")
	homeDirConfigFile = ""
	homeDirCacheFile  = ""
)

func init() {
	homeDir = SetHomeDir("msgraph")
	os.MkdirAll(homeDir, 0750)

	homeDirConfigFile = GetConfigPath(homeDir)
	if homeDir != "" && !pathExists(homeDirConfigFile) {
		WriteConfigFile([]byte("{}"))
	}

	homeDirCacheFile = GetCachePath(homeDir)
	if homeDir != "" && !pathExists(homeDirCacheFile) {
		WriteCacheFile("", make(map[string]string))
	}
}

func SetHomeDir(name string) string {
	envKey := strings.ToUpper(strings.ReplaceAll(name, "-", "_")) + "_HOME_DIR"
	dir := os.Getenv(envKey)
	if dir == "" {
		dir = path.Join(GetHomeDir(), "."+name)
		os.Setenv(envKey, dir)
	}

	return dir
}

func GetConfigPath(dir string) string {
	return path.Join(dir, "config.json")
}

func LoadConfigFile() []byte {
	data, _ := os.ReadFile(homeDirConfigFile)
	return data
}

func WriteConfigFile(data []byte) {
	if err := os.WriteFile(homeDirConfigFile, data, 0600); err != nil {
		msg := p.Format("couldn't write to %s: %v", homeDirConfigFile, err)
		panic(msg)
	}
}

func GetHomeDir() string {
	home := os.Getenv("HOME")
	switch runtime.GOOS {
	case "windows":
		home = p.Coalesce(os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH"), os.Getenv("USERPROFILE"), home)
	case "linux":
		home = p.Coalesce(os.Getenv("XDG_CONFIG_HOME"), home)
	}

	return home
}

func SafeOpen(path string) (*os.File, error) {
	root, err := os.OpenRoot(GetHomeDir())
	if err != nil {
		return nil, fmt.Errorf("os.OpenRoot(home): %v", err)
	}
	defer root.Close()

	return root.Open(path)
}

func pathExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
