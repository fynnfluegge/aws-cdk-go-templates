# AWS Lambda to Amazon DynamoDB

This example deploys a Lambda function that makes puts to a DynamoDB table.

## How it works

This is similar to the `lambda-dynamodb` example but is implemented in CDK.

- A Lambda function is invoked by an event integration or CLI command
- The Lambda function "stringifies" the event payload
- The Function uses the AWS SDK to perform a `put` command on a DynamoDB table
- The name of the DynamoDB table is passed to the Lambda function via an environment variable named `DatabaseTable`
- The Lambda function is granted `PutItem` permissions, defined in CDK via `dynamoTable.grantWriteData(lambdaPutDynamoDB);`

## Deploy
Run `cdk deploy`. This will deploy / redeploy your Stack to your AWS Account.

## Testing
Run the following Lambda CLI invoke command to invoke the function. Note, you must edit the {LambdFunctionArn} placeholder with the ARN of the deployed Lambda function. This is provided in the stack outputs. Note that this requires AWS CLI v2.


```bash
aws lambda invoke --function-name "LAMBDA_FUNCTION_ARN" \
--invocation-type Event \
--payload '{ "Metadata": "Hello" }' \
--cli-binary-format raw-in-base64-out \
response.json
```

## Synthesize Cloudformation Template
To see the Cloudformation template generated by the CDK, run `cdk synth`, then check the output file in the "cdk.out" directory.
