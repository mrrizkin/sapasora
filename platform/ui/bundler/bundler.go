package bundler

import (
	"strings"

	vitego "codeberg.org/mrrizkin/vite-go"
)

type Bundler struct {
	vite *vitego.Vite
}

func NewBundler() (*Bundler, error) {
	vite := vitego.NewVite().
		UseBuildDirectory("public/build").
		UseManifestFilename("manifest.json").
		UseHotFile("public/hot").
		CreateAssetPathsUsing(func(path string, secure bool) string {
			return strings.TrimPrefix(path, "public")
		})

	return &Bundler{
		vite: vite,
	}, nil
}

func (b *Bundler) GetVite() *vitego.Vite {
	return b.vite
}

func (b *Bundler) Entry(entries ...string) string {
	result, err := b.vite.Invoke(entries, "")
	if err != nil {
		return ""
	}
	return result
}

func (b *Bundler) ReactRefresh() string {
	if !b.vite.IsRunningHot() {
		return ""
	}

	result, err := b.vite.ReactRefresh()
	if err != nil {
		return ""
	}
	return result
}

func (b *Bundler) Asset(path string) (string, error) {
	return b.vite.Asset(path, "")
}

func (b *Bundler) Content(entry string) (string, error) {
	return b.vite.Content(entry, "")
}

func (b *Bundler) CspNonce() string {
	return b.vite.CspNonce()
}
