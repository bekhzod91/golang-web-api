package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/myproject/api/config"
	"io"
)

type Client struct {
	config    config.Config
	awsConfig aws.Config
}

func NewClient(cfg config.Config) (*Client, error) {
	awsCtx := context.Background()
	credentialsProvider := credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")
	awsCfg, err := awsConfig.LoadDefaultConfig(
		awsCtx,
		awsConfig.WithBaseEndpoint(cfg.AWSBaseEndpoint),
		awsConfig.WithRegion(cfg.AWSRegion),
		awsConfig.WithCredentialsProvider(credentialsProvider),
	)
	if err != nil {
		return nil, err
	}

	client := Client{config: cfg, awsConfig: awsCfg}

	return &client, nil
}

func (c *Client) UploadS3(ctx context.Context, key string, body io.Reader) (string, error) {
	client := s3.NewFromConfig(c.awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.config.AWSBucket),
		Key:    aws.String(key),
		Body:   body,
	})

	if err != nil {
		return "", fmt.Errorf("unable to upload %s to %s, %w", key, c.config.AWSBucket, err)
	}

	return fmt.Sprintf("%s/%s", c.config.AWSBucket, key), err
}
