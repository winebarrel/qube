package rds

import (
	"context"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

func BuildIAMAuthToken(ctx context.Context, endpoint string, user string) (string, error) {
	awscfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return "", err
	}

	return auth.BuildAuthToken(ctx, endpoint, awscfg.Region, user, awscfg.Credentials)
}

func ResolveCNAME(host string) (string, error) {
	if strings.HasSuffix(host, ".rds.amazonaws.com") {
		return host, nil
	}

	host, err := net.LookupCNAME(host)

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(host, "."), nil
}
