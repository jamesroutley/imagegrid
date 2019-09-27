package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

var (
	marginPercent  float64
	outputFilename string
	numCols        int
)

func init() {
	flag.Float64Var(&marginPercent, "margin", 5, "margin size (percentage)")
	flag.StringVar(
		&outputFilename, "output-filename", fmt.Sprintf("imagegrid-image-%s.png",
			time.Now().Format("2006-01-02")), "name of the file to save the image to",
	)
	flag.IntVar(&numCols, "cols", -1, "Number of columns of images. Negative number or zero indicates puts all images on the same row")

	flag.Parse()
}

func main() {
	files := flag.Args()

	if numCols <= 0 {
		numCols = len(files)
	}

	if err := run(files); err != nil {
		log.Fatal(err)
	}
}

func run(files []string) error {
	images := make([]image.Image, len(files))
	for i, filename := range files {
		image, err := decodeImageFile(filename)
		if err != nil {
			return err
		}
		images[i] = image
	}

	margin := calculateMargin(images)

	numRows := int(math.Ceil(float64(len(images)) / float64(numCols)))
	imageGroups := make([][]image.Image, numRows)
	row := 0
	col := 0
	for _, img := range images {
		imageGroups[row] = append(imageGroups[row], img)
		col++

		if col == numCols {
			col = 0
			row++
		}
	}

	yMax := height(imageGroups, margin)
	xMax := width(imageGroups, margin)

	outputImg := image.NewRGBA64(image.Rect(0, 0, xMax, yMax))

	offsets := calculateOffsets(imageGroups, margin)

	for i, img := range images {
		insertImage(outputImg, offsets[i], img)
	}

	// Output a PNG because it supports transparency
	f, err := os.Create(outputFilename)
	if err != nil {
		return err
	}

	if err := png.Encode(f, outputImg); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func decodeImageFile(filename string) (image.Image, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func calculateMargin(images []image.Image) (margin int) {
	height := maxHeight(images)
	return int(float64(height) * marginPercent * 0.01)
}

func height(imageGroups [][]image.Image, margin int) (max int) {
	for _, group := range imageGroups {
		max += maxHeight(group)
	}
	max += margin * (len(imageGroups) - 1)
	return max
}

func maxHeight(images []image.Image) (max int) {
	max = 0
	for _, image := range images {
		if image.Bounds().Min.Y != 0 {
			// Sanity check - all images are well formed, and
			// img.Bounds(),Min.Y == 0
			log.Fatal("image min y does not equal 0")
		}
		height := image.Bounds().Max.Y
		if height > max {
			max = height
		}
	}
	return max
}

func width(imageGroups [][]image.Image, margin int) int {
	largestWidth := 0
	for _, group := range imageGroups {
		width := sumWidths(group)
		width += margin * (len(group) - 1)
		if width > largestWidth {
			largestWidth = width
		}
	}

	return largestWidth
}

func sumWidths(images []image.Image) (sum int) {
	sum = 0
	for _, img := range images {
		if img.Bounds().Min.X != 0 {
			// Sanity check - all images are well formed, and
			// img.Bounds(),Min.X == 0
			log.Fatal("image min x does not equal 0")
		}
		sum += img.Bounds().Max.X
	}
	return sum
}

func calculateOffsets(imageGroups [][]image.Image, margin int) (offsets []image.Point) {
	yOffset := 0
	xOffset := 0
	for _, row := range imageGroups {
		for _, img := range row {
			offset := image.Point{
				X: xOffset,
				Y: yOffset,
			}
			offsets = append(offsets, offset)
			xOffset += img.Bounds().Max.X + margin
		}
		xOffset = 0
		yOffset += maxHeight(row) + margin
	}
	return
}

func insertImage(outputImg *image.RGBA64, offset image.Point, img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			outputImg.SetRGBA64(x+offset.X, y+offset.Y, color.RGBA64{
				R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a),
			})
		}
	}
}
