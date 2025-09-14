package tester

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/h2non/baloo.v3"
)

func NewAPITester() *baloo.Client {
	loadEnvConfig()
	return baloo.New(viper.GetString("API_URL"))
}

func loadEnvConfig() {
	rootDir := rootDir()
	file := path.Join(rootDir, "/.env")

	viper.SetConfigFile(file)
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return strings.TrimSpace(filepath.Dir(d))
}
