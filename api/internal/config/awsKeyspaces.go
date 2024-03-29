package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
)

func configureKeyspaces() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("cassandra.us-east-1.amazonaws.com")
	cluster.Port = 9142
	var auth sigv4.AwsAuthenticator = sigv4.NewAwsAuthenticator()

	auth.Region = os.Getenv("AWS_REGION")
	auth.AccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	auth.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	cluster.Authenticator = auth

	certFilePath := getCertificatePath()

	cluster.SslOpts = &gocql.SslOptions{
		CaPath:                 certFilePath,
		EnableHostVerification: false,
	}

	cluster.Consistency = gocql.LocalQuorum
	cluster.DisableInitialHostLookup = true

	return cluster
}

func getCertificatePath() string {
	ex, err := os.Executable()

	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)

	return fmt.Sprintf("%sAmazonRootCA1.pem", exPath+"/certs/")
}
