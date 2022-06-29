package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var (
	AWS aws
	K8S k8s
	WORKFLOW workflow
)

func Initialize(){

	var aws_default_region = os.Getenv("AWS_DEFAULT_REGION")
	if aws_default_region == "" {
		aws_default_region = "us-east-1"
		logrus.Infof("INFO environment variable AWS_DEFAULT_REGION is empty, setting default value: %s", aws_default_region)
	}

	var kubeconfig_path = os.Getenv("KUBECONFIG")
	if kubeconfig_path == "" {
		kubeconfig_path = "/etc/rancher/k3s/k3s.yaml"
		logrus.Warnf("WARNING environment variable KUBECONFIG is empty, setting default value: %s", kubeconfig_path)
	}


	var workflow_check_interval int

	var workflow_check_interval_str = os.Getenv("WORKFLOW_CHECK_INTERVAL")
	if workflow_check_interval_str == "" {
		workflow_check_interval = 5
		logrus.Warnf("WARNING environment variable WORKFLOW_CHECK_INTERVAL is empty, setting default value: %d", workflow_check_interval)
	}

	workflow_check_interval, err := strconv.Atoi(workflow_check_interval_str)
	if err != nil {
		workflow_check_interval = 5
		logrus.Errorf("ERROR parsing WORKFLOW_CHECK_INTERVAL value {%s} as int32: %v", workflow_check_interval_str, err)
		logrus.Warnf("WARNING unable to parse environmnet variable WORKFLOW_CHECK_INTERVAL, setting default value {%d}", workflow_check_interval)
	}

	AWS.Auth = awsAuth{
		AwsAccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AwsDefaultRegion:   aws_default_region,
	}

	AWS.DynamoDb = awsDynamoDb{
		Table: os.Getenv("DYNAMODB_TABLE"),
	}

	K8S = k8s{
		KubeConfigPath: kubeconfig_path,
	}

	WORKFLOW = workflow{
		CheckInterval: time.Duration(workflow_check_interval) * time.Second,
	}

}