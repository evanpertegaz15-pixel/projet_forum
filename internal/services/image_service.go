// upload, resize, suppression automatique d'une image

package services

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	//"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"golang.org/x/image/draw"
)

type ImageService struct {
	UploadDir string
	MaxWidth  int
	MaxHeight int
}

func NewImageService(uploadDir string) *ImageService {
	return &ImageService{
		UploadDir: uploadDir,
		MaxWidth:  800,
		MaxHeight: 800,
	}
}

// Upload + resize
func (s *ImageService) UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	// Vérifier type fichier
	ext := filepath.Ext(header.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", errors.New("format non supporté")
	}

	// Décoder image
	img, format, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Resize
	resized := s.resizeImage(img)

	// Nom unique
	filename := generateFileName(ext)
	path := filepath.Join(s.UploadDir, filename)

	// Créer fichier
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Encoder selon format
	switch format {
	case "jpeg":
		err = jpeg.Encode(out, resized, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(out, resized)
	default:
		return "", errors.New("format inconnu")
	}

	if err != nil {
		return "", err
	}

	return filename, nil
}

// Resize image
func (s *ImageService) resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Si déjà petite → rien à faire
	if width <= s.MaxWidth && height <= s.MaxHeight {
		return img
	}

	// Calcul ratio
	ratioW := float64(s.MaxWidth) / float64(width)
	ratioH := float64(s.MaxHeight) / float64(height)
	ratio := ratioW
	if ratioH < ratioW {
		ratio = ratioH
	}

	newW := int(float64(width) * ratio)
	newH := int(float64(height) * ratio)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst
}

// Supprimer une image
func (s *ImageService) DeleteImage(filename string) error {
	path := filepath.Join(s.UploadDir, filename)
	return os.Remove(path)
}

// 🧹 Suppression automatique (ex: images non utilisées)
func (s *ImageService) CleanupOldImages(maxAge time.Duration) error {
	files, err := os.ReadDir(s.UploadDir)
	if err != nil {
		return err
	}

	now := time.Now()

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			_ = os.Remove(filepath.Join(s.UploadDir, file.Name()))
		}
	}

	return nil
}

// 🏷 Générer nom unique
func generateFileName(ext string) string {
	return time.Now().Format("20060102150405") + "_" + randomString(6) + ext
}

// Petit générateur random
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
