package install

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/cavaliergopher/grab/v3"
)

func download(dest string, src string, desc string) string {
	grabClient := grab.NewClient()
	grabReq, err := grab.NewRequest(dest, src)

	if err != nil {
		log.Fatal(err)
	}

	grabRes := grabClient.Do(grabReq)

	<-grabRes.Done
	if grabRes.Err() != nil {
		log.Fatal(grabRes.Err())
	}

	fmt.Printf("downloaded %s to %s\n", desc, dest)

	return grabRes.Filename
}

var META_URLS = []string{"https://repo1.maven.org/maven2/ca/weblite/java-objc-bridge/maven-metadata.xml", "https://repo1.maven.org/maven2/org/lwjgl/lwjgl/maven-metadata.xml"}
var LWJGL_LIBS = []string{"lwjgl", "lwjgl-glfw", "lwjgl-jemalloc", "lwjgl-tinyfd", "lwjgl-stb", "lwjgl-opengl", "lwjgl-openal"}

// heavily inspired by https://github.com/Dreamail/M1MC/blob/master/main.go
func DownloadFiles() {
	dataDir := GetDataDir()
	tmpDir, err := os.MkdirTemp(os.TempDir(), "polym1")

	var wg sync.WaitGroup

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range []string{"libraries", "natives"} {
		err := os.Mkdir(path.Join(dataDir, v), 0770)

		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				log.Fatal(err)
			} else {
				os.RemoveAll(path.Join(dataDir, v))
				err := os.Mkdir(path.Join(dataDir, v), 0770)

				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	lwjglVersion := ""
	objcBridgeVersion := ""

	for _, v := range META_URLS {
		resp, err := http.Get(v)
		if err != nil {
			log.Fatal(err)
		}

		xmlbytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		xmlstr := string(xmlbytes)
		versionIndex := strings.Index(xmlstr, "<latest>") + len("<latest>")
		versionOutdex := strings.Index(xmlstr, "</latest>")

		version := string([]rune(xmlstr)[versionIndex:versionOutdex])
		if strings.Contains(v, "lwjgl") {
			lwjglVersion = version
		} else {
			objcBridgeVersion = version
		}
	}

	for _, v := range LWJGL_LIBS { // process lwjgl
		jarUrl := "https://repo1.maven.org/maven2/org/lwjgl/" + v + "/" + lwjglVersion + "/" + v + "-" + lwjglVersion + ".jar"
		nativeUrl := "https://repo1.maven.org/maven2/org/lwjgl/" + v + "/" + lwjglVersion + "/" + v + "-" + lwjglVersion + "-natives-macos-arm64.jar"

		wg.Add(1)
		go func() {
			defer wg.Done()

			download(path.Join(dataDir, "libraries"), jarUrl, v)

			if err != nil {
				log.Fatal(err)
			}
		}()

		wg.Add(1)
		go func(v string) {
			defer wg.Done()

			downloadFile := download(tmpDir, nativeUrl, v+"-natives")

			jarZip, err := zip.OpenReader(downloadFile)
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range jarZip.File {
				if strings.Contains(v.Name, "dylib") {
					dylibZip, err := v.Open()
					if err != nil {
						log.Fatal(err)
					}

					dylib, err := os.Create(path.Join(dataDir, "natives", strings.Split(v.Name, "/")[strings.Count(v.Name, "/")]))
					if err != nil {
						log.Fatal(err)
					}

					io.Copy(dylib, dylibZip)
					dylibZip.Close()
					dylib.Close()

					break
				}
			}

			jarZip.Close()
		}(v)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		jarUrl := "https://repo1.maven.org/maven2/ca/weblite/java-objc-bridge/" + objcBridgeVersion + "/" + "java-objc-bridge-" + objcBridgeVersion + ".jar"

		filename := download(path.Join(dataDir, "libraries"), jarUrl, "java-objc-bridge")

		jarZip, err := zip.OpenReader(filename)
		if err != nil {
			log.Fatal(err)
		}

		dylib, err := os.Create(path.Join(dataDir, "natives", "libjcocoa.dylib"))
		if err != nil {
			log.Fatal(err)
		}

		dylibZip, err := jarZip.Open("libjcocoa.dylib")

		if err != nil {
			log.Fatal(err)
		}

		io.Copy(dylib, dylibZip)
		jarZip.Close()
		dylib.Close()
	}()

	wg.Wait()
}
