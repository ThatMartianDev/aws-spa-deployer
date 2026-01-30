package s3

import (
	"context"
	"log"

	"github.com/ThatMartianDev/spa-deployer/internal/data"
	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	sdks3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// - allow public access to bucket (turning off "Block all public access")
func AllowPublicAccess(ctx context.Context, client *sdks3.Client, bucket string) error {
	_, err := client.PutPublicAccessBlock(ctx, &sdks3.PutPublicAccessBlockInput{
		Bucket: sdkaws.String(bucket),
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       sdkaws.Bool(false),
			BlockPublicPolicy:     sdkaws.Bool(false),
			IgnorePublicAcls:      sdkaws.Bool(false),
			RestrictPublicBuckets: sdkaws.Bool(false),
		},
	})
	if err != nil {
		return err
	}
	log.Println("Public access settings updated for bucket:", bucket)
	return nil
}

// - apply bucket policy to allow public read access
func ApplyBucketPolicy(ctx context.Context, client *sdks3.Client, bucket string) error {
	policy := data.BucketPolicy(bucket)

	_, err := client.PutBucketPolicy(ctx, &sdks3.PutBucketPolicyInput{
		Bucket: sdkaws.String(bucket),
		Policy: sdkaws.String(policy),
	})
	return err
}

// - configure static website hosting
func ConfigureStaticWebsite(ctx context.Context, client *sdks3.Client, bucket string) error {
	_, err := client.PutBucketWebsite(ctx, &sdks3.PutBucketWebsiteInput{
		Bucket: &bucket,
		WebsiteConfiguration: &types.WebsiteConfiguration{
			IndexDocument: &types.IndexDocument{
				Suffix: sdkaws.String("index.html"),
			},
			ErrorDocument: &types.ErrorDocument{
				Key: sdkaws.String("index.html"),
			},
		},
	})

	return err
}
