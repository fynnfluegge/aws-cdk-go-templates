package main

// import (
// 	"testing"

// 	"github.com/aws/aws-cdk-go/awscdk/v2"
// 	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
// 	"github.com/aws/jsii-runtime-go"
// )

// example tests. To run these tests, uncomment this file along with the
// example resource in httpapi-lambda_test.go
// func TestHttpapiLambdaStack(t *testing.T) {
// 	// GIVEN
// 	app := awscdk.NewApp(nil)

// 	// WHEN
// 	stack := NewHttpapiLambdaStack(app, "MyStack", nil)

// 	// THEN
// 	template := assertions.Template_FromStack(stack)

// 	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
// 		"VisibilityTimeout": 300,
// 	})
// }
