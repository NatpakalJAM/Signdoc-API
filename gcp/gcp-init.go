package gcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"signdoc_api/config"

	"cloud.google.com/go/storage"
)

type gcpClient struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	bucket     *storage.BucketHandle
	uploadPath string
	w          io.Writer
	ctx        context.Context
	// cleanUp is a list of filenames that need cleaning up at the end of the demo.
	cleanUp []string
	// failed indicates that one or more of the demo steps failed.
	failed bool
}

var Client *gcpClient

// var client *storage.Client

func Init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.C.GCP.GoogleApplicationCredentials)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	buf := &bytes.Buffer{}

	Client = &gcpClient{
		cl:         client,
		projectID:  config.C.GCP.ProjectID,
		bucketName: config.C.GCP.BucketName,
		uploadPath: fmt.Sprintf("%s/", config.C.GCP.UploadPath),

		w: buf,
	}

	fmt.Println("gcp init completed.")
}

func (c *gcpClient) dumpStats(obj *storage.ObjectAttrs) {
	fmt.Fprintf(c.w, "(filename: /%v/%v, ", obj.Bucket, obj.Name)
	fmt.Fprintf(c.w, "ContentType: %q, ", obj.ContentType)
	fmt.Fprintf(c.w, "ACL: %#v, ", obj.ACL)
	fmt.Fprintf(c.w, "Owner: %v, ", obj.Owner)
	fmt.Fprintf(c.w, "ContentEncoding: %q, ", obj.ContentEncoding)
	fmt.Fprintf(c.w, "Size: %v, ", obj.Size)
	fmt.Fprintf(c.w, "MD5: %q, ", obj.MD5)
	fmt.Fprintf(c.w, "CRC32C: %q, ", obj.CRC32C)
	fmt.Fprintf(c.w, "Metadata: %#v, ", obj.Metadata)
	fmt.Fprintf(c.w, "MediaLink: %q, ", obj.MediaLink)
	fmt.Fprintf(c.w, "StorageClass: %q, ", obj.StorageClass)
	if !obj.Deleted.IsZero() {
		fmt.Fprintf(c.w, "Deleted: %v, ", obj.Deleted)
	}
	fmt.Fprintf(c.w, "Updated: %v)\n", obj.Updated)
}
