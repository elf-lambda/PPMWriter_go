package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
}

func (pw *PPMWriter) setHeader() {
	// Set the header for the image
	pw.header = "P6\n" +
		strconv.Itoa(pw.imageWidth) + " " +
		strconv.Itoa(pw.imageHeight) + "\n255\n"
}

func (pw *PPMWriter) initImageArray() {
	pw.ppmimage = make([]byte, (pw.imageWidth * pw.imageHeight * 3))
}

func isPrintable(c byte) bool {
	// Check if char is printable
	return c >= 32 && c <= 126
}

func (pw *PPMWriter) save() {
	// Save the image to disk

	of, err := os.Create(pw.outputFile)
	check(err)
	defer of.Close()

	fmt.Fprintf(of, "%s", pw.header)
	of.Write(pw.ppmimage)
}

func (pw *PPMWriter) writerCharToArray(c byte, startX, startY int) {

	// Iterate over the character pixel data
	for y := 0; y < CHAR_IMAGE_HEIGHT; y++ {
		for x := 0; x < CHAR_IMAGE_WIDTH; x++ {
			imgX := startX + x
			imgY := startY + y
			charIndex := y*CHAR_IMAGE_WIDTH + x
			pixelIndex := (imgY*pw.imageWidth + imgX) * 3
			color := fontCharacters[c][charIndex]

			pw.ppmimage[pixelIndex] = byte((color >> 16) & 0xFF)  // Red
			pw.ppmimage[pixelIndex+1] = byte((color >> 8) & 0xFF) // Green
			pw.ppmimage[pixelIndex+2] = byte(color & 0xFF)        // Blue
		}
	}
}

func (pw *PPMWriter) writePPMImageArray() {
	// Write the input to the ppm image array
	startX := pw.padding
	startY := pw.padding
	var c byte // String char
	// Iterate over each char of the string, including \n
	for i := range len(pw.data.data) {
		c = pw.data.data[i] // Currenct char in string
		if c == '\n' {
			startX = pw.padding
			startY += CHAR_IMAGE_HEIGHT + 1
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
