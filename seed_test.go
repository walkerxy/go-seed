package seed

import (
	"testing"
)

func TestSeed(t *testing.T) {
	// testFile("./seeds", "admin_user.yaml")
	testDir("./seeds")
}

func testDir(dir string) {
	fileChan := make(chan string)
	go func() {
		WalkDir("./seeds", "yaml", fileChan)
		close(fileChan)
	}()

	for file := range fileChan {
		testFile(dir, file)
	}
}

func testFile(filepath string, filename string) {
	seed := NewSeed(filepath, filename)
	seed.SetTablePrefix("mc")
	seed.Fill()
}
