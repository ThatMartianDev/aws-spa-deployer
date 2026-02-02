package s3

import (
	"bufio"
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
)

func EnsureBucket(logger *logrus.Logger, ctx context.Context, client *s3.Client, bucket string, region string) error {
	for {
		_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: &bucket,
		})

		// Bucket exists and is accessible
		if err == nil {
			newBucket := PromptBucketName(logger, bucket)
			if newBucket == bucket {
				return nil
			}
			bucket = newBucket
			continue
		}

		var bucketExistsErr *types.BucketAlreadyOwnedByYou
		if errors.As(err, &bucketExistsErr) {
			logger.Infof("Bucket %s already exists and is owned by you", bucket)
			newBucket := PromptBucketName(logger, bucket)
			if newBucket == bucket {
				return nil
			}
			bucket = newBucket
			continue
		}

		var bucketExistsOtherAccountErr *types.BucketAlreadyExists
		if errors.As(err, &bucketExistsOtherAccountErr) {
			logger.Warnf("Bucket %s already exists and is owned by another account.", bucket)
			bucket = PromptBucketName(logger, bucket)
			continue
		}

		// Bucket does not exist or not accessible → try to create
		region := region
		createBucketInput := &s3.CreateBucketInput{
			Bucket: &bucket,
		}
		if region != "us-east-1" {
			createBucketInput.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(region),
			}
		}

		_, err = client.CreateBucket(ctx, createBucketInput)

		if err != nil {
			logger.Errorf("Failed to create bucket %s: %v", bucket, err)
			bucket = PromptBucketName(logger, "")
			continue
		}

		logger.Infof("Bucket %s created successfully", bucket)
		return nil
	}
}

func PromptBucketName(logger *logrus.Logger, existingBucket string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		if existingBucket != "" {
			logger.Infof("Bucket '%s' already exists. Do you want to use it? (yes/no): ", existingBucket)
			input, err := reader.ReadString('\n')
			if err != nil {
				logger.Error("Error reading input:", err)
				continue
			}
			input = strings.ToLower(strings.TrimSpace(input))
			switch input {
			case "yes", "y":
				return existingBucket
			case "no", "n":
				logger.Info("Please provide a new bucket name.")
				existingBucket = "" // Clear the existing bucket name
			}
		}

		logger.Info("Enter a new bucket name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			logger.Warn("Bucket name cannot be empty.")
			continue
		}
		return input
	}
}
