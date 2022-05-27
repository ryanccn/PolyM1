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
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	oldCmd := os.Args[1:]

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

	if classPathIdx == -1 {
		log.Fatal("No classpath found!")
	}
	if nativesDirIdx == -1 {
		log.Fatal("Couldn't find natives dir option!")
	}

	newClassPath := strings.ReplaceAll(originalClassPath, "abc", "DEF")

	newCmd := oldCmd[:]
	newCmd[nativesDirIdx] = "-Djava.library.path="
	newCmd[classPathIdx] = newClassPath

	finalExec := exec.Command(newCmd[0], newCmd[1:]...)

	// make the subprocess mostly passthrough
	finalExec.Stdin = nil
	finalExec.Stdout = os.Stdout
	finalExec.Stderr = os.Stderr

	err := finalExec.Run()

	if err != nil {
		log.Fatal(err)
	}
}
