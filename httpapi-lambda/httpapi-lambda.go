package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HttpapiLambdaStackProps struct {
	awscdk.StackProps
}

func NewHttpapiLambdaStack(scope constructs.Construct, id string, props *HttpapiLambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	getHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("myGoHandler"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./lambda-handler"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w"`)},
		},
	})

	httpApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("myHttpApi"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("myHttpApi"),
	})

	httpApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), getHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	awscdk.NewCfnOutput(stack, jsii.String("myHttpApiEndpoint"), &awscdk.CfnOutputProps{
		Value:       httpApi.ApiEndpoint(),
		Description: jsii.String("HTTP API Endpoint"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewHttpapiLambdaStack(app, "HttpapiLambdaStack", &HttpapiLambdaStackProps{
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
