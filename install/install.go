package install

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/cavaliergopher/grab/v3"
)

func GetDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(home, ".polym1")
}

func Install() {
	dataDir := GetDataDir()
	fmt.Println("installing PolyM1...")

	os.RemoveAll(dataDir)
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasSuffix(os.Args[0], "go") {
		Copy(os.Args[0], path.Join(dataDir, "polym1"))
		os.Chmod(path.Join(dataDir, "polym1"), 0755)
	}

	zipPath := path.Join(dataDir, "files.zip")

	_, err = grab.Get(zipPath, "https://github.com/PolyM1/files/archive/refs/heads/main.zip")

	if err != nil {
		log.Fatal(err)
	}

	Unzip(zipPath, dataDir)
	os.Remove(zipPath)
}
