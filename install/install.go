package install

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func GetDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(home, ".polym1")
}

func downloadFiles(zipPath string) {
	grabClient := grab.NewClient()
	grabReq, err := grab.NewRequest(zipPath, "https://github.com/PolyM1/files/archive/refs/heads/main.zip")

	if err != nil {
		log.Fatal(err)
	}

	grabRes := grabClient.Do(grabReq)
	bar := progressbar.Default(100)

	ticker := time.NewTicker(100)
	defer ticker.Stop()

OuterLoop:
	for {
		select {
		case <-ticker.C:
			bar.Set(int(grabRes.Progress() * 100))

		case <-grabRes.Done:
			break OuterLoop
		}
	}
}

func Install() {
	dataDir := GetDataDir()
	fmt.Println("installing PolyM1...")

	os.RemoveAll(dataDir)
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("copying binary...")
	if !strings.HasSuffix(os.Args[0], "go") {
		copy(os.Args[0], path.Join(dataDir, "polym1"))
		os.Chmod(path.Join(dataDir, "polym1"), 0755)
	}

	fmt.Println("downloading files...")
	zipPath := path.Join(dataDir, "files.zip")
	downloadFiles(zipPath)

	fmt.Println("unzipping...")
	unzip(zipPath, dataDir)
	os.Remove(zipPath)

	formatter := color.New(color.FgGreen)
	formatter.Println("done!")
	formatter.Print("add ")
	formatter.Print(color.New(color.Bold).Sprint(path.Join(dataDir, "polym1")))
	formatter.Println(" as your wrapper command and you're good to go!")
}
