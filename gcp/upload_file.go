package gcp

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

// UploadFile => uploads an object
// https://adityarama1210.medium.com/simple-golang-api-uploader-using-google-cloud-storage-3d5e45df74a5
func (c *gcpClient) UploadFile(file *multipart.FileHeader, path, fileName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	blobFile, _ := file.Open()

	// Upload an object with storage.Writer.
	wc := client.Bucket(c.bucketName).Object(path + fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, blobFile); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
