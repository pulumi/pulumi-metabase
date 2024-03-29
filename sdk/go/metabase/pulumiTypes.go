// Code generated by Pulumi SDK Generator DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package metabase

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Options for setting a custom domain.
type CustomDomain struct {
	DomainName     *string `pulumi:"domainName"`
	HostedZoneName *string `pulumi:"hostedZoneName"`
}

// CustomDomainInput is an input type that accepts CustomDomainArgs and CustomDomainOutput values.
// You can construct a concrete instance of `CustomDomainInput` via:
//
//          CustomDomainArgs{...}
type CustomDomainInput interface {
	pulumi.Input

	ToCustomDomainOutput() CustomDomainOutput
	ToCustomDomainOutputWithContext(context.Context) CustomDomainOutput
}

// Options for setting a custom domain.
type CustomDomainArgs struct {
	DomainName     pulumi.StringPtrInput `pulumi:"domainName"`
	HostedZoneName pulumi.StringPtrInput `pulumi:"hostedZoneName"`
}

func (CustomDomainArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*CustomDomain)(nil)).Elem()
}

func (i CustomDomainArgs) ToCustomDomainOutput() CustomDomainOutput {
	return i.ToCustomDomainOutputWithContext(context.Background())
}

func (i CustomDomainArgs) ToCustomDomainOutputWithContext(ctx context.Context) CustomDomainOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CustomDomainOutput)
}

func (i CustomDomainArgs) ToCustomDomainPtrOutput() CustomDomainPtrOutput {
	return i.ToCustomDomainPtrOutputWithContext(context.Background())
}

func (i CustomDomainArgs) ToCustomDomainPtrOutputWithContext(ctx context.Context) CustomDomainPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CustomDomainOutput).ToCustomDomainPtrOutputWithContext(ctx)
}

// CustomDomainPtrInput is an input type that accepts CustomDomainArgs, CustomDomainPtr and CustomDomainPtrOutput values.
// You can construct a concrete instance of `CustomDomainPtrInput` via:
//
//          CustomDomainArgs{...}
//
//  or:
//
//          nil
type CustomDomainPtrInput interface {
	pulumi.Input

	ToCustomDomainPtrOutput() CustomDomainPtrOutput
	ToCustomDomainPtrOutputWithContext(context.Context) CustomDomainPtrOutput
}

type customDomainPtrType CustomDomainArgs

func CustomDomainPtr(v *CustomDomainArgs) CustomDomainPtrInput {
	return (*customDomainPtrType)(v)
}

func (*customDomainPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**CustomDomain)(nil)).Elem()
}

func (i *customDomainPtrType) ToCustomDomainPtrOutput() CustomDomainPtrOutput {
	return i.ToCustomDomainPtrOutputWithContext(context.Background())
}

func (i *customDomainPtrType) ToCustomDomainPtrOutputWithContext(ctx context.Context) CustomDomainPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CustomDomainPtrOutput)
}

// Options for setting a custom domain.
type CustomDomainOutput struct{ *pulumi.OutputState }

func (CustomDomainOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*CustomDomain)(nil)).Elem()
}

func (o CustomDomainOutput) ToCustomDomainOutput() CustomDomainOutput {
	return o
}

func (o CustomDomainOutput) ToCustomDomainOutputWithContext(ctx context.Context) CustomDomainOutput {
	return o
}

func (o CustomDomainOutput) ToCustomDomainPtrOutput() CustomDomainPtrOutput {
	return o.ToCustomDomainPtrOutputWithContext(context.Background())
}

func (o CustomDomainOutput) ToCustomDomainPtrOutputWithContext(ctx context.Context) CustomDomainPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v CustomDomain) *CustomDomain {
		return &v
	}).(CustomDomainPtrOutput)
}

func (o CustomDomainOutput) DomainName() pulumi.StringPtrOutput {
	return o.ApplyT(func(v CustomDomain) *string { return v.DomainName }).(pulumi.StringPtrOutput)
}

func (o CustomDomainOutput) HostedZoneName() pulumi.StringPtrOutput {
	return o.ApplyT(func(v CustomDomain) *string { return v.HostedZoneName }).(pulumi.StringPtrOutput)
}

type CustomDomainPtrOutput struct{ *pulumi.OutputState }

func (CustomDomainPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**CustomDomain)(nil)).Elem()
}

func (o CustomDomainPtrOutput) ToCustomDomainPtrOutput() CustomDomainPtrOutput {
	return o
}

func (o CustomDomainPtrOutput) ToCustomDomainPtrOutputWithContext(ctx context.Context) CustomDomainPtrOutput {
	return o
}

func (o CustomDomainPtrOutput) Elem() CustomDomainOutput {
	return o.ApplyT(func(v *CustomDomain) CustomDomain {
		if v != nil {
			return *v
		}
		var ret CustomDomain
		return ret
	}).(CustomDomainOutput)
}

func (o CustomDomainPtrOutput) DomainName() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *CustomDomain) *string {
		if v == nil {
			return nil
		}
		return v.DomainName
	}).(pulumi.StringPtrOutput)
}

func (o CustomDomainPtrOutput) HostedZoneName() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *CustomDomain) *string {
		if v == nil {
			return nil
		}
		return v.HostedZoneName
	}).(pulumi.StringPtrOutput)
}

// The options for configuring your database.
type Database struct {
	// The database engine version. Updating this argument results in an outage. See the
	// [Aurora MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Updates.html)
	// documentation for your configured engine to determine this value. For example with Aurora MySQL 2,
	// a potential value for this argument is 5.7.mysql_aurora.2.03.2. The value can contain a partial version
	// where supported by the API.
	EngineVersion *string `pulumi:"engineVersion"`
}

// Defaults sets the appropriate defaults for Database
func (val *Database) Defaults() *Database {
	if val == nil {
		return nil
	}
	tmp := *val
	if isZero(tmp.EngineVersion) {
		engineVersion_ := "5.7.mysql_aurora.2.08.3"
		tmp.EngineVersion = &engineVersion_
	}
	return &tmp
}

// DatabaseInput is an input type that accepts DatabaseArgs and DatabaseOutput values.
// You can construct a concrete instance of `DatabaseInput` via:
//
//          DatabaseArgs{...}
type DatabaseInput interface {
	pulumi.Input

	ToDatabaseOutput() DatabaseOutput
	ToDatabaseOutputWithContext(context.Context) DatabaseOutput
}

// The options for configuring your database.
type DatabaseArgs struct {
	// The database engine version. Updating this argument results in an outage. See the
	// [Aurora MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Updates.html)
	// documentation for your configured engine to determine this value. For example with Aurora MySQL 2,
	// a potential value for this argument is 5.7.mysql_aurora.2.03.2. The value can contain a partial version
	// where supported by the API.
	EngineVersion pulumi.StringPtrInput `pulumi:"engineVersion"`
}

func (DatabaseArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Database)(nil)).Elem()
}

func (i DatabaseArgs) ToDatabaseOutput() DatabaseOutput {
	return i.ToDatabaseOutputWithContext(context.Background())
}

func (i DatabaseArgs) ToDatabaseOutputWithContext(ctx context.Context) DatabaseOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DatabaseOutput)
}

func (i DatabaseArgs) ToDatabasePtrOutput() DatabasePtrOutput {
	return i.ToDatabasePtrOutputWithContext(context.Background())
}

func (i DatabaseArgs) ToDatabasePtrOutputWithContext(ctx context.Context) DatabasePtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DatabaseOutput).ToDatabasePtrOutputWithContext(ctx)
}

// DatabasePtrInput is an input type that accepts DatabaseArgs, DatabasePtr and DatabasePtrOutput values.
// You can construct a concrete instance of `DatabasePtrInput` via:
//
//          DatabaseArgs{...}
//
//  or:
//
//          nil
type DatabasePtrInput interface {
	pulumi.Input

	ToDatabasePtrOutput() DatabasePtrOutput
	ToDatabasePtrOutputWithContext(context.Context) DatabasePtrOutput
}

type databasePtrType DatabaseArgs

func DatabasePtr(v *DatabaseArgs) DatabasePtrInput {
	return (*databasePtrType)(v)
}

func (*databasePtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**Database)(nil)).Elem()
}

func (i *databasePtrType) ToDatabasePtrOutput() DatabasePtrOutput {
	return i.ToDatabasePtrOutputWithContext(context.Background())
}

func (i *databasePtrType) ToDatabasePtrOutputWithContext(ctx context.Context) DatabasePtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DatabasePtrOutput)
}

// The options for configuring your database.
type DatabaseOutput struct{ *pulumi.OutputState }

func (DatabaseOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Database)(nil)).Elem()
}

func (o DatabaseOutput) ToDatabaseOutput() DatabaseOutput {
	return o
}

func (o DatabaseOutput) ToDatabaseOutputWithContext(ctx context.Context) DatabaseOutput {
	return o
}

func (o DatabaseOutput) ToDatabasePtrOutput() DatabasePtrOutput {
	return o.ToDatabasePtrOutputWithContext(context.Background())
}

func (o DatabaseOutput) ToDatabasePtrOutputWithContext(ctx context.Context) DatabasePtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v Database) *Database {
		return &v
	}).(DatabasePtrOutput)
}

// The database engine version. Updating this argument results in an outage. See the
// [Aurora MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Updates.html)
// documentation for your configured engine to determine this value. For example with Aurora MySQL 2,
// a potential value for this argument is 5.7.mysql_aurora.2.03.2. The value can contain a partial version
// where supported by the API.
func (o DatabaseOutput) EngineVersion() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Database) *string { return v.EngineVersion }).(pulumi.StringPtrOutput)
}

type DatabasePtrOutput struct{ *pulumi.OutputState }

func (DatabasePtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Database)(nil)).Elem()
}

func (o DatabasePtrOutput) ToDatabasePtrOutput() DatabasePtrOutput {
	return o
}

func (o DatabasePtrOutput) ToDatabasePtrOutputWithContext(ctx context.Context) DatabasePtrOutput {
	return o
}

func (o DatabasePtrOutput) Elem() DatabaseOutput {
	return o.ApplyT(func(v *Database) Database {
		if v != nil {
			return *v
		}
		var ret Database
		return ret
	}).(DatabaseOutput)
}

// The database engine version. Updating this argument results in an outage. See the
// [Aurora MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Updates.html)
// documentation for your configured engine to determine this value. For example with Aurora MySQL 2,
// a potential value for this argument is 5.7.mysql_aurora.2.03.2. The value can contain a partial version
// where supported by the API.
func (o DatabasePtrOutput) EngineVersion() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Database) *string {
		if v == nil {
			return nil
		}
		return v.EngineVersion
	}).(pulumi.StringPtrOutput)
}

// The options for networking.
type Networking struct {
	// The subnets to use for the RDS instance.
	DbSubnetIds []string `pulumi:"dbSubnetIds"`
	// The subnets to use for the Fargate task.
	EcsSubnetIds []string `pulumi:"ecsSubnetIds"`
	// The subnets to use for the load balancer.
	LbSubnetIds []string `pulumi:"lbSubnetIds"`
}

// NetworkingInput is an input type that accepts NetworkingArgs and NetworkingOutput values.
// You can construct a concrete instance of `NetworkingInput` via:
//
//          NetworkingArgs{...}
type NetworkingInput interface {
	pulumi.Input

	ToNetworkingOutput() NetworkingOutput
	ToNetworkingOutputWithContext(context.Context) NetworkingOutput
}

// The options for networking.
type NetworkingArgs struct {
	// The subnets to use for the RDS instance.
	DbSubnetIds pulumi.StringArrayInput `pulumi:"dbSubnetIds"`
	// The subnets to use for the Fargate task.
	EcsSubnetIds pulumi.StringArrayInput `pulumi:"ecsSubnetIds"`
	// The subnets to use for the load balancer.
	LbSubnetIds pulumi.StringArrayInput `pulumi:"lbSubnetIds"`
}

func (NetworkingArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Networking)(nil)).Elem()
}

func (i NetworkingArgs) ToNetworkingOutput() NetworkingOutput {
	return i.ToNetworkingOutputWithContext(context.Background())
}

func (i NetworkingArgs) ToNetworkingOutputWithContext(ctx context.Context) NetworkingOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NetworkingOutput)
}

func (i NetworkingArgs) ToNetworkingPtrOutput() NetworkingPtrOutput {
	return i.ToNetworkingPtrOutputWithContext(context.Background())
}

func (i NetworkingArgs) ToNetworkingPtrOutputWithContext(ctx context.Context) NetworkingPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NetworkingOutput).ToNetworkingPtrOutputWithContext(ctx)
}

// NetworkingPtrInput is an input type that accepts NetworkingArgs, NetworkingPtr and NetworkingPtrOutput values.
// You can construct a concrete instance of `NetworkingPtrInput` via:
//
//          NetworkingArgs{...}
//
//  or:
//
//          nil
type NetworkingPtrInput interface {
	pulumi.Input

	ToNetworkingPtrOutput() NetworkingPtrOutput
	ToNetworkingPtrOutputWithContext(context.Context) NetworkingPtrOutput
}

type networkingPtrType NetworkingArgs

func NetworkingPtr(v *NetworkingArgs) NetworkingPtrInput {
	return (*networkingPtrType)(v)
}

func (*networkingPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**Networking)(nil)).Elem()
}

func (i *networkingPtrType) ToNetworkingPtrOutput() NetworkingPtrOutput {
	return i.ToNetworkingPtrOutputWithContext(context.Background())
}

func (i *networkingPtrType) ToNetworkingPtrOutputWithContext(ctx context.Context) NetworkingPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(NetworkingPtrOutput)
}

// The options for networking.
type NetworkingOutput struct{ *pulumi.OutputState }

func (NetworkingOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Networking)(nil)).Elem()
}

func (o NetworkingOutput) ToNetworkingOutput() NetworkingOutput {
	return o
}

func (o NetworkingOutput) ToNetworkingOutputWithContext(ctx context.Context) NetworkingOutput {
	return o
}

func (o NetworkingOutput) ToNetworkingPtrOutput() NetworkingPtrOutput {
	return o.ToNetworkingPtrOutputWithContext(context.Background())
}

func (o NetworkingOutput) ToNetworkingPtrOutputWithContext(ctx context.Context) NetworkingPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v Networking) *Networking {
		return &v
	}).(NetworkingPtrOutput)
}

// The subnets to use for the RDS instance.
func (o NetworkingOutput) DbSubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v Networking) []string { return v.DbSubnetIds }).(pulumi.StringArrayOutput)
}

// The subnets to use for the Fargate task.
func (o NetworkingOutput) EcsSubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v Networking) []string { return v.EcsSubnetIds }).(pulumi.StringArrayOutput)
}

// The subnets to use for the load balancer.
func (o NetworkingOutput) LbSubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v Networking) []string { return v.LbSubnetIds }).(pulumi.StringArrayOutput)
}

type NetworkingPtrOutput struct{ *pulumi.OutputState }

func (NetworkingPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Networking)(nil)).Elem()
}

func (o NetworkingPtrOutput) ToNetworkingPtrOutput() NetworkingPtrOutput {
	return o
}

func (o NetworkingPtrOutput) ToNetworkingPtrOutputWithContext(ctx context.Context) NetworkingPtrOutput {
	return o
}

func (o NetworkingPtrOutput) Elem() NetworkingOutput {
	return o.ApplyT(func(v *Networking) Networking {
		if v != nil {
			return *v
		}
		var ret Networking
		return ret
	}).(NetworkingOutput)
}

// The subnets to use for the RDS instance.
func (o NetworkingPtrOutput) DbSubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *Networking) []string {
		if v == nil {
			return nil
		}
		return v.DbSubnetIds
	}).(pulumi.StringArrayOutput)
}

// The subnets to use for the Fargate task.
func (o NetworkingPtrOutput) EcsSubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *Networking) []string {
		if v == nil {
			return nil
		}
		return v.EcsSubnetIds
	}).(pulumi.StringArrayOutput)
}

// The subnets to use for the load balancer.
func (o NetworkingPtrOutput) LbSubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *Networking) []string {
		if v == nil {
			return nil
		}
		return v.LbSubnetIds
	}).(pulumi.StringArrayOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*CustomDomainInput)(nil)).Elem(), CustomDomainArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*CustomDomainPtrInput)(nil)).Elem(), CustomDomainArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*DatabaseInput)(nil)).Elem(), DatabaseArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*DatabasePtrInput)(nil)).Elem(), DatabaseArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*NetworkingInput)(nil)).Elem(), NetworkingArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*NetworkingPtrInput)(nil)).Elem(), NetworkingArgs{})
	pulumi.RegisterOutputType(CustomDomainOutput{})
	pulumi.RegisterOutputType(CustomDomainPtrOutput{})
	pulumi.RegisterOutputType(DatabaseOutput{})
	pulumi.RegisterOutputType(DatabasePtrOutput{})
	pulumi.RegisterOutputType(NetworkingOutput{})
	pulumi.RegisterOutputType(NetworkingPtrOutput{})
}
