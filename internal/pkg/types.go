package pkg

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-version"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Modinfo struct {
	Type              string `json:"type"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Modid             string `json:"modid"`
	FileName          string
	SelectedVersion   *Modversion
	LatestVersion     *Modversion
	AvailableVersions []*Modversion
}
type Modversion struct {
	DownloadLink          string
	SupportedGameVersions []string
	ModVersion            string
}

func (mi *Modinfo) GetMatchingVersion(v string) *Modversion {
	for _, x := range mi.AvailableVersions {
		if x.ModVersion == v {
			return x
		}
	}
	return nil
}
func (mi *Modinfo) ListAvailableStrings() []string {
	var strs []string
	for _, x := range mi.AvailableVersions {
		strs = append(strs, x.ModVersion)
	}
	return strs
}
func (mi *Modinfo) GetAvailableVerions() {
	mi.AvailableVersions = GetAvailableVersions(mi.Modid)
	mi.LatestVersion = mi.AvailableVersions[0]

}

func (m *Modinfo) UpdateToSelected(moddir string) {

	var sver *Modversion
	for _, v := range m.AvailableVersions {
		if v.ModVersion == m.SelectedVersion.ModVersion {
			sver = v
		}
	}
	resp, err := http.Get(sver.DownloadLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	os.Remove(moddir + m.Modid + ".zip")
	file, _ := os.Create(moddir + m.Modid + ".zip")
	defer file.Close()

	io.Copy(file, resp.Body)
	if m.FileName != m.Modid+".zip" {
		os.Remove(moddir + m.FileName)

	}
	m.FileName = file.Name()
	m.Version = m.SelectedVersion.ModVersion

}
func (m *Modinfo) UpdateAvailable() bool {
	vcur, err := version.NewVersion(m.Version)
	vlatest, err2 := version.NewVersion(m.LatestVersion.ModVersion)
	if err != nil {
		log.Fatal(err)

	}
	if err2 != nil {
		log.Fatal(err2)

	}
	return vlatest.GreaterThan(vcur)
}
func (mv *Modversion) FormatSupported() string {
	s := ""
	for _, v := range mv.SupportedGameVersions {
		s = s + v + ","
	}
	return s
}
func GetAvailableVersions(name string) []*Modversion {
	res, err := http.Get("https://mods.vintagestory.at/" + name + "#tab-files")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	availableVersions := []*Modversion{}
	downloadlink := "https://mods.vintagestory.at"
	// Find the review items
	doc.Find("#Connection\\ types > tbody > tr").Each(func(i int, s *goquery.Selection) {
		n := s.Find("td")
		nodes := n.Nodes
		var version = new(Modversion)
		version.ModVersion = strings.TrimSpace(nodes[0].FirstChild.Data)
		gameversionslist := n.Find(".tags").Find("a").Nodes
		for _, gv := range gameversionslist {
			version.SupportedGameVersions = append(version.SupportedGameVersions, gv.FirstChild.Data)
		}
		for _, a := range nodes[5].FirstChild.Attr {
			if a.Key == "href" {
				version.DownloadLink = downloadlink + a.Val
			}
		}
		availableVersions = append(availableVersions, version)

	})
	sort.SliceStable(availableVersions, func(i, j int) bool {
		vi, e1 := version.NewVersion(availableVersions[i].ModVersion)
		vj, e2 := version.NewVersion(availableVersions[j].ModVersion)
		if e1 != nil || e2 != nil {
			log.Println("cannot sort non-semver mod version")
			return false
		}
		return vi.GreaterThan(vj)
	})
	return availableVersions

}
