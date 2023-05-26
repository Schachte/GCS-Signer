package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

func main() {
	bucketName := flag.String("bucket", "<BUCKET>", "Name of the bucket")
	objectName := flag.String("object", "<OBJECT>", "Name of the object")
	serviceAccountKeyPath := flag.String("key", "access.json", "Path to the service account key file")

	flag.Parse()

	jsonKey, err := ioutil.ReadFile(*serviceAccountKeyPath)
	if err != nil {
		fmt.Printf("Failed to load service account key: %v\n", err)
		return
	}

	config, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		fmt.Printf("Failed to parse service account key: %v\n", err)
		return
	}

	opts := &storage.SignedURLOptions{
		Method:         "GET",
		GoogleAccessID: config.Email,
		PrivateKey:     config.PrivateKey,
		Expires:        time.Now().Add(24 * time.Hour), // URL valid for 24 hours
		Scheme:         storage.SigningSchemeV4,
	}

	url, err := storage.SignedURL(*bucketName, *objectName, &storage.SignedURLOptions{
		GoogleAccessID: opts.GoogleAccessID,
		PrivateKey:     opts.PrivateKey,
		Method:         "GET",
		Expires:        opts.Expires,
	})
	if err != nil {
		log.Fatalf("Failed to generate signed URL: %v", err)
	}
	fmt.Printf("%s\n", url)
}
