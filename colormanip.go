package main

import (
	"fmt"
)

// Color map for predefined color names
var colorMap = map[string][3]byte{
	"red":    {0xFF, 0x00, 0x00},
	"green":  {0x00, 0xFF, 0x00},
	"blue":   {0x00, 0x00, 0xFF},
	"purple": {0x80, 0x00, 0x80},
	"orange": {0xFF, 0xA5, 0x00},
	"yellow": {0xFF, 0xFF, 0x00},
	"cyan":   {0x00, 0xFF, 0xFF},
	"white":  {0xFF, 0xFF, 0xFF},
}

// Helper function to calculate luminance (simple average of RGB values)
func calculateLuminance(r, g, b byte) float64 {
	return (float64(r) + float64(g) + float64(b)) / 3.0
}

// Generalized function to transform a color to a target color while maintaining luminance
func adjustToTargetColor(r, g, b byte, targetColor [3]byte) (byte, byte, byte) {
	// Calculate the luminance of the current color
	luminance := calculateLuminance(r, g, b)

	// Scale the target color by the luminance
	newR := byte(float64(targetColor[0]) * luminance / 255.0)
	newG := byte(float64(targetColor[1]) * luminance / 255.0)
	newB := byte(float64(targetColor[2]) * luminance / 255.0)

	return newR, newG, newB
}

// Get RGB by color name
func getRGBForColor(colorName string) ([3]byte, error) {
	if rgb, exists := colorMap[colorName]; exists {
		return rgb, nil
	}
	return [3]byte{}, fmt.Errorf("color %s not found", colorName)
}
