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
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/route53"
	"github.com/pulumi/pulumi-metabase/pkg/metabase"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AuthenticationStrategy string

const (
	MetabaseIdentifier = "metabase:index:Metabase"

	// The default port that the `metabase/metabase` Docker image exposes it's HTTP endpoint.
	metabasePort = 3000

	GoogleAuthentication AuthenticationStrategy = "google"
)

type MetabaseEmailConfig struct {
	Host     pulumi.StringInput `pulumi:"host"`
	Port     pulumi.IntInput    `pulumi:"port"`
	Security pulumi.StringInput `pulumi:"security"`
	Username pulumi.StringInput `pulumi:"username"`
	Password pulumi.StringInput `pulumi:"password"`
}

type CustomDomain struct {
	HostedZoneName *pulumi.StringInput `pulumi:"hostedZoneName"`
	DomainName     *pulumi.StringInput `pulumi:"domainName"`
}

type MetabaseArgs struct {
	VpcID           pulumi.StringInput      `pulumi:"vpcId"`
	ECSSubnetIDs    pulumi.StringArrayInput `pulumi:"ecsSubnetIds"`
	DBSubnetIDs     pulumi.StringArrayInput `pulumi:"dbSubnetIds"`
	LBSubnetIDs     pulumi.StringArrayInput `pulumi:"lbSubnetIds"`
	MetabaseVersion pulumi.StringInput      `pulumi:"metabaseVersion"`

	EmailConfig MetabaseEmailConfig `pulumi:"emailConfig"`

	// Additional args [WIP]
	Domain CustomDomain `pulumi:"domain"`
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

	metabaseBuilder := metabase.NewMetabaseResourceConstructor(ctx, name, opts...)

	// Security Group for the MySQL database and Metabase Task
	metabaseSecurityGroup, err := metabaseBuilder.NewMetabaseSecurityGroup(args.VpcID)
	if err != nil {
		return nil, errors.Wrap(err, "Creating Metabase Security Group")
	}

	// Security Group for the load balancer.
	loadBalancerSecurityGroup, err := metabaseBuilder.NewLoadBalancerSecurityGroup(args.VpcID, metabasePort, metabaseSecurityGroup.ID().ToStringOutput())
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
	metabaseMysqlCluster, err := metabaseBuilder.NewMySQLCluster(args.DBSubnetIDs, metabasePassword, metabaseSecurityGroup.ID())
	if err != nil {
		return nil, errors.Wrap(err, "Creating MySQL Cluster")
	}

	var certificate *acm.Certificate
	var certificateValidation *acm.CertificateValidation
	var hostedZoneID pulumi.StringOutput
	if args.Domain.DomainName != nil && args.Domain.HostedZoneName != nil {
		hostedZoneID = metabaseBuilder.GetHostedZoneId(*args.Domain.HostedZoneName)

		certificate, certificateValidation, err = metabaseBuilder.NewDomainCertificate(*args.Domain.DomainName, hostedZoneID)
		if err != nil {
			return nil, errors.Wrap(err, "Create Domain Certificate")
		}
	}

	loadBalancer, targetGroup, lbListener, err := metabaseBuilder.NewLoadBalancer(
		args.VpcID, args.LBSubnetIDs, loadBalancerSecurityGroup.ID(), metabasePort,
		certificateValidation, certificate,
	)
	if err != nil {
		return nil, err
	}

	regionName := aws.GetRegionOutput(ctx, aws.GetRegionOutputArgs{}, pulumi.Parent(component)).Name()

	metabaseImageName := pulumi.String("metabase/metabase:latest").ToStringOutput()
	if args.MetabaseVersion != nil {
		metabaseImageName = pulumi.Sprintf("metabase/metabase:%s", args.MetabaseVersion)
	}

	metabaseContainerDef := newMetabaseContainer(*metabaseMysqlCluster, metabaseImageName, regionName, args.EmailConfig)

	err = metabaseBuilder.NewMetabaseService(
		args.MetabaseVersion, regionName, metabaseContainerDef, args.ECSSubnetIDs,
		metabaseSecurityGroup.ID(), metabasePort, targetGroup.Arn, lbListener,
	)
	if err != nil {
		return nil, err
	}

	var metabaseDnsRecord *route53.Record
	if args.Domain.DomainName != nil && args.Domain.HostedZoneName != nil {
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

func newMetabaseContainer(cluster rds.Cluster, metabaseImageName, regionName pulumi.StringOutput, emailConfig MetabaseEmailConfig) pulumi.StringOutput {
	return pulumi.All(
		cluster.Endpoint, cluster.MasterUsername, cluster.MasterPassword,
		cluster.Port, cluster.DatabaseName, regionName, metabaseImageName,
		emailConfig.Host, emailConfig.Username, emailConfig.Password,
		emailConfig.Port, emailConfig.Security,
	).ApplyT(func(values []interface{}) (string, error) {
		hostname := values[0].(string)
		username := values[1].(string)
		password := values[2].(*string)
		port := values[3].(int)
		dbName := values[4].(string)
		region := values[5].(string)
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

		if values[7] != nil {
			emailHost := values[7].(string)
			emailUsername := values[8].(string)
			emailPassword := values[9].(string)
			emailPort := values[10].(int)
			emailSecurity := values[11].(string)

			metabaseEnv = append(metabaseEnv, newMetabaseEnvironmentVariable("MB_EMAIL_SMTP_USERNAME", emailUsername))
			metabaseEnv = append(metabaseEnv, newMetabaseEnvironmentVariable("MB_EMAIL_SMTP_PASSWORD", emailPassword))
			metabaseEnv = append(metabaseEnv, newMetabaseEnvironmentVariable("MB_EMAIL_SMTP_HOST", emailHost))
			metabaseEnv = append(metabaseEnv, newMetabaseEnvironmentVariable("MB_EMAIL_SMTP_PORT", fmt.Sprintf("%d", emailPort)))
			metabaseEnv = append(metabaseEnv, newMetabaseEnvironmentVariable("MB_EMAIL_SMTP_SECURITY", emailSecurity))
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
				"logConfiguration": map[string]interface{}{
					"logDriver": "awslogs",
					"options": map[string]interface{}{
						"awslogs-group":         "/ecs/metabase",
						"awslogs-region":        region,
						"awslogs-stream-prefix": "ecs",
					},
				},
			},
		})
		if err != nil {
			return "", err
		}

		return string(containerJSON), nil
	}).(pulumi.StringOutput)
}
