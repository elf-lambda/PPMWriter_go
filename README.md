![Alt text](screenshots/out.png)

# PPMWriter

- Used to write text to a ppmimage
- Only supports ascii characters defined in mfont.go, any other codepoint is replaced with the inverted "?"
- No external FONTS/Libraries used, only depends on stdlib

## How to run

> go run .

## Supports displaying a pepe image

![Alt text](screenshots/pepe.png)

## Also supports colored text

![Alt text](screenshots/green.png)
![Alt text](screenshots/cyan.png)

## Convert .ppm to .png/.jpeg

Depends on ffmpeg
