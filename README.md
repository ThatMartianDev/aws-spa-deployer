# SPA Deployer

`SPA Deployer` is a CLI tool for deploying Single Page Applications (SPAs) to AWS S3 and CloudFront. It automates the process of uploading static assets to an S3 bucket, configuring the bucket for public access and static website hosting, and setting up a CloudFront distribution for global content delivery.

> **Note**: This project is currently in development. Future updates will include ACM (AWS Certificate Manager) integration for HTTPS, domain name configuration, CI/CD pipelines, and potentially a GitHub App or Dockerized version for easier accessibility.

## Why make it?
for a couple of reasons: I regularly deploy frontend apps at work, I was learning go, I always wanted to make a full cli app. So finally decided to spend some of my free time building this tool! This tool will eventually have everything in this guide I wrote + more:
https://www.notion.so/Deploying-SPA-React-Apps-on-AWS-292418753f0180989534f7bdc0d9fa32

---

## Features

- Automatically creates and configures an S3 bucket for hosting SPAs.
- Configures S3 bucket for public access and static website hosting.
- Uploads static files from a specified build directory to the S3 bucket.
- Creates a CloudFront distribution for the S3 bucket to enable global content delivery.
- Validates AWS region and S3 bucket name inputs.

### Planned Features

- **ACM Integration**: Automate the creation and management of SSL/TLS certificates for secure HTTPS connections.
- **Custom Domain Support**: Configure custom domains for your deployed SPAs using Route 53 or other DNS providers.
- **CI/CD Integration**:
  - GitHub Actions workflows for automated deployments.
  - AWS CodeBuild pipelines for seamless deployment.
- **GitHub App**: Explore the possibility of creating a GitHub App for easier integration with repositories.
- **Dockerized Version**: Provide a Docker image for running the tool in containerized environments.

---

## Project Structure

```
spa-deployer/
├── cmd/
│   └── spa-deploy/
│       └── main.go                # Entry point for the CLI application
├── internal/
│   ├── aws/
│   │   ├── session.go             # AWS session configuration
│   │   ├── cloudfront/
│   │   │   └── create_cfd.go      # CloudFront distribution creation logic
│   │   └── s3/
│   │       ├── configure_bucket.go # S3 bucket configuration (public access, static hosting)
│   │       ├── create_bucket.go    # S3 bucket creation and validation
│   │       └── upload_files.go     # Upload files to S3 with progress bar
│   ├── config/
│   │   └── config.go              # Configuration struct for the application
│   ├── data/
│   │   └── s3_policies.go         # S3 bucket policy generation
│   ├── deploy/
│   │   └── deploy.go              # Main deployment logic
│   └── helpers/
│       ├── cli_menu.go            # CLI menu for user input
│       └── validate_flags.go      # Validation for input flags (region, bucket name, etc.)
├── go.mod                         # Go module dependencies
└── README.md                      # Project documentation
```

---

## Prerequisites

- [Go](https://golang.org/) (version 1.20 or later)
- AWS account with appropriate permissions to create and manage S3 buckets and CloudFront distributions.
- AWS CLI configured with credentials and default region.

---


## Usage

The CLI tool requires the following flags:

- `--app`: The name of your application.
- `--bucket`: The name of the S3 bucket to use or create.
- `--region`: The AWS region where the S3 bucket and CloudFront distribution will be created. Defaults to `us-east-1`.
- `--dist`: The path to the build output directory of your SPA. Defaults to `./dist`.

Example usage:
```bash
./spa-deploy --app my-spa --bucket my-spa-bucket --region us-east-1 --dist ./dist
```

If any of the required flags are missing, the CLI will prompt you to input the missing values interactively.

---

## Deployment Steps

1. **Verify Build Directory**: The tool checks if the specified build directory exists.
2. **AWS Configuration**: Loads AWS configuration for the specified region.
3. **S3 Bucket Setup**:
   - Ensures the specified S3 bucket exists or creates a new one.
   - Configures the bucket for public access.
   - Applies a public read bucket policy.
   - Configures the bucket for static website hosting.
4. **Upload Files**: Uploads the contents of the build directory to the S3 bucket with appropriate content types and cache control headers.
5. **CloudFront Distribution**:
   - Creates a CloudFront distribution for the S3 bucket.
   - Outputs the CloudFront domain name for accessing the deployed SPA.

---

## AWS Permissions

To use this tool, ensure your AWS IAM user/role has the following permissions:

- `s3:CreateBucket`
- `s3:PutBucketPolicy`
- `s3:PutBucketWebsite`
- `s3:PutObject`
- `s3:PutPublicAccessBlock`
- `cloudfront:CreateDistribution`

---

## Future Plans

- **ACM Integration**: Automate SSL/TLS certificate creation and management for secure HTTPS connections.
- **Custom Domain Support**: Add support for custom domains using AWS Route 53 or other DNS providers.
- **CI/CD Pipelines**:
  - Automate deployments using GitHub Actions.
  - Provide AWS CodeBuild integration for CI/CD workflows.
- **GitHub App**: Explore creating a GitHub App for seamless integration with repositories.
- **Docker Support**: Provide a Docker image for running the tool in containerized environments.

---

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.
```