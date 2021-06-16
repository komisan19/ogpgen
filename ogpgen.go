package main

import (
	"flag"
	"fmt"

	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
  "golang.org/x/image/math/fixed"
)


// Create image size
const fontsize = 30
const imageWidth = 1200
const imageHeight = 630

func flagUsage() {
	usageText := `ogpgen is an OGP generater cli tool.
  Usage:

  ogpgen command [arguments]
    --url string      Target url 
    --filename string   output filename(ex: hoge.png...)

  Use "ogpgen [command] --help" for more information about a command`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func queryUrl(url *string) string {
	doc, err := goquery.NewDocument(*url)
	if err != nil {
		return "Invalid URL"
	}
	title := doc.Find("title").Text()
	return title
}

func backgroud(img *image.RGBA) {
	rect := img.Rect
	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for v := rect.Min.X; v < rect.Max.X; v++ {
			img.Set(v, h, color.RGBA{
				0,
				2,
				85,
				255,
			})
		}
	}
}

func createDraw(title string, fileName *string) {
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		fmt.Println("font", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	backgroud(img)

	opt := truetype.Options{Size: fontsize}
	face := truetype.NewFace(ft, &opt)
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: face,
	}
	d.Dot.X = (fixed.I(1200) - d.MeasureString(title)) / 2
	d.Dot.Y = fixed.I(315+int(fontsize/2))
	d.DrawString(title)

	file, err := os.Create(*fileName)
	defer file.Close()
	if err != nil {
		fmt.Println("error:file\n", err)
	}

	if err := png.Encode(file, img); err != nil {
		fmt.Println("error:png\n", err)
		return
	}
}

func main() {
	flag.Usage = flagUsage
	url := flag.String("url", " ", "Target url")
	fileName := flag.String("filename", "example.png", "Output file")
	flag.Parse()
	title := queryUrl(url)
	createDraw(title, fileName)

	if len(os.Args) == 1 {
		flagUsage()
		return
	}
}
