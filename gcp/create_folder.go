package gcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
)

func (c *gcpClient) CreateFolder(bucket, path, folderPrefix string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(fmt.Sprintf("%s/%s", path, folderPrefix)).NewWriter(ctx)
	if _, err := io.Copy(wc, bytes.NewReader([]byte{})); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
