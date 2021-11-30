package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/awssqs"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type EventbridgeSqsStackProps struct {
	awscdk.StackProps
}

func NewEventbridgeSqsStack(scope constructs.Construct, id string, props *EventbridgeSqsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// example resource
	queue := awssqs.NewQueue(stack, jsii.String("EventbridgeSqsQueue"), &awssqs.QueueProps{
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	})

	eventBus := awsevents.NewEventBus(stack, jsii.String("MyEventBus"), &awsevents.EventBusProps{
		EventBusName: jsii.String("MyEventBus"),
	})

	rule := awsevents.NewRule(stack, jsii.String("myEventBusRule"), &awsevents.RuleProps{
		Description: jsii.String("Log all events"),
		EventBus:    eventBus,
		EventPattern: &awsevents.EventPattern{
			Source:     &[]*string{jsii.String("MyCdkApp")},
			DetailType: &[]*string{jsii.String("message-for-queue")},
			Region:     &[]*string{jsii.String("eu-central-1")},
		},
	})

	rule.AddTarget(awseventstargets.NewSqsQueue(queue, &awseventstargets.SqsQueueProps{}))

	awscdk.NewCfnOutput(stack, jsii.String("queueUrl"), &awscdk.CfnOutputProps{
		Value:       queue.QueueUrl(),
		Description: jsii.String("URL of SQS Queue"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewEventbridgeSqsStack(app, "EventbridgeSqsStack", &EventbridgeSqsStackProps{
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
