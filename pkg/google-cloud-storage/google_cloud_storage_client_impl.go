package google_cloud_storage

import (
	"bulk/config"
	"bulk/models/response"
	"bulk/logger"
	"bulk/tracer"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"net/http"
	"time"
)

func (gcs *GoogleCloudStorage) SignedUrl(ctx context.Context, filePath string) (string, error) {
	spanName := "GoogleCloudStorage.SignedUrl"
	childCtx, span := tracer.StartSpan(ctx, spanName)

	_, cancel := context.WithTimeout(childCtx, 10*time.Second)
	defer cancel()

	object := response.Object{}
	bucketName := gcs.cfg.GCSBulk.Bucket

	gcs.log.Debugf("Bucket Name: [%v]", bucketName)

	// Get object through looping by having bucketName
	data := gcs.client.Bucket(bucketName).Objects(childCtx, &storage.Query{
		Prefix:    filePath,
		Delimiter: "/",
	})
	for {
		attrs, errData := data.Next()
		if errData != nil {
			if errData == iterator.Done {
				span.RecordError(errData)
				break
			}
			span.RecordError(errData)
			return "", fmt.Errorf(spanName+" Bucket(%v).Objects(): %v", bucketName, errData)
		}
		// Storing object Name to struct
		object.Name = attrs.Name
	}

	timeDuration := time.Duration(config.Cfg.GCSBulk.SignedUrlExpiredTimeInMinutes)

	// Configuration for SignedUrl so that we have 15 min for time expiration.
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         http.MethodGet,
		GoogleAccessID: config.Cfg.GCSBulk.GoogleAccessId,
		Expires:        time.Now().Add(timeDuration * time.Minute),
	}

	gcs.log.Debugf("Object Selected from Storage: [%v]", object.Name)

	url, err := gcs.client.Bucket(bucketName).SignedURL(object.Name, opts)
	if err != nil {
		span.RecordError(err)
		logger.Log.Errorf("Error while generate SignedUrl %v", err)
		return "", err
	}

	return url, nil
}
