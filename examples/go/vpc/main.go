package main

import (
	"github.com/pulumi/pulumi-metabase/sdk/go/metabase"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		metabaseService, err := metabase.NewMetabase(ctx, "metabaseService", &metabase.MetabaseArgs{
			VpcId: pulumi.String("vpc-123"),
		})
		if err != nil {
			return err
		}
		ctx.Export("url", metabaseService.DnsName)
		return nil
	})
}
