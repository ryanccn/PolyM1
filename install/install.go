package install

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/fatih/color"
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

	err := os.RemoveAll(path.Join(dataDir, "libraries"))
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll(path.Join(dataDir, "natives"))
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("copying binary...")
	// if !strings.HasSuffix(os.Args[0], "go") {
	// 	copy(os.Args[0], path.Join(dataDir, "polym1"))
	// 	os.Chmod(path.Join(dataDir, "polym1"), 0755)
	// }

	fmt.Println("downloading files...")
	DownloadFiles()

	formatter := color.New(color.FgGreen)
	formatter.Println("done!")
	formatter.Print("add ")
	formatter.Print(color.New(color.Bold).Sprint(os.Args[0]))
	formatter.Println(" as your wrapper command and you're good to go!")
}
