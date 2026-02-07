package env

import (
	"embed"
	"encoding/json"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/s-hammon/msgraph/auth"
	"github.com/s-hammon/p"
)

var (
	homeDir        = os.Getenv("GRAPH_HOME_DIR")
	homeDirDotFile = ""
)

//go:embed *
var envFolder embed.FS

func init() {
	homeDir = SetHomeDir("ms-graph")
	homeDirDotFile = GetDotFilePath(homeDir)

	os.MkdirAll(homeDir, 0755)
	if homeDir != "" && !pathExists(homeDirDotFile) {
		WriteDotFile(auth.AzureADConfig{})
	}
}

func SetHomeDir(name string) string {
	envKey := strings.ToUpper(strings.ReplaceAll(name, "-", "_")) + "_HOME_DIR"
	dir := os.Getenv(envKey)
	if dir == "" {
		dir = path.Join(getHomeDir(), "."+name)
		os.Setenv(envKey, dir)
	}

	return dir
}

func GetDotFilePath(dir string) string {
	return path.Join(dir, "config.json")
}

func LoadDotFile() auth.AzureADConfig {
	data, _ := os.ReadFile(homeDirDotFile)

	var cfg auth.AzureADConfig
	json.Unmarshal(data, &cfg)
	return cfg
}

func WriteDotFile(cfg auth.AzureADConfig) {
	if cfg.Scopes == nil {
		cfg.Scopes = make([]string, 0)
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(homeDirDotFile, data, 0644); err != nil {
		msg := p.Format("couldn't write to %s: %v", homeDirDotFile, err)
		panic(msg)
	}
}

func getHomeDir() string {
	home := os.Getenv("HOME")
	switch runtime.GOOS {
	case "windows":
		home = p.Coalesce(os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH"), os.Getenv("USERPROFILE"), home)
	case "linux":
		home = p.Coalesce(os.Getenv("XDG_CONFIG_HOME"), home)
	}

	return home
}

func pathExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
