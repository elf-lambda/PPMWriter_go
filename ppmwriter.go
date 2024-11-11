package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type PPMData struct {
	data       string
	rows       int
	maxLineLen int
}

func (id *PPMData) read(path string) {
	// Read a file into the struct keeping track of rows and line length
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		id.data += line
		id.data += "\n"
		id.rows++

		tmp := len(line)
		if tmp > id.maxLineLen {
			id.maxLineLen = tmp
		}
	}
	fmt.Printf("Populated struct with %d rows and maxlen %d\n", id.rows, id.maxLineLen)
}

type PPMWriter struct {
	header      string
	imageWidth  int
	imageHeight int
	padding     int
	outputFile  string
	data        PPMData
	ppmimage    []byte
}

func (pw *PPMWriter) init(data PPMData) {
	pw.padding = 15
	pw.outputFile = "output.ppm"
	pw.data = data
	pw.setImageSize()
	pw.setHeader()
	pw.initImageArray()
}

func (pw *PPMWriter) setImageSize() {
	// Calculate the image width and height based on the input
	pw.imageWidth = pw.data.maxLineLen*CHAR_IMAGE_WIDTH + (pw.padding * 2)
	pw.imageHeight = pw.data.rows*CHAR_IMAGE_HEIGHT + (pw.padding * 2)

	lines := strings.Split(pw.data.data, "\n")
	for _, line := range lines {
		// for each line
		count := strings.Count(line, "![pepe1]")
		if count > 0 {
			if len(line) > pw.data.maxLineLen {
				pw.imageWidth += 32 * count
			}
			pw.imageHeight += 96
		}
	}

}

func (pw *PPMWriter) setHeader() {
	// Set the header for the image
	pw.header = "P6\n" +
		strconv.Itoa(pw.imageWidth) + " " +
		strconv.Itoa(pw.imageHeight) + "\n255\n"
}

func (pw *PPMWriter) initImageArray() {
	pw.ppmimage = make([]byte, (pw.imageWidth * pw.imageHeight * 3)) // wasting memory
}

func (pw *PPMWriter) save() {
	// Save the image to disk

	of, err := os.Create(pw.outputFile)
	check(err)
	defer of.Close()

	fmt.Fprintf(of, "%s", pw.header)
	of.Write(pw.ppmimage)
}

func (pw *PPMWriter) writeImageToArray(c byte, startX, startY int) {
	// TODO: to be implemented fully
	fontChar := fontCharacters[c]
	// Iterate over the character pixel data
	for y := 0; y < fontChar.height; y++ {
		for x := 0; x < fontChar.width; x++ {
			imgX := startX + x
			imgY := startY + y
			charIndex := y*fontChar.width + x
			pixelIndex := (imgY*pw.imageWidth + imgX) * 3
			color := fontChar.data[charIndex]
			r := byte((color >> 16) & 0xFF) // red
			g := byte((color >> 8) & 0xFF)  // green
			b := byte(color & 0xFF)         // blue

			pw.ppmimage[pixelIndex] = r
			pw.ppmimage[pixelIndex+1] = g
			pw.ppmimage[pixelIndex+2] = b
		}
	}
}

func (pw *PPMWriter) writerCharToArray(c byte, startX, startY int) {
	fontChar := fontCharacters[c]

	// Iterate over the character pixel data
	for y := 0; y < fontChar.height; y++ {
		for x := 0; x < fontChar.width; x++ {
			imgX := startX + x
			imgY := startY + y
			charIndex := y*fontChar.width + x
			pixelIndex := (imgY*pw.imageWidth + imgX) * 3
			color := fontChar.data[charIndex]
			r := byte((color >> 16) & 0xFF) // red
			g := byte((color >> 8) & 0xFF)  // green
			b := byte(color & 0xFF)         // blue

			selectedColor := "green" // Change to the desired color
			targetRGB, err := getRGBForColor(selectedColor)
			check(err)

			// If the color is not black (0x000000), transform it to green while maintaining luminance
			if color != 0x000000 {
				newR, newG, newB := adjustToTargetColor(r, g, b, targetRGB)
				pw.ppmimage[pixelIndex] = newR
				pw.ppmimage[pixelIndex+1] = newG
				pw.ppmimage[pixelIndex+2] = newB
			} else {
				// Keep black as is
				pw.ppmimage[pixelIndex] = r
				pw.ppmimage[pixelIndex+1] = g
				pw.ppmimage[pixelIndex+2] = b
			}
		}
	}
}

func isPrintable(c byte) bool {
	// Check if char is printable
	return c >= 32 && c <= 126
}

func (pw *PPMWriter) writePPMImageArray() {
	// Write the input to the ppm image array
	startX := pw.padding
	startY := pw.padding
	var c byte // Current char

	// Iterate over each char of the string, including \n
	imagePattern := []byte("![pepe1]") // The pattern to look for
	imagePatternIndex := 0             // Index to track progress through the pattern

	for i := 0; i < len(pw.data.data); i++ {
		c = pw.data.data[i] // Current char in string

		// Handle the pattern matching for ![image]
		if c == imagePattern[imagePatternIndex] {
			// Move to the next character in the pattern
			imagePatternIndex++
			if imagePatternIndex == len(imagePattern) {
				// Once the whole pattern is matched, write the image and reset
				pw.writeImageToArray(2, startX, startY)
				startX += fontCharacters[2].width
				imagePatternIndex = 0
				continue
			}
			continue
		} else {
			// Reset if the pattern is broken
			imagePatternIndex = 0
		}

		// Handle normal characters and newline
		if c == '\n' {
			startX = pw.padding
			startY += CHAR_IMAGE_HEIGHT
			continue
		}
		if isPrintable(c) {
			pw.writerCharToArray(c, startX, startY)
			startX += CHAR_IMAGE_WIDTH
		} else {
			pw.writerCharToArray(1, startX, startY)
			startX += CHAR_IMAGE_WIDTH
		}
	}
}
