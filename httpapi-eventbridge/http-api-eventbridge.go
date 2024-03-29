package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HttpApiEventbridgeStackProps struct {
	awscdk.StackProps
}

func NewHttpApiEventbridgeStack(scope constructs.Construct, id string, props *HttpApiEventbridgeStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	eventBus := awsevents.NewEventBus(stack, jsii.String("myEventBus"), &awsevents.EventBusProps{
		EventBusName: jsii.String("MyEventBus"),
	})

	eventLoggerRule := awsevents.NewRule(stack, jsii.String("myEventLoggerRule"), &awsevents.RuleProps{
		Description: jsii.String("Log all events"),
		EventBus:    eventBus,
		EventPattern: &awsevents.EventPattern{
			Region: jsii.Strings(*props.Env.Region),
		},
	})

	logGroup := awslogs.NewLogGroup(stack, jsii.String("MyEventLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/events/MyEventBus"),
	})

	eventLoggerRule.AddTarget(awseventstargets.NewCloudWatchLogGroup(logGroup, &awseventstargets.LogGroupProps{}))

	httpApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("myHttpApi"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("myHttpApi"),
	})

	apiRole := awsiam.NewRole(stack, jsii.String("myEventBridgeIntegrationRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})

	apiRole.AddToPolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Resources: jsii.Strings(*eventBus.EventBusArn()),
			Actions:   jsii.Strings("events:putEvents"),
		}),
	)

	eventbridgeIntegration := awsapigatewayv2.NewCfnIntegration(stack, jsii.String("myEventbridgeIntegration"), &awsapigatewayv2.CfnIntegrationProps{
		ApiId:              httpApi.HttpApiId(),
		IntegrationType:    jsii.String("AWS_PROXY"),
		IntegrationSubtype: jsii.String("EventBridge-PutEvents"),
		CredentialsArn:     apiRole.RoleArn(),
		RequestParameters: &map[string]string{
			"Source":       "WebApp",
			"DetailType":   "MyDetailType",
			"Detail":       "$request.body",
			"EventBusName": *eventBus.EventBusArn(),
		},
		PayloadFormatVersion: jsii.String("1.0"),
		TimeoutInMillis:      jsii.Number(10000.0),
	})

	awsapigatewayv2.NewCfnRoute(stack, jsii.String("myEventRoute"), &awsapigatewayv2.CfnRouteProps{
		ApiId:    httpApi.HttpApiId(),
		RouteKey: jsii.String("POST /"),
		Target:   jsii.String("integrations/" + *eventbridgeIntegration.Ref()),
	})

	awscdk.NewCfnOutput(stack, jsii.String("apiUrl"), &awscdk.CfnOutputProps{
		Value:       httpApi.Url(),
		Description: jsii.String("HTTP API endpoint URL"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewHttpApiEventbridgeStack(app, "HttpApiEventbridgeStack", &HttpApiEventbridgeStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String(os.Getenv("AWS_REGION")),
	}
}
