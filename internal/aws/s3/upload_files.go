package s3

import (
	"context"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/schollz/progressbar/v3"
)

func DetectContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".js":
		return "application/javascript"
	case ".mjs":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".json":
		return "application/json"
	case ".svg":
		return "image/svg+xml"
	case ".wasm":
		return "application/wasm"
	}

	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	// Fallback to content sniffing
	file, err := os.Open(path)
	if err != nil {
		return "application/octet-stream"
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	return http.DetectContentType(buf[:n])
}

func UploadFolderContents(ctx context.Context, client *s3.Client, bucket, folder string) error {
	uploader := manager.NewUploader(client)

	folder = filepath.Clean(folder)

	// Count the total number of files to upload
	var totalFiles int
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Initialize the progress bar
	bar := progressbar.Default(int64(totalFiles), "Uploading files")

	// Walk through the folder and upload files
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}

		key := strings.ReplaceAll(relPath, "\\", "/")

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		ctype := DetectContentType(path)

		cacheControl := "max-age=31536000,public"
		if key == "index.html" {
			cacheControl = "no-cache,no-store,must-revalidate"
		}

		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket:       sdkaws.String(bucket),
			Key:          sdkaws.String(key),
			Body:         file,
			ContentType:  sdkaws.String(ctype),
			CacheControl: sdkaws.String(cacheControl),
		})
		if err != nil {
			return err
		}

		// Increment the progress bar
		bar.Add(1)
		return nil
	})
}
