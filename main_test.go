package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Setenv("CONFIG_DIR", "templates/nginxConf.conf")
	os.Setenv("SETTING_FILE_NAME", "config.yaml")
	if !logic(os.Stdout) {
		t.FailNow()
	}
}
