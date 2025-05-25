package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/joho/godotenv"

	arabic "github.com/abdullahdiaa/garabic"
)

func AddIDToImage(filePath string, id string, signature string) error {
	_ = godotenv.Load()

	file, err := os.Open(filePath)
	if err != nil {
		// return fmt.Errorf("failed to open image: %w", err)
		return HandleError(err, "Failed to open image", Error)
	}
	defer file.Close()

	var img image.Image
	if strings.HasSuffix(strings.ToLower(filePath), ".png") {
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		// return fmt.Errorf("failed to decode image: %w", err)
		return HandleError(err, "Failed to decode image", Error)
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	dc := gg.NewContext(w, h)
	dc.DrawImage(img, 0, 0)

	textPrefix := os.Getenv("IMAGE_TEXT_PREFIX")
	if textPrefix == "" {
		textPrefix = "Document ID:"
	}

	rawText := fmt.Sprintf("%s %s", textPrefix, id)
	text := arabic.Shape(rawText)

	fontPath := os.Getenv("IMAGE_FONT_PATH")
	if fontPath == "" {
		fontPath = "assets/fonts/Cairo.ttf"
	}

	fontSize, _ := strconv.ParseFloat(os.Getenv("IMAGE_FONT_SIZE"), 64)
	if fontSize == 0 {
		fontSize = 18
	}

	textColor := parseRGB(os.Getenv("IMAGE_TEXT_COLOR"), 255, 255, 255)
	bgColor := parseRGB(os.Getenv("IMAGE_BG_COLOR"), 0, 0, 0)
	bgOpacity := parseOpacity(os.Getenv("IMAGE_BG_OPACITY"))

	align := strings.ToLower(os.Getenv("IMAGE_TEXT_ALIGN"))
	var anchorX float64
	switch align {
	case "left":
		anchorX = 0
	case "right":
		anchorX = 1
	default:
		anchorX = 0.5
	}

	boxHeight := fontSize + 20
	dc.SetRGBA(bgColor[0], bgColor[1], bgColor[2], bgOpacity)
	dc.DrawRectangle(0, float64(h)-boxHeight, float64(w), boxHeight)
	dc.Fill()

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		// return fmt.Errorf("failed to load font: %w", err)
		return HandleError(err, "Failed to load font", Error)
	}
	dc.SetRGB(textColor[0], textColor[1], textColor[2])
	x := float64(w) * anchorX
	y := float64(h) - boxHeight/2
	dc.DrawStringAnchored(text, x, y, anchorX, 0.5)

	generateQR := strings.ToLower(os.Getenv("QR_GENERATOR")) == "true"
	if generateQR {
		qrBytes, err := GenerateQRCodeImage(id)
		if err == nil {
			qrImg, err := png.Decode(bytes.NewReader(qrBytes))
			if err == nil {
				qrSize := 100.0
				marginX := parseFloatEnv("QR_MARGIN_X", 10)
				marginY := parseFloatEnv("QR_MARGIN_Y", 10)
				position := strings.ToLower(os.Getenv("QR_POSITION"))

				var qrX, qrY float64

				switch position {
				case "top-left":
					qrX = marginX
					qrY = marginY
				case "top-right":
					qrX = float64(w) - qrSize - marginX
					qrY = marginY
				case "bottom-left":
					qrX = marginX
					qrY = float64(h) - qrSize - marginY
				default: // "bottom-right"
					qrX = float64(w) - qrSize - marginX
					qrY = float64(h) - qrSize - marginY
				}

				dc.DrawImageAnchored(qrImg, int(qrX), int(qrY), 0, 0)
			}
		}
	}

	out, err := os.Create(filePath)
	if err != nil {
		// return fmt.Errorf("failed to create output image: %w", err)
		return HandleError(err, "Failed to create output image", Error)
	}
	defer out.Close()

	if strings.HasSuffix(strings.ToLower(filePath), ".png") {
		err = png.Encode(out, dc.Image())
	} else {
		err = jpeg.Encode(out, dc.Image(), &jpeg.Options{Quality: 90})
	}
	if err != nil {
		// return fmt.Errorf("failed to encode output image: %w", err)
		return HandleError(err, "Failed to encode output image", Error)
	}

	comment := fmt.Sprintf("ID:%s;SIG:%s", id, signature)
	cmd := exec.Command("exiftool", "-overwrite_original", "-UserComment="+comment, filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// return fmt.Errorf("failed to write Exif: %v | Output: %s", err, string(output))
		return HandleError(err, fmt.Sprintf("Failed to write Exif: %s", string(output)), Error)
	}

	return nil
}

func parseRGB(input string, defR, defG, defB int) [3]float64 {
	parts := strings.Split(input, ",")
	if len(parts) != 3 {
		return [3]float64{float64(defR) / 255, float64(defG) / 255, float64(defB) / 255}
	}
	r, _ := strconv.Atoi(parts[0])
	g, _ := strconv.Atoi(parts[1])
	b, _ := strconv.Atoi(parts[2])
	return [3]float64{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

func parseOpacity(input string) float64 {
	val, err := strconv.ParseFloat(input, 64)
	if err != nil || val < 0 || val > 1 {
		return 0.5
	}
	return val
}

func parseFloatEnv(key string, def float64) float64 {
	val, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		return def
	}
	return val
}
