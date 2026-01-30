package helpers

import (
	"fmt"
	"regexp"

	"github.com/ThatMartianDev/spa-deployer/internal/config"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func ValidateFlags(cfg *config.Config) (bool, retry bool) {
	// Validate Region
	if cfg.Region != "" {
		valid, suggestion := ValidateRegion(cfg.Region)
		if !valid {
			fmt.Printf("Invalid AWS region: %s\n", cfg.Region)
			if suggestion != "" {
				// Use the CLI menu to display options
				options := []string{
					fmt.Sprintf("Yes, use the suggestion: '%s'", suggestion),
					"No, use what I input",
					"No, use what I input and don't ask again",
					"Cancel",
				}
				choice := DisplayMenu(suggestion, options)

				switch choice {
				case options[0]:
					cfg.Region = suggestion
					return true, false
				case options[1]:
					return true, false
				case options[2]:
					// You can implement logic to save this preference if needed
					return true, false
				case options[3]:
					return false, false
				}
			} else {
				options := []string{
					fmt.Sprintf("Use what I input: '%s'", cfg.Region),
					"Use what I input and don't ask again",
					"Input a different region",
					"Cancel",
				}
				choice := DisplayMenu(suggestion, options)

				switch choice {
				case options[0]:
					return true, false
				case options[1]:
					// You can implement logic to save this preference if needed
					return true, false
				case options[2]:
					// Clear the region to prompt for re-entry on main
					cfg.Region = ""
					return false, true
				case options[3]:
					return false, false
				}
			}
			return false, false
		}
	}

	// Validate Bucket Name
	if cfg.Bucket != "" {
		valid, errorMessage := ValidateBucketName(cfg.Bucket)
		if !valid {
			fmt.Printf("Invalid S3 bucket name: %s\n", cfg.Bucket)
			options := []string{
				"Re-enter bucket name",
				"Cancel",
			}
			choice := DisplayMenu(errorMessage, options)

			switch choice {
			case options[0]:
				cfg.Bucket = ""
				return false, true
			case options[1]:
				return false, false
			}
		}
	}
	return true, false
}

// ValidateRegion validates the AWS region and suggests a fix if the input is close to a valid region
func ValidateRegion(region string) (bool, string) {
	regions := `(af|il|ap|ca|eu|me|sa|us|cn|us-gov|us-iso|us-isob)-(central|north|(north(?:east|west))|south|south(?:east|west)|east|west)-\d{1}`
	validRegions := []string{
		"af-south-1", "il-central-1", "ap-east-1", "ap-south-1", "ap-south-2", "ap-southeast-1", "ap-southeast-2",
		"ap-southeast-3", "ap-southeast-4", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ca-central-1",
		"eu-central-1", "eu-central-2", "eu-west-1", "eu-west-2", "eu-west-3", "eu-north-1", "eu-south-1",
		"eu-south-2", "me-central-1", "me-south-1", "sa-east-1", "us-east-1", "us-east-2", "us-west-1",
		"us-west-2", "cn-north-1", "cn-northwest-1", "us-gov-east-1", "us-gov-west-1", "us-iso-east-1",
		"us-isob-east-1",
	}

	// Check if the region matches the regex
	if regexp.MustCompile("^" + regions + "$").MatchString(region) {
		return true, ""
	}

	// If not valid, suggest the closest match
	closestRegion := ""
	minDistance := len(region)
	for _, validRegion := range validRegions {
		distance := levenshtein.DistanceForStrings([]rune(region), []rune(validRegion), levenshtein.DefaultOptions)
		if distance < minDistance {
			minDistance = distance
			closestRegion = validRegion
		}
	}

	// If the closest match is within a reasonable distance, suggest it
	if minDistance <= 3 {
		return false, closestRegion
	}

	return false, ""
}

func ValidateBucketName(bucket string) (bool, string) {
	bucketRegex := `^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`
	ipAddressRegex := `^\d{1,3}(\.\d{1,3}){3}$`

	// Check if the bucket name matches the S3 bucket naming rules
	if !regexp.MustCompile(bucketRegex).MatchString(bucket) {
		return false, "Bucket name must be 3-63 characters long, only contain lowercase letters, numbers, dots, and hyphens, and must start and end with a letter or number."
	}

	// Check if the bucket name is formatted as an IP address
	if regexp.MustCompile(ipAddressRegex).MatchString(bucket) {
		return false, "Bucket name cannot be formatted as an IP address (e.g., 192.168.1.1)."
	}

	return true, ""
}
