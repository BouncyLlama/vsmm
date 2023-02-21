package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strings"
	"vs-mm/internal/pkg"
	"vs-mm/internal/pkg/local"
)

func listmods(moddir string) []*pkg.Modinfo {
	data := local.ListMods(moddir)
	for _, row := range data {
		row.GetAvailableVerions()
		for _, v := range row.AvailableVersions {
			if v.ModVersion == row.Version {
				row.SelectedVersion = v
			}
		}
	}
	return data
}
func main() {
	myApp := app.NewWithID("com.bouncyllama.vsmm")
	myWindow := myApp.NewWindow("Vintage Story Mod Manager")
	moddir := myApp.Preferences().StringWithFallback("vsmoddir", "nowhere")

	var data []*pkg.Modinfo

	vert := container.NewVBox()
	//makeRow(data, vert)
	infinite := widget.NewProgressBarInfinite()
	infinite.Hide()
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			infinite.Show()
			data = listmods(moddir)

			makeRow(data, vert)
			vert.Refresh()
			infinite.Hide()
		}),
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
				path := uri.Path()
				if !strings.HasSuffix(path, "/") {
					path = path + "/"
				}
				myApp.Preferences().SetString("vsmoddir", path)
				moddir = myApp.Preferences().String("vsmoddir")
				infinite.Show()
				data = listmods(moddir)

				makeRow(data, vert)
				vert.Refresh()
				infinite.Hide()
			}, myWindow)

		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			infinite.Show()
			for _, d := range data {
				if d.Version != d.SelectedVersion.ModVersion {
					d.UpdateToSelected(moddir)
					//ht.Refresh()
				}
			}
			vert.RemoveAll()

			data = listmods(moddir)

			makeRow(data, vert)
			vert.Refresh()
			infinite.Hide()
		}),
	)

	content := container.NewBorder(toolbar, infinite, nil, nil, container.NewVScroll(vert))

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func makeRow(data []*pkg.Modinfo, vert *fyne.Container) {
	for _, d := range data {
		name := widget.NewLabel(d.Name)
		installedVer := widget.NewLabel(d.Version)
		sver := widget.NewSelect(d.SelectedVersion.SupportedGameVersions, func(s string) {

		})
		sver.PlaceHolder = "Supported Game Versions"
		selectver := widget.NewSelect(d.ListAvailableStrings(), func(s string) {
			d.SelectedVersion = d.GetMatchingVersion(s)
			sver.Options = d.SelectedVersion.SupportedGameVersions
			sver.Selected = sver.Options[0]
			sver.Refresh()
		})
		selectver.PlaceHolder = "Desired Mod Version"
		c := container.NewHBox(name, installedVer, sver, selectver)
		vert.Add(c)
	}
}
