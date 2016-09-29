package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitURLPatterns()
	InitSessionsMap()

	ErrorTemplate = InitErrorTemplate("/home/valya/myprogs/yorablog/templates")

	os.Exit(m.Run())
}
