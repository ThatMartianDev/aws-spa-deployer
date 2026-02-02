package deploy

import (
	"context"
	"fmt"
	"os"

	sdkCloudfront "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	sdkS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"

	"github.com/ThatMartianDev/spa-deployer/internal/aws"
	"github.com/ThatMartianDev/spa-deployer/internal/aws/cloudfront"
	"github.com/ThatMartianDev/spa-deployer/internal/aws/s3"
	"github.com/ThatMartianDev/spa-deployer/internal/config"
)

func Run(logger *logrus.Logger, ctx context.Context, cfg config.Config) error {

	// - verify dist folder exists
	if _, err := os.Stat(cfg.DistDir); os.IsNotExist(err) {
		return fmt.Errorf("build folder not found: %s", cfg.DistDir)
	}

	logger.Infof("Build folder found: %s", cfg.DistDir)

	// - setup AWS session
	awsCfg, err := aws.LoadAWSConfig(ctx, cfg.Region)
	if err != nil {
		return err
	}

	s3Client := sdkS3.NewFromConfig(awsCfg)

	for {
		if err := s3.EnsureBucket(logger, ctx, s3Client, cfg.Bucket, cfg.Region); err != nil {
			logger.Warnf("Bucket creation failed: %v", err)
			cfg.Bucket = s3.PromptBucketName(logger, cfg.Bucket)
			continue
		}
		break
	}

	// - allow public access
	logger.Info("Configuring bucket for public access...")
	if err := s3.AllowPublicAccess(ctx, s3Client, cfg.Bucket); err != nil {
		return err
	}

	// - apply bucket policy
	logger.Info("Applying public read bucket policy...")
	if err := s3.ApplyBucketPolicy(ctx, s3Client, cfg.Bucket); err != nil {
		return fmt.Errorf("failed to apply bucket policy: %w", err)
	}

	// - configure NEW bucket
	logger.Info("Configuring static website hosting")
	if err := s3.ConfigureStaticWebsite(ctx, s3Client, cfg.Bucket); err != nil {
		return err
	}
	logger.Info("✔ Bucket Ready")

	// - upload dist
	logger.Infof("Uploading build folder: %s", cfg.DistDir)
	if err := s3.UploadFolderContents(ctx, s3Client, cfg.Bucket, cfg.DistDir); err != nil {
		return err
	}
	logger.Info("✔ Files uploaded to S3 bucket")

	// - cloudfront setup
	cfClient := sdkCloudfront.NewFromConfig(awsCfg)

	logger.Info("Creating CloudFront distribution...")
	domain, err := cloudfront.CreateCloudFrontDistribution(
		ctx,
		cfClient,
		cfg.Bucket,
		cfg.Region,
		cfg.AppName,
	)
	if err != nil {
		return fmt.Errorf("cloudfront creation failed: %w", err)
	}

	logger.Info("✔ CloudFront distribution created")
	logger.Info("🌍 https://" + domain)

	return nil
}
