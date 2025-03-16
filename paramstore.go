package main

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/ptr"
)

type SSMClient struct {
	*ssm.Client
	keyID    string
	basePath string
}

func newSSMClient(ctx context.Context, keyID string, basePath string) (*SSMClient, error) {
	options := [](func(*config.LoadOptions) error){}
	awscfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	return &SSMClient{Client: ssm.NewFromConfig(awscfg), keyID: keyID, basePath: basePath}, err
}

func (client *SSMClient) fullPath(name string) string {
	return path.Join(client.basePath, name)
}

func (client *SSMClient) get(ctx context.Context, name string) (*ssm.GetParameterOutput, error) {
	return client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           ptr.String(client.fullPath(name)),
		WithDecryption: ptr.Bool(client.keyID != ""),
	})
}

func (client *SSMClient) put(ctx context.Context, name string, value string) (*ssm.PutParameterOutput, error) {
	params := ssm.PutParameterInput{
		Name:      ptr.String(client.fullPath(name)),
		Value:     ptr.String(value),
		DataType:  ptr.String("text"),
		Overwrite: ptr.Bool(true),
		Type:      types.ParameterTypeString,
	}

	if client.keyID != "" {
		params.KeyId = ptr.String(client.keyID)
		params.Type = types.ParameterTypeSecureString
	}

	return client.PutParameter(ctx, &params)
}
