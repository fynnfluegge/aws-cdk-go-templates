package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CognitoHttpapiStackProps struct {
	awscdk.StackProps
}

func NewCognitoHttpapiStack(scope constructs.Construct, id string, props *CognitoHttpapiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	userpool := awscognito.NewUserPool(stack, jsii.String("myCognitoUserPool"), &awscognito.UserPoolProps{
		UserPoolName: jsii.String("myCognitoUserPool"),
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength: jsii.Number(8),
		},
		AccountRecovery: awscognito.AccountRecovery_EMAIL_ONLY,
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(true),
		},
		StandardAttributes: &awscognito.StandardAttributes{
			Email: &awscognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(false),
			},
		},
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
		SelfSignUpEnabled: jsii.Bool(true),
	})

	awscognito.NewUserPoolClient(stack, jsii.String("MyUserPoolClient"), &awscognito.UserPoolClientProps{
		UserPoolClientName: jsii.String("MyUserPoolClient"),
		UserPool:           userpool,
		GenerateSecret:     jsii.Bool(false),
		SupportedIdentityProviders: &[]awscognito.UserPoolClientIdentityProvider{
			awscognito.UserPoolClientIdentityProvider_COGNITO(),
		},
		AuthFlows: &awscognito.AuthFlow{
			UserPassword:      jsii.Bool(true),
			AdminUserPassword: jsii.Bool(true),
		},
		OAuth: &awscognito.OAuthSettings{
			LogoutUrls:   jsii.Strings("https://oauth.pstmn.io/v1/callback"),
			CallbackUrls: jsii.Strings("https://oauth.pstmn.io/v1/callback"),
			Flows: &awscognito.OAuthFlows{
				AuthorizationCodeGrant: jsii.Bool(true),
				ImplicitCodeGrant:      jsii.Bool(true),
			},
			Scopes: &[]awscognito.OAuthScope{
				awscognito.OAuthScope_EMAIL(),
				awscognito.OAuthScope_OPENID(),
				awscognito.OAuthScope_PROFILE(),
			},
		},
	})

	userPoolDomain := awscognito.NewUserPoolDomain(stack, jsii.String("MyUserPoolDomain"), &awscognito.UserPoolDomainProps{
		UserPool: userpool,
		CognitoDomain: &awscognito.CognitoDomainOptions{
			DomainPrefix: jsii.String("myauth"),
		},
	})

	httpApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("MyHttpApi"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("MyHttpApi"),
		CorsPreflight: &awscdkapigatewayv2alpha.CorsPreflightOptions{
			AllowMethods: &[]awscdkapigatewayv2alpha.CorsHttpMethod{
				awscdkapigatewayv2alpha.CorsHttpMethod_GET,
			},
		},
	})

	awscdkapigatewayv2alpha.NewHttpAuthorizer(stack, jsii.String("MyHttpAuthorizer"), &awscdkapigatewayv2alpha.HttpAuthorizerProps{
		AuthorizerName: jsii.String("MyHttpAuthorizer"),
		Type:           awscdkapigatewayv2alpha.HttpAuthorizerType_JWT,
		HttpApi:        httpApi,
		JwtIssuer:      jsii.String("https://cognito-idp.eu-central-1.amazonaws.com/" + *userpool.UserPoolId()),
		JwtAudience:    jsii.Strings(*userpool.UserPoolId()),
		IdentitySource: jsii.Strings("$request.header.Authorization"),
	})

	awscdk.NewCfnOutput(stack, jsii.String("authUrl"), &awscdk.CfnOutputProps{
		Value: jsii.String("https://" + *userPoolDomain.DomainName() + ".auth.eu-central-1.amazoncognito.com/login"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCognitoHttpapiStack(app, "CognitoHttpapiStack", &CognitoHttpapiStackProps{
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
