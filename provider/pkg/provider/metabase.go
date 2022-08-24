// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/route53"
	"github.com/pulumi/pulumi-metabase/pkg/metabase"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	MetabaseIdentifier = "metabase:index:Metabase"

	// The default port that the `metabase/metabase` Docker image exposes it's HTTP endpoint.
	metabasePort = 3000
)

type CustomDomain struct {
	HostedZoneName *pulumi.StringInput `pulumi:"hostedZoneName"`
	DomainName     *pulumi.StringInput `pulumi:"domainName"`
}

type Networking struct {
	ECSSubnetIDs pulumi.StringArrayInput `pulumi:"ecsSubnetIds"`
	DBSubnetIDs  pulumi.StringArrayInput `pulumi:"dbSubnetIds"`
	LBSubnetIDs  pulumi.StringArrayInput `pulumi:"lbSubnetIds"`
}

type Database struct {
	EngineVersion pulumi.StringInput `pulumi:"engineVersion"`
}

type MetabaseArgs struct {
	VpcID           pulumi.StringInput `pulumi:"vpcId"`
	MetabaseVersion pulumi.StringInput `pulumi:"metabaseVersion"`

	// Additional args
	Domain   CustomDomain `pulumi:"domain"`
	Network  Networking   `pulumi:"networking"`
	Database Database     `pulumi:"database"`
}

type Metabase struct {
	pulumi.ResourceState

	DNSName         pulumi.StringOutput `pulumi:"dnsName"`
	SecurityGroupID pulumi.StringOutput `pulumi:"securityGroupId"`
}

func NewMetabase(ctx *pulumi.Context, name string, args *MetabaseArgs, opts ...pulumi.ResourceOption) (*Metabase, error) {
	if args == nil {
		args = &MetabaseArgs{}
	}

	component := &Metabase{}
	err := ctx.RegisterComponentResource(MetabaseIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	attachDomainName := args.Domain.DomainName != nil && args.Domain.HostedZoneName != nil

	metabaseBuilder := metabase.NewMetabaseResourceConstructor(ctx, name, opts...)

	vpcID := args.VpcID
	if vpcID == nil {
		vpc, err := ec2.NewDefaultVpc(ctx, name, &ec2.DefaultVpcArgs{}, opts...)
		if err != nil {
			return nil, err
		}
		vpcID = vpc.ID()
	}

	// If the network options are not provided we just run everything in the same
	// public subnets. If there are exactly two public subnets available we throw
	// an error letting the user know they need to configure their VPC.
	defaultSubnetIDs := vpcID.ToStringOutput().ApplyT(func(id string) ([]string, error) {
		subnets, err := ec2.GetSubnets(ctx, &ec2.GetSubnetsArgs{
			Filters: []ec2.GetSubnetsFilter{
				{
					Name:   "vpc-id",
					Values: []string{id},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		var result []string
		azMap := make(map[string]string, len(subnets.Ids))
		for _, subnetID := range subnets.Ids {
			s, err := ec2.LookupSubnet(ctx, &ec2.LookupSubnetArgs{
				Id: &subnetID,
			})
			if err != nil {
				return nil, err
			}

			_, ok := azMap[s.AvailabilityZone]
			if !ok && s.MapPublicIpOnLaunch {
				azMap[s.AvailabilityZone] = s.AvailabilityZone
				result = append(result, subnetID)

				if len(result) == 2 {
					break
				}
			}
		}

		if len(result) != 2 {
			return nil, fmt.Errorf("Your VPC must have at least two public subnets available. You will need to correctly configure your VPC or alternatively provide specific subnet ids in the 'Networking' options.")
		}

		return result, nil
	}).(pulumi.StringArrayOutput)

	lbSubnetIDs := args.Network.LBSubnetIDs
	if lbSubnetIDs == nil {
		lbSubnetIDs = defaultSubnetIDs
	}

	ecsAssignPublicIP := false
	ecsSubnetIDs := args.Network.ECSSubnetIDs
	if ecsSubnetIDs == nil {
		ecsAssignPublicIP = true
		ecsSubnetIDs = defaultSubnetIDs
	}

	dbSubnetIDs := args.Network.DBSubnetIDs
	if dbSubnetIDs == nil {
		dbSubnetIDs = defaultSubnetIDs
	}

	// Security Group for the MySQL database and Metabase Task
	metabaseSecurityGroup, err := metabaseBuilder.NewMetabaseSecurityGroup(vpcID)
	if err != nil {
		return nil, errors.Wrap(err, "Creating Metabase Security Group")
	}

	// Security Group for the load balancer.
	loadBalancerSecurityGroup, err := metabaseBuilder.NewLoadBalancerSecurityGroup(vpcID, metabasePort, metabaseSecurityGroup.ID().ToStringOutput())
	if err != nil {
		return nil, errors.Wrap(err, "Creating Load Balancer Security Group")
	}

	// Create the security group rules.
	err = metabaseBuilder.NewSecurityGroupRules(metabasePort, metabaseSecurityGroup.ID(), loadBalancerSecurityGroup.ID())
	if err != nil {
		return nil, errors.Wrap(err, "Creating Security Group Rules")
	}

	// Create a password for the MySQL cluster.
	metabasePassword, err := metabaseBuilder.NewMetabasePassword()
	if err != nil {
		return nil, errors.Wrap(err, "Creating Metabase Password")
	}

	// Create the MySQL cluster.
	metabaseMysqlCluster, err := metabaseBuilder.NewMySQLCluster(dbSubnetIDs, metabasePassword, metabaseSecurityGroup.ID(), args.Database.EngineVersion)
	if err != nil {
		return nil, errors.Wrap(err, "Creating MySQL Cluster")
	}

	var certificate *acm.Certificate
	var certificateValidation *acm.CertificateValidation
	var hostedZoneID pulumi.StringOutput
	if attachDomainName {
		hostedZoneID = metabaseBuilder.GetHostedZoneId(*args.Domain.HostedZoneName)

		certificate, certificateValidation, err = metabaseBuilder.NewDomainCertificate(*args.Domain.DomainName, hostedZoneID)
		if err != nil {
			return nil, errors.Wrap(err, "Create Domain Certificate")
		}
	}

	loadBalancer, targetGroup, lbListener, err := metabaseBuilder.NewLoadBalancer(
		vpcID, lbSubnetIDs, loadBalancerSecurityGroup.ID(), metabasePort,
		certificateValidation, certificate, attachDomainName,
	)
	if err != nil {
		return nil, err
	}

	regionName := aws.GetRegionOutput(ctx, aws.GetRegionOutputArgs{}, pulumi.Parent(component)).Name()

	metabaseImageName := pulumi.String("metabase/metabase:latest").ToStringOutput()
	if args.MetabaseVersion != nil {
		metabaseImageName = pulumi.Sprintf("metabase/metabase:%s", args.MetabaseVersion)
	}

	metabaseContainerDef := newMetabaseContainer(*metabaseMysqlCluster, metabaseImageName, regionName)

	err = metabaseBuilder.NewMetabaseService(
		args.MetabaseVersion, regionName, metabaseContainerDef, ecsSubnetIDs,
		metabaseSecurityGroup.ID(), metabasePort, targetGroup.Arn, lbListener,
		ecsAssignPublicIP,
	)
	if err != nil {
		return nil, err
	}

	var metabaseDnsRecord *route53.Record
	if attachDomainName {
		metabaseDnsRecord, err = route53.NewRecord(ctx, "metabase-dns", &route53.RecordArgs{
			ZoneId: hostedZoneID,
			Type:   pulumi.String("A"),
			Name:   *args.Domain.DomainName,
			Aliases: route53.RecordAliasArray{
				route53.RecordAliasArgs{
					Name:                 loadBalancer.DnsName,
					ZoneId:               loadBalancer.ZoneId,
					EvaluateTargetHealth: pulumi.Bool(true),
				},
			},
		}, opts...)
		if err != nil {
			return nil, err
		}
	}

	component.SecurityGroupID = metabaseSecurityGroup.ID().ToStringOutput()

	component.DNSName = loadBalancer.DnsName
	if metabaseDnsRecord != nil {
		component.DNSName = pulumi.Sprintf("https://%s", metabaseDnsRecord.Name)
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"securityGroupId": component.SecurityGroupID,
		"dnsName":         component.DNSName,
	}); err != nil {
		return nil, err
	}

	return component, nil
}

type metabaseEnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func newMetabaseEnvironmentVariable(name, value string) metabaseEnvironmentVariable {
	return metabaseEnvironmentVariable{Name: name, Value: value}
}

func newMetabaseContainer(cluster rds.Cluster, metabaseImageName, regionName pulumi.StringOutput) pulumi.StringOutput {
	return pulumi.All(
		cluster.Endpoint, cluster.MasterUsername, cluster.MasterPassword,
		cluster.Port, cluster.DatabaseName, regionName, metabaseImageName,
	).ApplyT(func(values []interface{}) (string, error) {
		hostname := values[0].(string)
		username := values[1].(string)
		password := values[2].(*string)
		port := values[3].(int)
		dbName := values[4].(string)
		//region := values[5].(string)
		imageName := values[6].(string)

		metabaseEnv := []metabaseEnvironmentVariable{
			newMetabaseEnvironmentVariable("JAVA_TIMEZONE", "US/Pacific"),
			newMetabaseEnvironmentVariable("MB_DB_TYPE", "mysql"),
			newMetabaseEnvironmentVariable("MB_DB_DBNAME", dbName),
			newMetabaseEnvironmentVariable("MB_DB_PORT", fmt.Sprintf("%d", port)),
			newMetabaseEnvironmentVariable("MB_DB_USER", username),
			newMetabaseEnvironmentVariable("MB_DB_PASS", *password),
			newMetabaseEnvironmentVariable("MB_DB_HOST", hostname),
		}

		containerJSON, err := json.Marshal([]interface{}{
			map[string]interface{}{
				"name":  "metabase",
				"image": imageName,
				"portMappings": []map[string]interface{}{
					{
						"containerPort": metabasePort,
					},
				},
				"environment": metabaseEnv,
				// "logConfiguration": map[string]interface{}{
				// 	"logDriver": "awslogs",
				// 	"options": map[string]interface{}{
				// 		"awslogs-group":         "/ecs/metabase",
				// 		"awslogs-region":        region,
				// 		"awslogs-stream-prefix": "ecs",
				// 	},
				// },
			},
		})
		if err != nil {
			return "", err
		}

		return string(containerJSON), nil
	}).(pulumi.StringOutput)
}
