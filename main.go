package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type fileProcessor struct {
	pulumi.ResourceState

	Bucket *s3.Bucket
	Lambda *lambda.Function
	Table  *dynamodb.Table
}

func NewFileProcessor(ctx *pulumi.Context, name string, args *fileProcessorArgs, opts ...pulumi.ResourceOption) (*fileProcessor, error) {

	component := &fileProcessor{}
	err := ctx.RegisterComponentResource("components:fileProcessor:fileProcessor", name, component, opts...)
	if err != nil {
		return nil, err
	}

	bucket, err := s3.NewBucket(ctx, name, &s3.BucketArgs{
		Bucket: args.BucketName,
	})
	if err != nil {
		return nil, err
	}

	bucketPolicy := pulumi.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}`, bucket.ID()).ToStringOutput()

	dynamoPolicy := `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Action": ["dynamodb:PutItem"],
			"Resource": "*"
		}]
	}`

	logPolicy := pulumi.String(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"logs:CreateLogGroup",
					"logs:CreateLogStream",
					"logs:PutLogEvents",
					"logs:DescribeLogGroups",
					"logs:DescribeLogStreams"
				],
				"Resource": "*"
			}
		]
	}`)

	lambdaRole, err := iam.NewRole(ctx, fmt.Sprintf("%s-role", name), &iam.RoleArgs{
		Name: pulumi.String(name),
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": {
				"Effect": "Allow",
				"Principal": {"Service": "lambda.amazonaws.com"},
				"Action": "sts:AssumeRole"
			}
		}`),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("%s-bucketPolicy", name), &iam.RolePolicyArgs{
		Role:   lambdaRole.Name,
		Policy: bucketPolicy,
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("%s-dynamoPolicy", name), &iam.RolePolicyArgs{
		Role:   lambdaRole.Name,
		Policy: pulumi.String(dynamoPolicy),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("%s-logPolicy", name), &iam.RolePolicyArgs{
		Role:   lambdaRole.Name,
		Policy: pulumi.String(logPolicy),
	})
	if err != nil {
		return nil, err
	}

	function, err := lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Name:    pulumi.String(name),
		Handler: pulumi.String("index.handler"),
		Role:    lambdaRole.Arn,
		Runtime: pulumi.String("nodejs16.x"),
		Code:    pulumi.NewFileArchive("./lambda"),
		Environment: lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{"TABLE_NAME": args.TableName},
		},
	}, pulumi.DependsOn([]pulumi.Resource{bucket}))
	if err != nil {
		return nil, err
	}

	bucketLambdaPermission, err := lambda.NewPermission(ctx, fmt.Sprintf("%s-permission", name), &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  function.Name,
		Principal: pulumi.String("s3.amazonaws.com"),
		SourceArn: bucket.Arn,
	})
	if err != nil {
		return nil, err
	}

	_, err = s3.NewBucketNotification(ctx, fmt.Sprintf("%s-objectCreated", name), &s3.BucketNotificationArgs{
		Bucket: bucket.ID(),
		LambdaFunctions: s3.BucketNotificationLambdaFunctionArray{
			&s3.BucketNotificationLambdaFunctionArgs{
				LambdaFunctionArn: function.Arn,
				Events:            pulumi.StringArray{pulumi.String("s3:ObjectCreated:*")},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{bucketLambdaPermission}))
	if err != nil {
		return nil, err
	}

	table, err := dynamodb.NewTable(ctx, name, &dynamodb.TableArgs{
		Name: pulumi.String(name),
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("Key"),
				Type: pulumi.String("S"),
			},
		},
		HashKey:        pulumi.String("Key"),
		BillingMode:    pulumi.String("PAY_PER_REQUEST"),
		StreamEnabled:  pulumi.Bool(true),
		StreamViewType: pulumi.String("NEW_IMAGE"),
	})
	if err != nil {
		return nil, err
	}

	component.Bucket = bucket
	component.Lambda = function
	component.Table = table

	ctx.RegisterResourceOutputs(component, pulumi.Map{
		"bucket": component.Bucket,
		"lambda": component.Lambda,
		"table":  component.Table,
	})

	return component, nil
}

type fileProcessorArgs struct {
	TableName    pulumi.StringInput
	BucketName   pulumi.StringInput
	FunctionName pulumi.StringInput
}

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		config := config.New(ctx, "")

		table := pulumi.String(config.Require("table"))
		bucket := pulumi.String(config.Require("bucket"))
		function := pulumi.String(config.Require("function"))

		_, err := NewFileProcessor(ctx, "sync-stream", &fileProcessorArgs{
			TableName:    table,
			BucketName:   bucket,
			FunctionName: function,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
