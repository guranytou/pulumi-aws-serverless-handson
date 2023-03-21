package main

import (
	"io"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		tableName := "users"
		_, err := dynamodb.NewTable(ctx, tableName, &dynamodb.TableArgs{
			Attributes: dynamodb.TableAttributeArray{
				&dynamodb.TableAttributeArgs{
					Name: pulumi.String("id"),
					Type: pulumi.String("S"),
				},
			},
			HashKey:       pulumi.String("id"),
			ReadCapacity:  pulumi.Int(5),
			WriteCapacity: pulumi.Int(5),
			Name:          pulumi.String(tableName),
		})
		if err != nil {
			return err
		}

		f, err := os.Open("assume-role-policy.json")
		if err != nil {
			return err
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		roleName := "users-role"
		role, err := iam.NewRole(ctx, roleName, &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(b),
			Name:             pulumi.String(roleName),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicyAttachment(ctx, roleName+"-1", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
			Role:      role.Name,
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicyAttachment(ctx, roleName+"-2", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"),
			Role:      role.Name,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
