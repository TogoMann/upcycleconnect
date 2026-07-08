package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const MaxUploadSize = 5 * 1024 * 1024

func SaveImage(r *http.Request, fieldName string) (string, error) {

	r.Body = http.MaxBytesReader(nil, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		return "", fmt.Errorf("Fichier trop volumineux ou invalide (max 5 Mo)")
	}

	file, handler, err := r.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("Erreur lors de la récupération du fichier: %v", err)
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/webp" {
		return "", fmt.Errorf("Type de fichier non supporté (JPG, PNG, WEBP uniquement)")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), handler.Filename)
	uploadDir := "./uploads"
	filePath := filepath.Join(uploadDir, filename)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

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

const MaxDocumentUploadSize = 10 * 1024 * 1024

func SaveDocument(r *http.Request, fieldName string) (filename string, originalName string, err error) {
	r.Body = http.MaxBytesReader(nil, r.Body, MaxDocumentUploadSize)
	if err := r.ParseMultipartForm(MaxDocumentUploadSize); err != nil {
		return "", "", fmt.Errorf("Fichier trop volumineux ou invalide (max 10 Mo)")
	}

	file, handler, err := r.FormFile(fieldName)
	if err != nil {
		return "", "", fmt.Errorf("Erreur lors de la récupération du fichier: %v", err)
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", "", err
	}
	fileType := http.DetectContentType(buff)
	allowed := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
		"image/webp":      true,
	}
	if !allowed[fileType] {
		return "", "", fmt.Errorf("Type de fichier non supporté (PDF, JPG, PNG, WEBP uniquement)")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", err
	}

	savedName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), handler.Filename)
	uploadDir := "./uploads"
	filePath := filepath.Join(uploadDir, savedName)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", "", err
	}

	return savedName, handler.Filename, nil
}
