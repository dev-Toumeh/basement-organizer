package database

import (
	"basement/main/internal/logg"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"

	"golang.org/x/image/draw"
)

var UnsupportedImageFormat error = errors.New("Unsupported picture format")

// updatePicture checks for valid base64 encoding and creates a resized preview image
// in `previewPicture` with max side length of 50 pixel.
// In case of error the strings will be set to empty string.
func updatePicture(picture *string, previewPicture *string) error {
	if *picture != "" {
		_, err := Base64StringToByte(*picture)
		if err != nil {
			*picture = ""
			*previewPicture = ""
			return logg.Errorf("CreateShelf: invalid base64 in Picture %w", err)
		}

		*previewPicture, err = ResizePNG(*picture, 50)
		if err != nil {
			*previewPicture = ""
			return logg.WrapErr(err)
		}
	}
	return nil
}

// ResizePNG2 resizes image while keeping aspect ratio.
// Longest side will fit the pixel value of `fitLongestSideToPixel`.
func ResizePNG(input64 string, fitLongestSideToPixel int) (string, error) {
	pic, err := Base64StringToByte(input64)
	if err != nil {
		return "", logg.WrapErr(err)
	}

	picreader := bytes.NewReader(pic)
	var output []byte
	buf := bytes.NewBuffer(output)

	err = ResizePNG2(picreader, buf, fitLongestSideToPixel)
	if err != nil {
		return "", logg.WrapErr(err)
	}
	return ByteToBase64String(buf.Bytes()), nil
}

// ResizePNG2 resizes image while keeping aspect ratio.
// Longest side will fit the pixel value of `fitLongestSideToPixel`.
func ResizePNG2(input io.Reader, output io.Writer, fitLongestSideToPixel int) error {
	// Decode the image (from PNG to image.Image):
	src, err := png.Decode(input)
	if err != nil {
		return logg.WrapErr(err)
	}

	// Set the expected size that you want:
	var scale float64
	var newX, newY int
	if src.Bounds().Max.X > src.Bounds().Max.Y {
		scale = float64(src.Bounds().Max.X) / float64(fitLongestSideToPixel)
	} else {
		scale = float64(src.Bounds().Max.Y) / float64(fitLongestSideToPixel)
	}
	newX = int(math.Round(float64(src.Bounds().Max.X) / scale))
	newY = int(math.Round(float64(src.Bounds().Max.Y) / scale))

	dst := image.NewRGBA(image.Rect(0, 0, newX, newY))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:
	err = png.Encode(output, dst)
	if err != nil {
		return logg.WrapErr(err)
	}

	return nil
}

func ResizeImage(input64 string, fitLongestSideToPixel int, format string) (string, error) {
	switch format {
	case "image/jpeg":
		return ResizeJPEG(input64, 50)
	case "image/png":
		return ResizePNG(input64, fitLongestSideToPixel)
	case "":
		logg.Debug("remove image")
		return "", nil
	default:
		// logg.Warningf("%s", UnsupportedImageFormat.Error())
		return "", logg.WrapErr(fmt.Errorf("%w %s", UnsupportedImageFormat, format))
	}
}

// ResizeJPEG resizes a JPEG image while keeping its aspect ratio.
// The longest side will fit the pixel value of `fitLongestSideToPixel`.
func ResizeJPEG(input64 string, fitLongestSideToPixel int) (string, error) {
	pic, err := Base64StringToByte(input64)
	if err != nil {
		return "", logg.WrapErr(err)
	}

	picreader := bytes.NewReader(pic)
	var output []byte
	buf := bytes.NewBuffer(output)

	err = ResizeJPEG2(picreader, buf, fitLongestSideToPixel)
	if err != nil {
		return "", logg.WrapErr(err)
	}
	return ByteToBase64String(buf.Bytes()), nil
}

// ResizeJPEG2 resizes a JPEG image (read from `input`) while keeping its aspect ratio.
// The longest side will fit the pixel value of `fitLongestSideToPixel`.
func ResizeJPEG2(input io.Reader, output io.Writer, fitLongestSideToPixel int) error {
	// Decode the image (from JPEG to image.Image):
	src, err := jpeg.Decode(input)
	if err != nil {
		return logg.WrapErr(err)
	}

	// Calculate scale and new dimensions:
	var scale float64
	var newX, newY int
	if src.Bounds().Max.X > src.Bounds().Max.Y {
		scale = float64(src.Bounds().Max.X) / float64(fitLongestSideToPixel)
	} else {
		scale = float64(src.Bounds().Max.Y) / float64(fitLongestSideToPixel)
	}
	newX = int(math.Round(float64(src.Bounds().Max.X) / scale))
	newY = int(math.Round(float64(src.Bounds().Max.Y) / scale))

	dst := image.NewRGBA(image.Rect(0, 0, newX, newY))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output` in JPEG format:
	// You can pass &jpeg.Options{Quality: 80} (or another value) for quality customization.
	err = jpeg.Encode(output, dst, nil)
	if err != nil {
		return logg.WrapErr(err)
	}

	return nil
}

func Base64StringToByte(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func ByteToBase64String(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
