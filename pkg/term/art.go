package term

import (
	"fmt"
	"image/color"
	"io"
	"path/filepath"
	"strings"

	"github.com/batmac/ccat/pkg/log"

	"github.com/eliukblau/pixterm/pkg/ansimage"
)

const (
	artScale = 0.5
)

var (
	ext = []string{
		"jpeg", "jpg", "gif", "png", "heic", "tiff", "svg",
	}
	extMap = make(map[string]bool)
)

func init() {
	for _, e := range ext {
		extMap["."+e] = true
	}
}

func IsArt(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	_, ok := extMap[extension]

	log.Debugf("testing %v(%v) IsArt? %v", path, extension, ok)

	return ok
}

func PrintArt(r io.Reader) {
	// fmt.Println()
	if IsITerm2() {
		log.Debugf("  art: printITerm2Art\n")
		_ = PrintITerm2Art(r)
	} else {
		log.Debugf("  art: printANSIArt\n")
		_ = PrintANSIArt(r)
	}
	fmt.Println()
}

/* func PrintArtFromURL(url string) {
	fmt.Println()
	if IsITerm2() {
		printITerm2ArtFromURL(url)
	} else {
		printANSIArtFromURL(url)
	}
	fmt.Println()

} */
func PrintANSIArt(r io.Reader) error {
	tx, ty, err := GetTerminalSize()
	if err != nil {
		log.Println(err)
	}

	sfy, sfx := 2, 1 // 2x1 --> without dithering

	img, err := ansimage.NewScaledFromReader(
		r,
		int(float32(ty*sfy)*artScale), int(float32(tx*sfx)*artScale),
		color.Transparent,
		ansimage.ScaleModeFit,
		ansimage.NoDithering,
	)
	if err != nil {
		log.Println(err)
	}
	img.Draw()
	return nil
}

/* func printANSIArtFromURL(url string) {

tx, ty, err := GetTerminalSize()
if err != nil {
	log.Println(err)
}
// set image scale factor for ANSIPixel grid
//sfy, sfx := ansimage.BlockSizeY, ansimage.BlockSizeX // 8x4 --> with dithering
//if ansimage.DitheringMode(flagDither) == ansimage.NoDithering {
sfy, sfx := 2, 1 // 2x1 --> without dithering
//}

img, err := ansimage.NewScaledFromURL(
	url,
	int(float32(ty*sfy)*artScale), int(float32(tx*sfx)*artScale),
	color.Transparent,
	ansimage.ScaleModeFit,
	ansimage.NoDithering,
)

if err != nil {
	log.Println(err)
}
img.Draw()
} */
