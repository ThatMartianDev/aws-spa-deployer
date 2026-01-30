package data

import "fmt"

// BucketPolicy generates a public read policy for the specified bucket
func BucketPolicy(bucket string) string {
	return fmt.Sprintf(`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::%s/*"
        }
    ]
}`, bucket)
}
