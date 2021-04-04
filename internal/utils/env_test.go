package utils

import (
	"os"
	"testing"
)

func TestGetEnvWithNotNone(t *testing.T) {
	os.Setenv("CLOUD_TYPE", "azure")
	os.Setenv("AWS_REGION", "ap-northeast-1")
	os.Setenv("AWS_SECRET_NAME", "prod/profile")
	cloudType := GetEnv("CLOUD_TYPE", "aws")
	awsRegion := GetEnv("AWS_REGION", "ap-southeast-1")
	awsSecretName := GetEnv("AWS_SECRET_NAME", "")
	if cloudType != "azure" {
		t.Fatalf("Could not get environment variable value as expected [%s], actual [%s]", "azure", cloudType)
	}
	if awsRegion != "ap-northeast-1" {
		t.Fatalf("Could not get environment variable value as expected [%s], actual [%s]", "ap-northeast-1", awsRegion)
	}
	if awsSecretName != "prod/profile" {
		t.Fatalf("Could not get environment variable value as expected [%s], actual [%s]", "prod/profile", awsSecretName)
	}
}

func TestGetEnvWithNone(t *testing.T) {
	os.Setenv("CLOUD_TYPE", "")
	os.Setenv("AWS_REGION", "")
	os.Setenv("AWS_SECRET_NAME", "")
	cloudType := GetEnv("CLOUD_TYPE", "aws")
	awsRegion := GetEnv("AWS_REGION", "ap-southeast-1")
	awsSecretName := GetEnv("AWS_SECRET_NAME", "")
	if cloudType != "aws" {
		t.Fatalf("Could not get environment variable value as expected [%s], actual [%s]", "aws", cloudType)
	}
	if awsRegion != "ap-southeast-1" {
		t.Fatalf("Could not get environment variable value as expected [%s], actual [%s]", "ap-southeast-1", awsRegion)
	}
	if awsSecretName != "" {
		t.Fatalf("Could not get environment variable value as expected [%s], actual [%s]", "", awsSecretName)
	}
}
