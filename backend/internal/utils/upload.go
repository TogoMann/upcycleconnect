package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const MaxUploadSize = 5 * 1024 * 1024 // 5 MB

func SaveImage(r *http.Request, fieldName string) (string, error) {
	// Limit request size
	r.Body = http.MaxBytesReader(nil, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		return "", fmt.Errorf("Fichier trop volumineux ou invalide (max 5 Mo)")
	}

	file, handler, err := r.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("Erreur lors de la récupération du fichier: %v", err)
	}
	defer file.Close()

	// Check file type (optional but good)
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/webp" {
		return "", fmt.Errorf("Type de fichier non supporté (JPG, PNG, WEBP uniquement)")
	}

	// Reset file pointer after reading buff
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	// Create unique filename
	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), handler.Filename)
	uploadDir := "./uploads"
	filePath := filepath.Join(uploadDir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	// Save to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filename, nil
}
