package main

import "os/exec"

func main() {
	t := PPMData{}
	t.read("input.txt")
	writer := PPMWriter{}
	writer.init(t)
	writer.writePPMImageArray()
	writer.save()
	convToPng("output.ppm", "out.png")
}

func convToPng(in, out string) {
	cmd := exec.Command("ffmpeg", "-y", "-i", in, out)
	err := cmd.Run()
	check(err)
}
