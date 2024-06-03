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

func UploadToCloudinary(data []byte) (string, error) {
	cld, err := cloudinary.NewFromParams("djkotlye3", "888558296647534", "ruR8pPWSzFXyfD5dGv4GuWNDpYg")
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, bytes.NewReader(data), uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func DeleteFromCloudinary(url string) error {

	publicID := GetPublicIDFromURL(url)

	if publicID == "" {
		return errors.New("publicID is empty")
	}

	cld, err := cloudinary.NewFromParams("djkotlye3", "888558296647534", "ruR8pPWSzFXyfD5dGv4GuWNDpYg")
	if err != nil {
		return err
	}

	params := uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	}

	ctx := context.Background()
	result, err := cld.Upload.Destroy(ctx, params)
	if err != nil {
		return err
	}

	fmt.Println("Delete result:", result)

	return nil
}

func GetPublicIDFromURL(url string) string {
	parts := strings.Split(url, "/")

	lastSegment := parts[len(parts)-1]
	publicID := strings.TrimSuffix(lastSegment, filepath.Ext(lastSegment))

	return publicID
}
