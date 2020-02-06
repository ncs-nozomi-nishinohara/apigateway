package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	if !logic(os.Stdout) {
		t.FailNow()
	}
}
