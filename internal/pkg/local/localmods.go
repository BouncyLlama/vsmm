package local

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"vs-mm/internal/pkg"
)

func ListMods(moddir string) []*pkg.Modinfo {
	var mods []*pkg.Modinfo

	files, err := os.ReadDir(moddir)
	if err != nil {
		log.Println(err)
		return mods
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".zip") {
			zipfile, err := zip.OpenReader(moddir + file.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			insidefiles := zipfile.File
			for _, innerfile := range insidefiles {
				if innerfile.Name == "modinfo.json" {
					openedjsonfile, _ := innerfile.Open()
					bytes, _ := io.ReadAll(openedjsonfile)
					var m pkg.Modinfo
					json.Unmarshal(bytes, &m)
					m.FileName = file.Name()
					if !strings.HasPrefix(m.Version, "v") {
						m.Version = "v" + m.Version
					}

					//fmt.Println(m)
					mods = append(mods, &m)

				}
			}
		}

	}
	return mods
}
