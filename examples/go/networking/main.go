package main

import (
	"github.com/pulumi/pulumi-metabase/sdk/go/metabase"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		metabaseService, err := metabase.NewMetabase(ctx, "metabaseService", &metabase.MetabaseArgs{
			VpcId: pulumi.String("vpc-123"),
			Networking: &metabase.NetworkingArgs{
				EcsSubnetIds: pulumi.StringArray{
					pulumi.String("subnet-123"),
					pulumi.String("subnet-456"),
				},
				DbSubnetIds: pulumi.StringArray{
					pulumi.String("subnet-789"),
					pulumi.String("subnet-abc"),
				},
				LbSubnetIds: pulumi.StringArray{
					pulumi.String("subnet-def"),
					pulumi.String("subnet-ghi"),
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("url", metabaseService.DnsName)
		return nil
	})
}
