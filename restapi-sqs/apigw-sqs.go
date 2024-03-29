package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApigwSqsStackProps struct {
	awscdk.StackProps
}

func NewApigwSqsStack(scope constructs.Construct, id string, props *ApigwSqsStackProps) awscdk.Stack {

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	queue := awssqs.NewQueue(stack, jsii.String("EventbridgeSqsQueue"), &awssqs.QueueProps{
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
		QueueName:         jsii.String("MySqsQueue"),
	})

	restApiRole := awsiam.NewRole(stack, jsii.String("myRestAPIRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSQSFullAccess")),
		},
	})

	integrationResponse := &awsapigateway.IntegrationResponse{
		StatusCode:        jsii.String("200"),
		ResponseTemplates: &map[string]*string{"application/json": jsii.String("")},
	}

	integration := awsapigateway.NewAwsIntegration(&awsapigateway.AwsIntegrationProps{
		Service:               jsii.String("sqs"),
		IntegrationHttpMethod: jsii.String("POST"),
		Path:                  jsii.String(*props.Env.Account + *queue.QueueName()),
		Options: &awsapigateway.IntegrationOptions{
			CredentialsRole: restApiRole,
			IntegrationResponses: &[]*awsapigateway.IntegrationResponse{
				integrationResponse,
			},
			RequestTemplates:    &map[string]*string{"application/json": jsii.String("Action=SendMessage&MessageBody=$input.body")},
			PassthroughBehavior: awsapigateway.PassthroughBehavior_NEVER,
			RequestParameters:   &map[string]*string{"integration.request.header.Content-Type": jsii.String("'application/x-www-form-urlencoded'")},
		},
	})

	restApi := awsapigateway.NewRestApi(stack, jsii.String("myRestApi"), &awsapigateway.RestApiProps{
		DefaultIntegration: integration,
	})
	apiResource := restApi.Root().AddResource(jsii.String("test"), &awsapigateway.ResourceOptions{})

	methodResponse := awsapigateway.MethodResponse{StatusCode: jsii.String("200")}

	apiResource.AddMethod(jsii.String("POST"), integration, &awsapigateway.MethodOptions{
		MethodResponses: &[]*awsapigateway.MethodResponse{&methodResponse},
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewApigwSqsStack(app, "ApigwSqsStack", &ApigwSqsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT")),
	}
}
