package iauth

import (
	"common/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
)

var _authSession *session.Session

func Authenticate() *session.Session{

	var sess *session.Session
	var err error

	if config.AWS.Auth.AwsAccessKeyId != "" && config.AWS.Auth.AwsSecretAccessKey != "" {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(config.AWS.Auth.AwsDefaultRegion),
			Credentials: credentials.NewStaticCredentials(config.AWS.Auth.AwsAccessKeyId, config.AWS.Auth.AwsSecretAccessKey, ""),
			MaxRetries:  aws.Int(5),
		})

	} else {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(config.AWS.Auth.AwsDefaultRegion),
			MaxRetries: aws.Int(5),
		})
	}

	if err != nil {
		logrus.Fatalf("FATAL unable to authenticate in AWS: %s", err)
	} else {
		_authSession = sess
	}

	return _authSession
}

func GetCurrentAuthSession() *session.Session{
	return _authSession
}
