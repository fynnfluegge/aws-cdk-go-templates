package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type AppServerlessCdkGoStackProps struct {
	awscdk.StackProps
}

func NewAppServerlessCdkGoStack(scope constructs.Construct, id string, props *AppServerlessCdkGoStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	getHandler := awslambdago.NewGoFunction(stack, jsii.String("myGoHandler"), &awslambdago.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./basic-api-app/handler-function-get"),
		Bundling: &awslambdago.BundlingOptions{
			GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w"`)},
		},
	})

	restApi := awsapigateway.NewRestApi(stack, jsii.String("myGoApi"), &awsapigateway.RestApiProps{
		RestApiName:    jsii.String("myGoApi"),
		CloudWatchRole: jsii.Bool(false),
	})

	restApi.Root().AddResource(jsii.String("hello-world"), &awsapigateway.ResourceOptions{}).AddMethod(
		jsii.String("GET"),
		awsapigateway.NewLambdaIntegration(getHandler, &awsapigateway.LambdaIntegrationOptions{}),
		restApi.Root().DefaultMethodOptions(),
	)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewAppServerlessCdkGoStack(app, "AppServerlessCdkGoStack", &AppServerlessCdkGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
