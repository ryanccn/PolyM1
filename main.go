/*
PolyM1 - M1 hack for PolyMC
Copyright (C) 2022 Ryan Cao

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/ryanccn/PolyM1/install"
)

func main() {
	if len(os.Args) == 1 {
		os.Exit(0)
	}

	oldCmd := os.Args[1:]

	if oldCmd[0] == "install" {
		install.Install()
		return
	}

	envClassPath, isEnv := os.LookupEnv("CLASSPATH")

	if isEnv {
		fmt.Println("[polym1] using classpath from environment variable")
	} else {
		fmt.Println("[polym1] using classpath from option -cp")
	}

	classPathIdx := -1
	originalClassPath := ""
	nativesDirIdx := -1

	for idx, v := range oldCmd {
		if v == "-cp" {
			classPathIdx = idx + 1
			originalClassPath = oldCmd[classPathIdx]
		}

		if strings.HasPrefix(v, "-Djava.library.path=") {
			nativesDirIdx = idx
		}
	}

	if classPathIdx == -1 && !isEnv {
		log.Fatal("[polym1] no classpath found!")
	}
	if nativesDirIdx == -1 {
		log.Fatal("[polym1] couldn't find natives dir option!")
	}

	newClassPath := make([]string, 0)
	oldClassPath := originalClassPath
	if isEnv {
		oldClassPath = envClassPath
	}

	for _, val := range strings.Split(oldClassPath, ":") {
		if !strings.Contains(val, "lwjgl") && !strings.Contains(val, "java-objc-bridge") {
			newClassPath = append(newClassPath, val)
		} else {
			fmt.Println("[polym1] removed library", val)
		}
	}

	err := filepath.Walk(path.Join(install.GetDataDir(), "libraries"), func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".jar") {
			newClassPath = append(newClassPath, path)
			fmt.Println("[polym1] added library", path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	newCmd := oldCmd[:]
	newCmd[nativesDirIdx] = "-Djava.library.path=" + path.Join(install.GetDataDir(), "natives")
	if !isEnv {
		newCmd[classPathIdx] = strings.Join(newClassPath, ":")
	}

	fmt.Println("[polym1] patched command:", newCmd)

	finalExec := exec.Command(newCmd[0], newCmd[1:]...)

	// make the subprocess fully passthrough
	finalExec.Stdin = os.Stdin
	finalExec.Stdout = os.Stdout
	finalExec.Stderr = os.Stderr

	if isEnv {
		finalExec.Env = append(os.Environ(), "CLASSPATH="+strings.Join(newClassPath, ":"))
	}

	err = finalExec.Run()

	if err != nil {
		log.Fatal(err)
	}
}
