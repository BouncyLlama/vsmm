package util

import (
	"github.com/hashicorp/go-version"
	"log"
	"vs-mm/internal/pkg"
)

func supportsCurrentGameVersion(currentgameversion string, v pkg.Modversion) bool {
	vcur, err := version.NewVersion(currentgameversion)
	if err != nil {
		log.Fatal(err)

	}
	for _, sv := range v.SupportedGameVersions {
		supported, err2 := version.NewVersion(sv)
		if err2 != nil {
			log.Fatal(err2)
		}
		if supported.Equal(vcur) {
			return true
		} else if supported.Segments()[0] == vcur.Segments()[0] && supported.Segments()[1] == vcur.Segments()[1] {

			return true
		}
	}

	return false
}
