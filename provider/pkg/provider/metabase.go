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

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/route53"
	random "github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	MetabaseIdentifier = "metabase:index:Metabase"

	// The default port that the `metabase/metabase` Docker image exposes it's HTTP endpoint.
	metabasePort = 3000
)

type MetabaseEmailConfig struct {
	Host     pulumi.StringInput `pulumi:"host"`
	Port     pulumi.IntInput    `pulumi:"port"`
	Security pulumi.StringInput `pulumi:"security"`
	Username pulumi.StringInput `pulumi:"username"`
	Password pulumi.StringInput `pulumi:"password"`
}

type MetabaseArgs struct {
	VpcID            pulumi.StringInput      `pulumi:"vpcId"`
	ECSSubnetIDs     pulumi.StringArrayInput `pulumi:"ecsSubnetIDs"`
	DBSubnetIDs      pulumi.StringArrayInput `pulumi:"dbSubnetIds"`
	LBSubnetIDs      pulumi.StringArrayInput `pulumi:"lbSubnetIds"`
	HostedZoneName   pulumi.StringInput      `pulumi:"hostedZoneName"`
	DomainName       pulumi.StringInput      `pulumi:"domainName"`
	MetabaseVersion  pulumi.StringInput      `pulumi:"metabaseVersion"`
	EmailConfig      MetabaseEmailConfig     `pulumi:"emailConfig"`
	OIDCClientId     pulumi.StringInput      `pulumi:"oidcClientId"`
	OIDCClientSecret pulumi.StringInput      `pulumi:"oidcClientSecret"`
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

	// Security Group for the MySQL database and Metabase Task
	metabaseSecurityGroupName := fmt.Sprintf("%s-metabase-sg", name)
	metabaseSecurityGroup, err := ec2.NewSecurityGroup(ctx, metabaseSecurityGroupName, &ec2.SecurityGroupArgs{
		VpcId: args.VpcID,
	}, opts...)
	if err != nil {
		return nil, err
	}

	loadBalancerSecurityGroupName := fmt.Sprintf("%s-metabase-lb-sg", name)
	loadBalancerSecurityGroup, err := ec2.NewSecurityGroup(ctx, loadBalancerSecurityGroupName, &ec2.SecurityGroupArgs{
		VpcId: args.VpcID,
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				ToPort:     pulumi.Int(443),
				FromPort:   pulumi.Int(443),
				CidrBlocks: pulumi.ToStringArray([]string{"0.0.0.0/0"}),
			},
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				ToPort:     pulumi.Int(80),
				FromPort:   pulumi.Int(80),
				CidrBlocks: pulumi.ToStringArray([]string{"0.0.0.0/0"}),
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			ec2.SecurityGroupEgressArgs{
				Protocol:       pulumi.String("tcp"),
				ToPort:         pulumi.Int(metabasePort),
				FromPort:       pulumi.Int(metabasePort),
				SecurityGroups: pulumi.ToStringArrayOutput([]pulumi.StringOutput{metabaseSecurityGroup.ID().ToStringOutput()}),
			},
			ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("tcp"),
				ToPort:     pulumi.Int(443),
				FromPort:   pulumi.Int(443),
				CidrBlocks: pulumi.ToStringArray([]string{"0.0.0.0/0"}),
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseSecurityGroupSegmentRuleName := fmt.Sprintf("%s-metabase-segment", name)
	_, err = ec2.NewSecurityGroupRule(ctx, metabaseSecurityGroupSegmentRuleName, &ec2.SecurityGroupRuleArgs{
		Description:           pulumi.String("Allow access to Metabase from the Load Balancer"),
		SecurityGroupId:       metabaseSecurityGroup.ID(),
		Type:                  pulumi.String("ingress"),
		Protocol:              pulumi.String("tcp"),
		FromPort:              pulumi.Int(metabasePort),
		ToPort:                pulumi.Int(metabasePort),
		SourceSecurityGroupId: loadBalancerSecurityGroup.ID(),
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseSecurityGroupSelfRuleName := fmt.Sprintf("%s-metabase-self", name)
	_, err = ec2.NewSecurityGroupRule(ctx, metabaseSecurityGroupSelfRuleName, &ec2.SecurityGroupRuleArgs{
		Description:           pulumi.String("Allow access to anything from within the Security Group"),
		SecurityGroupId:       metabaseSecurityGroup.ID(),
		Type:                  pulumi.String("ingress"),
		Protocol:              pulumi.String("tcp"),
		FromPort:              pulumi.Int(0),
		ToPort:                pulumi.Int(65535),
		SourceSecurityGroupId: metabaseSecurityGroup.ID(),
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseSecurityGroupEgressRuleName := fmt.Sprintf("%s-metabase-egress", name)
	_, err = ec2.NewSecurityGroupRule(ctx, metabaseSecurityGroupEgressRuleName, &ec2.SecurityGroupRuleArgs{
		Description:     pulumi.String("Allow egress to anywhere"),
		SecurityGroupId: metabaseSecurityGroup.ID(),
		Type:            pulumi.String("egress"),
		Protocol:        pulumi.String("tcp"),
		FromPort:        pulumi.Int(0),
		ToPort:          pulumi.Int(65535),
		CidrBlocks:      pulumi.ToStringArray([]string{"0.0.0.0/0"}),
	}, opts...)
	if err != nil {
		return nil, err
	}

	// Aurora Serverless MySQL for Metabase query/dashboard state
	metabaseResourceName := fmt.Sprintf("%s-metabase", name)
	metabasePassword, err := random.NewRandomString(ctx, metabaseResourceName, &random.RandomStringArgs{
		Special: pulumi.BoolPtr(false),
		Length:  pulumi.Int(20),
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseMysqlSubnetGroup, err := rds.NewSubnetGroup(ctx, metabaseResourceName, &rds.SubnetGroupArgs{
		SubnetIds: args.DBSubnetIDs,
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseMysqlCluster, err := rds.NewCluster(ctx, metabaseResourceName, &rds.ClusterArgs{
		ClusterIdentifier:       pulumi.Sprintf("%smetabasemysql", name),
		DatabaseName:            pulumi.String("metabase"),
		MasterUsername:          pulumi.String("admin"),
		MasterPassword:          metabasePassword.Result,
		Engine:                  pulumi.String("aurora"),
		EngineMode:              pulumi.String("serverless"),
		EngineVersion:           pulumi.String("5.6.10a"),
		DbSubnetGroupName:       metabaseMysqlSubnetGroup.Name,
		VpcSecurityGroupIds:     pulumi.ToStringArrayOutput([]pulumi.StringOutput{metabaseSecurityGroup.ID().ToStringOutput()}),
		FinalSnapshotIdentifier: pulumi.Sprintf("%smetabasefinalsnapshot", name),
		EnableHttpEndpoint:      pulumi.BoolPtr(true),
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseMySqlSnapshot, err := rds.NewClusterSnapshot(ctx, metabaseResourceName, &rds.ClusterSnapshotArgs{
		DbClusterIdentifier:         metabaseMysqlCluster.ClusterIdentifier,
		DbClusterSnapshotIdentifier: pulumi.String("snaphostfor57migration"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseMysql57ClusterName := fmt.Sprintf("%smetabase57", name)
	metabaseMysql57Cluster, err := rds.NewCluster(ctx, metabaseMysql57ClusterName, &rds.ClusterArgs{
		SnapshotIdentifier:      metabaseMySqlSnapshot.ID(),
		ClusterIdentifier:       pulumi.Sprintf("%smetabasemysql57", name),
		DatabaseName:            pulumi.String("metabase"),
		MasterUsername:          pulumi.String("admin"),
		MasterPassword:          metabasePassword.Result,
		Engine:                  pulumi.String("aurora-mysql"),
		EngineMode:              pulumi.String("serverless"),
		EngineVersion:           pulumi.String("5.7.mysql_aurora.2.07.1"),
		DbSubnetGroupName:       metabaseMysqlSubnetGroup.Name,
		VpcSecurityGroupIds:     pulumi.ToStringArrayOutput([]pulumi.StringOutput{metabaseSecurityGroup.ID().ToStringOutput()}),
		FinalSnapshotIdentifier: pulumi.Sprintf("%s57metabasefinalsnapshot", name),
		EnableHttpEndpoint:      pulumi.BoolPtr(true),
	}, opts...)
	if err != nil {
		return nil, err
	}

	certificate, err := acm.NewCertificate(ctx, metabaseResourceName, &acm.CertificateArgs{
		DomainName:       args.DomainName,
		ValidationMethod: pulumi.String("DNS"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	hostedZoneID := args.HostedZoneName.ToStringOutput().ApplyT(func(name string) (string, error) {
		hostedZone, err := route53.LookupZone(ctx, &route53.LookupZoneArgs{
			Name: &name,
		})
		if err != nil {
			return "", nil
		}
		return hostedZone.Id, nil
	}).(pulumi.StringOutput)

	certificateValidationRecordName := fmt.Sprintf("%s-metabase-certvalidation", name)
	certificateValidationRecord, err := route53.NewRecord(ctx, certificateValidationRecordName, &route53.RecordArgs{
		Name: certificate.ValidationOptions.ApplyT(func(opts []acm.CertificateDomainValidationOption) string {
			return *opts[0].ResourceRecordName
		}).(pulumi.StringOutput),
		Type: certificate.ValidationOptions.ApplyT(func(opts []acm.CertificateDomainValidationOption) string {
			return *opts[0].ResourceRecordType
		}).(pulumi.StringOutput),
		ZoneId: hostedZoneID,
		Records: certificate.ValidationOptions.ApplyT(func(opts []acm.CertificateDomainValidationOption) []string {
			return []string{*opts[0].ResourceRecordType}
		}).(pulumi.StringArrayOutput),
		Ttl: pulumi.IntPtr(60),
	}, opts...)
	if err != nil {
		return nil, err
	}

	certificateValidation, err := acm.NewCertificateValidation(ctx, metabaseResourceName, &acm.CertificateValidationArgs{
		CertificateArn:        certificate.Arn,
		ValidationRecordFqdns: pulumi.ToStringArrayOutput([]pulumi.StringOutput{certificateValidationRecord.Fqdn}),
	}, opts...)
	if err != nil {
		return nil, err
	}

	// Stable load balancer endpoint (no other way to get a consistent IP for an ECS service!!!)
	loadBalancer, err := lb.NewLoadBalancer(ctx, metabaseResourceName, &lb.LoadBalancerArgs{
		LoadBalancerType: pulumi.String("application"),
		Subnets:          args.LBSubnetIDs,
		SecurityGroups:   pulumi.ToStringArrayOutput([]pulumi.StringOutput{loadBalancerSecurityGroup.ID().ToStringOutput()}),
		IdleTimeout:      pulumi.IntPtr(600),
	}, opts...)
	if err != nil {
		return nil, err
	}

	targetGroup, err := lb.NewTargetGroup(ctx, metabaseResourceName, &lb.TargetGroupArgs{
		TargetType: pulumi.String("ip"),
		Port:       pulumi.Int(metabasePort),
		Protocol:   pulumi.String("HTTP"),
		VpcId:      args.VpcID,
		// Since this is a user facing tool, and we only have 0 or 1 running instances, we don't need to wait to
		// drain connections, and instead want to ensure we have as little downtime as possible.
		DeregistrationDelay: pulumi.Int(0),
	}, opts...)
	if err != nil {
		return nil, err
	}

	listenerOpts := append(opts, pulumi.DependsOn([]pulumi.Resource{certificateValidation}))
	_, err = lb.NewListener(ctx, metabaseResourceName, &lb.ListenerArgs{
		LoadBalancerArn: loadBalancer.Arn,
		Port:            pulumi.Int(443),
		Protocol:        pulumi.String("HTTP"),
		CertificateArn:  certificate.Arn,
		SslPolicy:       pulumi.String("ELBSecurityPolicy-TLS-1-2-2017-01"),
		DefaultActions: lb.ListenerDefaultActionArray{
			lb.ListenerDefaultActionArgs{
				Type:  pulumi.String("authenticate-oidc"),
				Order: pulumi.IntPtr(1),
				AuthenticateOidc: &lb.ListenerDefaultActionAuthenticateOidcArgs{
					OnUnauthenticatedRequest: pulumi.String("authenticate"),
					Issuer:                   pulumi.String("https://accounts.google.com"),
					AuthorizationEndpoint:    pulumi.String("https://accounts.google.com/o/oauth2/v2/auth"),
					TokenEndpoint:            pulumi.String("https://oauth2.googleapis.com/token"),
					UserInfoEndpoint:         pulumi.String("https://openidconnect.googleapis.com/v1/userinfo"),
					ClientId:                 args.OIDCClientId,
					ClientSecret:             args.OIDCClientSecret,
				},
			},
			lb.ListenerDefaultActionArgs{
				Type:           pulumi.String("forward"),
				Order:          pulumi.IntPtr(2),
				TargetGroupArn: targetGroup.Arn,
			},
		},
	}, listenerOpts...)
	if err != nil {
		return nil, err
	}

	httpRedirectListenerName := fmt.Sprintf("%s-metabase-redirecthttp", name)
	_, err = lb.NewListener(ctx, httpRedirectListenerName, &lb.ListenerArgs{
		LoadBalancerArn: loadBalancer.Arn,
		Port:            pulumi.IntPtr(80),
		Protocol:        pulumi.String("HTTP"),
		DefaultActions: lb.ListenerDefaultActionArray{
			lb.ListenerDefaultActionArgs{
				Type: pulumi.String("redirect"),
				Redirect: &lb.ListenerDefaultActionRedirectArgs{
					Protocol:   pulumi.String("HTTPS"),
					Port:       pulumi.String("443"),
					StatusCode: pulumi.String("HTTP_301"),
				},
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseImageName := args.MetabaseVersion.ToStringOutput().ApplyT(func(version string) string {
		if version == "" {
			version = "latest"
		}
		return fmt.Sprintf("metabase/metabase:%s", version)
	}).(pulumi.StringOutput)

	metabaseCluster, err := ecs.NewCluster(ctx, metabaseResourceName, &ecs.ClusterArgs{}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseExecutionRole, err := iam.GetRole(ctx, metabaseResourceName, pulumi.ID("ecsTaskExecutionRole"), nil, opts...)
	if err != nil {
		return nil, err
	}

	regionName := aws.GetRegionOutput(ctx, aws.GetRegionOutputArgs{}, pulumi.Parent(component)).Name()

	metabaseTaskDefinition, err := ecs.NewTaskDefinition(ctx, metabaseResourceName, &ecs.TaskDefinitionArgs{
		Family:                  pulumi.String("metabase"),
		Cpu:                     pulumi.String("2048"),
		Memory:                  pulumi.String("8192"),
		RequiresCompatibilities: pulumi.ToStringArray([]string{"FARGATE"}),
		NetworkMode:             pulumi.StringPtr("awsvpc"),
		ExecutionRoleArn:        metabaseExecutionRole.Arn,
		ContainerDefinitions:    newMetabaseContainer(*metabaseMysql57Cluster, metabaseImageName, regionName, args.EmailConfig),
	}, opts...)
	if err != nil {
		return nil, err
	}

	_, err = ecs.NewService(ctx, metabaseResourceName, &ecs.ServiceArgs{
		Cluster:                         metabaseCluster.Arn,
		TaskDefinition:                  metabaseTaskDefinition.Arn,
		DesiredCount:                    pulumi.Int(1),
		DeploymentMaximumPercent:        pulumi.IntPtr(100),
		DeploymentMinimumHealthyPercent: pulumi.IntPtr(0),
		LaunchType:                      pulumi.String("FARGATE"),
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			AssignPublicIp: pulumi.BoolPtr(false),
			Subnets:        args.ECSSubnetIDs,
			SecurityGroups: pulumi.ToStringArrayOutput([]pulumi.StringOutput{metabaseSecurityGroup.ID().ToStringOutput()}),
		},
		LoadBalancers: ecs.ServiceLoadBalancerArray{
			ecs.ServiceLoadBalancerArgs{
				ContainerName:  pulumi.String("metabase"),
				ContainerPort:  pulumi.Int(metabasePort),
				TargetGroupArn: targetGroup.Arn,
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}

	metabaseDnsRecord, err := route53.NewRecord(ctx, metabaseResourceName, &route53.RecordArgs{
		ZoneId: hostedZoneID,
		Type:   pulumi.String("A"),
		Name:   args.DomainName,
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

	component.SecurityGroupID = metabaseSecurityGroup.ID().ToStringOutput()
	component.DNSName = pulumi.Sprintf("https://%s", metabaseDnsRecord.Name)

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"securityGroupId": metabaseSecurityGroup.ID().ToStringOutput(),
		"dnsName":         pulumi.Sprintf("https://%s", metabaseDnsRecord.Name),
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
		cluster.Port, cluster.DatabaseName, regionName, emailConfig.Host,
		emailConfig.Username, emailConfig.Password, emailConfig.Port, emailConfig.Security,
		metabaseImageName,
	).ApplyT(func(values []interface{}) (string, error) {
		hostname := values[0].(string)
		username := values[1].(string)
		password := values[2].(*string)
		port := values[3].(int)
		dbName := values[4].(string)
		region := values[5].(string)
		emailHost := values[6].(string)
		emailUsername := values[7].(string)
		emailPassword := values[8].(string)
		emailPort := values[9].(int)
		emailSecurity := values[10].(string)
		imageName := values[11].(string)

		metabaseEnv := []metabaseEnvironmentVariable{
			newMetabaseEnvironmentVariable("JAVA_TIMEZONE", "US/Pacific"),
			newMetabaseEnvironmentVariable("MB_DB_TYPE", "mysql"),
			newMetabaseEnvironmentVariable("MB_DB_DBNAME", dbName),
			newMetabaseEnvironmentVariable("MB_DB_PORT", fmt.Sprintf("%d", port)),
			newMetabaseEnvironmentVariable("MB_DB_USER", username),
			newMetabaseEnvironmentVariable("MB_DB_PASS", *password),
			newMetabaseEnvironmentVariable("MB_DB_HOST", hostname),
		}

		if emailHost != "" {
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
