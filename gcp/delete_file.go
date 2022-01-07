package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
)

// DeleteFiles => delete
// https://cloud.google.com/storage/docs/deleting-objects#code-samples
func (c *gcpClient) DeleteFiles(bucket, path, fileName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(path + fileName)
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", path+fileName, err)
	}
	// fmt.Printf("Blob %v deleted.\n", object)
	return nil
}
