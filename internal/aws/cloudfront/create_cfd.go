package cloudfront

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func CreateCloudFrontDistribution(
	ctx context.Context,
	client *cf.Client,
	bucket string,
	region string,
	appName string,
) (string, error) {

	originID := "s3-website-" + bucket
	originDomain := fmt.Sprintf(
		"%s.s3-website-%s.amazonaws.com",
		bucket,
		region,
	)

	distConfig := &types.DistributionConfig{
		CallerReference: aws.String(
			fmt.Sprintf("%s-%d", appName, time.Now().Unix()),
		),
		Comment:           aws.String("SPA CloudFront distribution"),
		Enabled:           aws.Bool(true),
		DefaultRootObject: aws.String("index.html"),

		Origins: &types.Origins{
			Quantity: aws.Int32(1),
			Items: []types.Origin{
				{
					Id:         aws.String(originID),
					DomainName: aws.String(originDomain),
					CustomOriginConfig: &types.CustomOriginConfig{
						HTTPPort:             aws.Int32(80),
						HTTPSPort:            aws.Int32(443),
						OriginProtocolPolicy: types.OriginProtocolPolicyHttpOnly,
					},
				},
			},
		},

		DefaultCacheBehavior: &types.DefaultCacheBehavior{
			TargetOriginId:       aws.String(originID),
			ViewerProtocolPolicy: types.ViewerProtocolPolicyRedirectToHttps,
			Compress:             aws.Bool(true),

			AllowedMethods: &types.AllowedMethods{
				Quantity: aws.Int32(2),
				Items:    []types.Method{types.MethodGet, types.MethodHead},
			},

			ForwardedValues: &types.ForwardedValues{
				QueryString: aws.Bool(false),
				Cookies: &types.CookiePreference{
					Forward: types.ItemSelectionNone,
				},
			},

			MinTTL:     aws.Int64(0),
			DefaultTTL: aws.Int64(86400),
			MaxTTL:     aws.Int64(31536000),
		},

		PriceClass: types.PriceClassPriceClass100,

		ViewerCertificate: &types.ViewerCertificate{
			CloudFrontDefaultCertificate: aws.Bool(true),
		},

		HttpVersion:   types.HttpVersionHttp2,
		IsIPV6Enabled: aws.Bool(true),
	}

	out, err := client.CreateDistribution(ctx, &cf.CreateDistributionInput{
		DistributionConfig: distConfig,
	})
	if err != nil {
		return "", err
	}

	return *out.Distribution.DomainName, nil
}
