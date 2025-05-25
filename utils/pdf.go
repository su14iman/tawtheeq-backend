package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
	"github.com/signintech/gopdf"
)

func AddIDToPDF(filePath string, id string, signature string) error {

	doc, err := fitz.New(filePath)
	if err != nil {
		// return fmt.Errorf("failed to open PDF: %w", err)
		return HandleError(err, "Failed to open PDF", Error)
	}
	defer doc.Close()

	tempDir := "./temp_pdf_pages"
	os.MkdirAll(tempDir, os.ModePerm)

	imagePaths := []string{}

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			// return fmt.Errorf("failed to render page %d: %w", n+1, err)
			return HandleError(err, fmt.Sprintf("Failed to render page %d", n+1), Error)
		}

		imgPath := filepath.Join(tempDir, fmt.Sprintf("page_%d.jpg", n))
		f, err := os.Create(imgPath)
		if err != nil {
			// return fmt.Errorf("failed to create image file: %w", err)
			return HandleError(err, "Failed to create image file", Error)
		}
		if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 90}); err != nil {
			f.Close()
			// return fmt.Errorf("failed to encode image: %w", err)
			return HandleError(err, "Failed to encode image", Error)
		}
		f.Close()

		err = AddIDToImage(imgPath, id, signature)
		if err != nil {
			// return fmt.Errorf("failed to annotate image page %d: %w", n+1, err)
			return HandleError(err, fmt.Sprintf("Failed to annotate image page %d", n+1), Error)
		}

		imagePaths = append(imagePaths, imgPath)
	}

	newPDF := gopdf.GoPdf{}
	newPDF.Start(gopdf.Config{})

	for _, imgPath := range imagePaths {
		imgFile, err := os.Open(imgPath)
		if err != nil {
			// return fmt.Errorf("failed to open image: %w", err)
			return HandleError(err, "Failed to open image", Error)
		}
		imgConf, _, err := image.DecodeConfig(imgFile)
		imgFile.Close()
		if err != nil {
			// return fmt.Errorf("failed to decode image size: %w", err)
			return HandleError(err, "Failed to decode image size", Error)
		}

		newPDF.AddPageWithOption(gopdf.PageOption{
			PageSize: &gopdf.Rect{
				W: float64(imgConf.Width),
				H: float64(imgConf.Height),
			},
		})
		err = newPDF.Image(imgPath, 0, 0, &gopdf.Rect{
			W: float64(imgConf.Width),
			H: float64(imgConf.Height),
		})
		if err != nil {
			// return fmt.Errorf("failed to add image to PDF: %w", err)
			return HandleError(err, "Failed to add image to PDF", Error)
		}
	}

	tempOutput := filePath + ".signed.pdf"
	err = newPDF.WritePdf(tempOutput)
	if err != nil {
		// return fmt.Errorf("failed to write final PDF: %w", err)
		return HandleError(err, "Failed to write final PDF", Error)
	}

	err = os.Rename(tempOutput, filePath)
	if err != nil {
		// return fmt.Errorf("failed to overwrite original PDF: %w", err)
		return HandleError(err, "Failed to overwrite original PDF", Error)
	}

	for _, imgPath := range imagePaths {
		_ = os.Remove(imgPath)
	}
	_ = os.RemoveAll(tempDir)

	comment := fmt.Sprintf("ID:%s;SIG:%s", id, signature)
	cmd := exec.Command("exiftool", "-overwrite_original", "-UserComment="+comment, filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// return fmt.Errorf("failed to write Exif: %v | Output: %s", err, string(output))
		return HandleError(err, fmt.Sprintf("Failed to write Exif: %s", string(output)), Error)
	}

	return nil
}
