package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ThatMartianDev/spa-deployer/internal/config"
	"github.com/ThatMartianDev/spa-deployer/internal/deploy"
	"github.com/ThatMartianDev/spa-deployer/internal/helpers"
	"github.com/sirupsen/logrus"
)

const (
	defaultRegion = "us-east-1"
	defaultDist   = "./dist"
)

func SetupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	return logger
}

func main() {
	logger := SetupLogger()

	cfg := config.Config{}

	flag.StringVar(&cfg.AppName, "app", "", "App name")
	flag.StringVar(&cfg.Bucket, "bucket", "", "S3 bucket name")
	flag.StringVar(&cfg.Region, "region", defaultRegion, "AWS region")
	flag.StringVar(&cfg.DistDir, "dist", defaultDist, "Build output directory")

	flag.Parse()

	ctx := context.Background()

	if cfg.AppName == "" || cfg.Bucket == "" || cfg.Region == "" || cfg.DistDir == "" {
		flagsInputs(&cfg)
	}

	validateRetry(&cfg, logger)

	if err := deploy.Run(logger, ctx, cfg); err != nil {
		logger.WithError(err).Error("Deployment failed")
		os.Exit(1)
	}

	logger.Info("✔ Deployment finished successfully")
}

func validateRetry(cfg *config.Config, logger *logrus.Logger) {
	for {
		valid, retry := helpers.ValidateFlags(cfg)
		if valid {
			break
		}
		if retry {
			flagsInputs(cfg)
		} else {
			logger.Error("Invalid flags provided. Exiting.")
			os.Exit(1)
		}
	}
}

func flagsInputs(cfg *config.Config) {
	reader := bufio.NewReader(os.Stdin)

	if cfg.AppName == "" {
		fmt.Print("Enter App Name: ")
		appName, _ := reader.ReadString('\n')
		cfg.AppName = strings.TrimSpace(appName)
	}

	if cfg.Bucket == "" {
		fmt.Print("Enter S3 Bucket Name: ")
		bucket, _ := reader.ReadString('\n')
		cfg.Bucket = strings.TrimSpace(bucket)
	}

	if cfg.Region == "" {
		fmt.Print("Enter AWS Region (default us-east-1): ")
		region, _ := reader.ReadString('\n')
		cfg.Region = strings.TrimSpace(region)
	}

	if cfg.DistDir == "" {
		fmt.Print("Enter Build Output Directory (default ./dist): ")
		distDir, _ := reader.ReadString('\n')
		cfg.DistDir = strings.TrimSpace(distDir)
	}
}
