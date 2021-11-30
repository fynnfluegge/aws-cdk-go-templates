package main

// import (
// 	"testing"

// 	"github.com/aws/aws-cdk-go/awscdk"
// 	"github.com/aws/aws-cdk-go/awscdk/assertions"
// 	"github.com/aws/jsii-runtime-go"
// )

// example tests. To run these tests, uncomment this file along with the
// example resource in eventbridge-sqs_test.go
// func TestEventbridgeSqsStack(t *testing.T) {
// 	// GIVEN
// 	app := awscdk.NewApp(nil)

// 	// WHEN
// 	stack := NewEventbridgeSqsStack(app, "MyStack", nil)

// 	// THEN
// 	template := assertions.Template_FromStack(stack)

// 	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
// 		"VisibilityTimeout": 300,
// 	})
// }
