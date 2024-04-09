package helpers

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func ImageToCloud(imageBase64 string, w http.ResponseWriter) (string, error) {
	imageParts := strings.Split(imageBase64, ",")
	if len(imageParts) != 2 {
		http.Error(w, "Invalid base64 image format", http.StatusBadRequest)
		return "", errors.New("invalid base64 image format")
	}

	// decode image from base64
	imageData, err := base64.StdEncoding.DecodeString(imageParts[1])
	if err != nil {
		http.Error(w, "Error decoding base64 image", http.StatusInternalServerError)
		return "", err
	}

	// upload image to Cloudinary
	cloudinaryURL, err := UploadToCloudinary(imageData)
	if err != nil {
		http.Error(w, "Error uploading image to Cloudinary", http.StatusInternalServerError)
		return "", err
	}

	return cloudinaryURL, nil
}

// Загружает изображение на Cloudinary и возвращает URL загруженного изображения
func UploadToCloudinary(data []byte) (string, error) {
	// Инициализация клиента Cloudinary
	cld, err := cloudinary.NewFromParams("djkotlye3", "888558296647534", "ruR8pPWSzFXyfD5dGv4GuWNDpYg")
	if err != nil {
		return "", err
	}

	// Загрузка изображения на Cloudinary
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, bytes.NewReader(data), uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	// Возврат URL загруженного изображения
	return uploadResult.SecureURL, nil
}

func DeleteFromCloudinary(url string) error {
	fmt.Println("entered deleteFromCloudinary")

	// Получаем публичный идентификатор из URL
	publicID := GetPublicIDFromURL(url)
	fmt.Println("public id: ", publicID)

	// Проверяем, не пустой ли publicID
	if publicID == "" {
		return errors.New("publicID is empty")
	}

	// Инициализация клиента Cloudinary
	cld, err := cloudinary.NewFromParams("djkotlye3", "888558296647534", "ruR8pPWSzFXyfD5dGv4GuWNDpYg")
	if err != nil {
		return err
	}

	// Подготовка параметров для удаления
	params := uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	}

	// Удаление ресурса из Cloudinary
	ctx := context.Background()
	result, err := cld.Upload.Destroy(ctx, params)
	if err != nil {
		return err
	}

	// Печатаем результат удаления
	fmt.Println("Delete result:", result)

	return nil
}

func GetPublicIDFromURL(url string) string {
	// Разбиваем URL по "/"
	parts := strings.Split(url, "/")

	// По умолчанию публичный идентификатор - последний сегмент URL без расширения
	lastSegment := parts[len(parts)-1]
	publicID := strings.TrimSuffix(lastSegment, filepath.Ext(lastSegment))

	return publicID
}
