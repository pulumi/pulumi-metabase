package metabase

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/route53"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type MetabaseResourceConstructor struct {
	ctx              *pulumi.Context
	name             string
	baseResourceName string
	opts             []pulumi.ResourceOption
}

func NewMetabaseResourceConstructor(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) *MetabaseResourceConstructor {
	return &MetabaseResourceConstructor{
		ctx:              ctx,
		name:             name,
		baseResourceName: fmt.Sprintf("%s-metabase", name),
		opts:             opts,
	}
}

func (m *MetabaseResourceConstructor) NewMetabaseSecurityGroup(vpcID pulumi.StringInput) (*ec2.SecurityGroup, error) {
	metabaseSecurityGroupName := fmt.Sprintf("%s-metabase-sg", m.name)
	return ec2.NewSecurityGroup(m.ctx, metabaseSecurityGroupName, &ec2.SecurityGroupArgs{
		VpcId: vpcID,
	}, m.opts...)
}

func (m *MetabaseResourceConstructor) NewLoadBalancerSecurityGroup(vpcID pulumi.StringInput, metabasePort int, metabaseSecurityGroupID pulumi.StringOutput) (*ec2.SecurityGroup, error) {
	loadBalancerSecurityGroupName := fmt.Sprintf("%s-metabase-lb-sg", m.name)
	return ec2.NewSecurityGroup(m.ctx, loadBalancerSecurityGroupName, &ec2.SecurityGroupArgs{
		VpcId: vpcID,
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
				SecurityGroups: pulumi.ToStringArrayOutput([]pulumi.StringOutput{metabaseSecurityGroupID}),
			},
			ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("tcp"),
				ToPort:     pulumi.Int(443),
				FromPort:   pulumi.Int(443),
				CidrBlocks: pulumi.ToStringArray([]string{"0.0.0.0/0"}),
			},
		},
	}, m.opts...)
}

func (m *MetabaseResourceConstructor) NewSecurityGroupRules(metabasePort int, metabaseSecurityGroupID, loadBalancerSecurityGroupID pulumi.IDOutput) error {
	metabaseSecurityGroupSegmentRuleName := fmt.Sprintf("%s-metabase-segment", m.name)
	_, err := ec2.NewSecurityGroupRule(m.ctx, metabaseSecurityGroupSegmentRuleName, &ec2.SecurityGroupRuleArgs{
		Description:           pulumi.String("Allow access to Metabase from the Load Balancer"),
		SecurityGroupId:       metabaseSecurityGroupID,
		Type:                  pulumi.String("ingress"),
		Protocol:              pulumi.String("tcp"),
		FromPort:              pulumi.Int(metabasePort),
		ToPort:                pulumi.Int(metabasePort),
		SourceSecurityGroupId: loadBalancerSecurityGroupID,
	}, m.opts...)
	if err != nil {
		return err
	}

	metabaseSecurityGroupSelfRuleName := fmt.Sprintf("%s-metabase-self", m.name)
	_, err = ec2.NewSecurityGroupRule(m.ctx, metabaseSecurityGroupSelfRuleName, &ec2.SecurityGroupRuleArgs{
		Description:           pulumi.String("Allow access to anything from within the Security Group"),
		SecurityGroupId:       metabaseSecurityGroupID,
		Type:                  pulumi.String("ingress"),
		Protocol:              pulumi.String("tcp"),
		FromPort:              pulumi.Int(0),
		ToPort:                pulumi.Int(65535),
		SourceSecurityGroupId: metabaseSecurityGroupID,
	}, m.opts...)
	if err != nil {
		return err
	}

	metabaseSecurityGroupEgressRuleName := fmt.Sprintf("%s-metabase-egress", m.name)
	_, err = ec2.NewSecurityGroupRule(m.ctx, metabaseSecurityGroupEgressRuleName, &ec2.SecurityGroupRuleArgs{
		Description:     pulumi.String("Allow egress to anywhere"),
		SecurityGroupId: metabaseSecurityGroupID,
		Type:            pulumi.String("egress"),
		Protocol:        pulumi.String("tcp"),
		FromPort:        pulumi.Int(0),
		ToPort:          pulumi.Int(65535),
		CidrBlocks:      pulumi.ToStringArray([]string{"0.0.0.0/0"}),
	}, m.opts...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetabaseResourceConstructor) NewMetabasePassword() (*random.RandomString, error) {
	return random.NewRandomString(m.ctx, m.baseResourceName, &random.RandomStringArgs{
		Special: pulumi.BoolPtr(false),
		Length:  pulumi.Int(20),
	}, m.opts...)
}

func (m *MetabaseResourceConstructor) NewMySQLCluster(dbSubnetIDs pulumi.StringArrayInput, metabasePassword *random.RandomString, metabaseSecurityGroupID pulumi.IDOutput) (*rds.Cluster, error) {
	metabaseMysqlSubnetGroup, err := rds.NewSubnetGroup(m.ctx, m.baseResourceName, &rds.SubnetGroupArgs{
		SubnetIds: dbSubnetIDs,
	}, m.opts...)
	if err != nil {
		return nil, err
	}

	return rds.NewCluster(m.ctx, m.baseResourceName, &rds.ClusterArgs{
		ClusterIdentifier:       pulumi.Sprintf("%smetabasemysql", m.name),
		DatabaseName:            pulumi.String("metabase"),
		MasterUsername:          pulumi.String("admin"),
		MasterPassword:          metabasePassword.Result,
		Engine:                  pulumi.String("aurora-mysql"),
		EngineMode:              pulumi.String("serverless"),
		EngineVersion:           pulumi.String("5.7.mysql_aurora.2.07.1"),
		DbSubnetGroupName:       metabaseMysqlSubnetGroup.Name,
		VpcSecurityGroupIds:     pulumi.ToStringArrayOutput([]pulumi.StringOutput{metabaseSecurityGroupID.ToStringOutput()}),
		FinalSnapshotIdentifier: pulumi.Sprintf("%smetabasefinalsnapshot", m.name),
		EnableHttpEndpoint:      pulumi.BoolPtr(true),
	}, m.opts...)
}

func (m *MetabaseResourceConstructor) GetHostedZoneId(hostedZoneName pulumi.StringInput) pulumi.StringOutput {
	return hostedZoneName.ToStringOutput().ApplyT(func(name string) (string, error) {
		hostedZone, err := route53.LookupZone(m.ctx, &route53.LookupZoneArgs{
			Name: &name,
		})
		if err != nil {
			return "", nil
		}
		return hostedZone.Id, nil
	}).(pulumi.StringOutput)
}

func (m *MetabaseResourceConstructor) NewDomainCertificate(domainName pulumi.StringInput, hostedZoneID pulumi.StringOutput) (*acm.Certificate, *acm.CertificateValidation, error) {
	certificate, err := acm.NewCertificate(m.ctx, m.baseResourceName, &acm.CertificateArgs{
		DomainName:       domainName,
		ValidationMethod: pulumi.String("DNS"),
	}, m.opts...)
	if err != nil {
		return nil, nil, err
	}

	certificateValidationRecordName := fmt.Sprintf("%s-metabase-certvalidation", m.name)
	certificateValidationRecord, err := route53.NewRecord(m.ctx, certificateValidationRecordName, &route53.RecordArgs{
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
	}, m.opts...)
	if err != nil {
		return nil, nil, err
	}

	certificateValidation, err := acm.NewCertificateValidation(m.ctx, m.baseResourceName, &acm.CertificateValidationArgs{
		CertificateArn:        certificate.Arn,
		ValidationRecordFqdns: pulumi.ToStringArrayOutput([]pulumi.StringOutput{certificateValidationRecord.Fqdn}),
	}, m.opts...)
	if err != nil {
		return nil, nil, err
	}

	return certificate, certificateValidation, nil
}

func (m *MetabaseResourceConstructor) NewLoadBalancer(
	vpcID pulumi.StringInput, lbSubnetIDs pulumi.StringArrayInput,
	loadBalancerSecurityGroupID pulumi.IDOutput, metabasePort int,
	certificateValidation *acm.CertificateValidation, certificate *acm.Certificate) (*lb.LoadBalancer, *lb.TargetGroup, *lb.Listener, error) {
	// Stable load balancer endpoint (no other way to get a consistent IP for an ECS service!!!)
	loadBalancer, err := lb.NewLoadBalancer(m.ctx, m.baseResourceName, &lb.LoadBalancerArgs{
		LoadBalancerType: pulumi.String("application"),
		Subnets:          lbSubnetIDs,
		SecurityGroups:   pulumi.ToStringArrayOutput([]pulumi.StringOutput{loadBalancerSecurityGroupID.ToStringOutput()}),
		IdleTimeout:      pulumi.IntPtr(600),
	}, m.opts...)
	if err != nil {
		return nil, nil, nil, err
	}

	targetGroup, err := lb.NewTargetGroup(m.ctx, m.baseResourceName, &lb.TargetGroupArgs{
		TargetType: pulumi.String("ip"),
		Port:       pulumi.Int(metabasePort),
		Protocol:   pulumi.String("HTTP"),
		VpcId:      vpcID,
		// Since this is a user facing tool, and we only have 0 or 1 running instances, we don't need to wait to
		// drain connections, and instead want to ensure we have as little downtime as possible.
		DeregistrationDelay: pulumi.Int(0),
	}, m.opts...)
	if err != nil {
		return nil, nil, nil, err
	}

	listenerArgs := &lb.ListenerArgs{
		LoadBalancerArn: loadBalancer.Arn,
		Port:            pulumi.Int(443),
		Protocol:        pulumi.String("HTTP"),
		DefaultActions: lb.ListenerDefaultActionArray{
			lb.ListenerDefaultActionArgs{
				Type:           pulumi.String("forward"),
				Order:          pulumi.IntPtr(1),
				TargetGroupArn: targetGroup.Arn,
			},
		},
	}

	if certificate != nil {
		listenerArgs.CertificateArn = certificate.Arn
		listenerArgs.SslPolicy = pulumi.String("ELBSecurityPolicy-TLS-1-2-2017-01")
	}

	listenerOpts := m.opts
	if certificateValidation != nil {
		listenerOpts = append(m.opts, pulumi.DependsOn([]pulumi.Resource{certificateValidation}))
	}

	listener, err := lb.NewListener(m.ctx, m.baseResourceName, listenerArgs, listenerOpts...)
	if err != nil {
		return nil, nil, nil, err
	}

	httpRedirectListenerName := fmt.Sprintf("%s-metabase-redirecthttp", m.name)
	_, err = lb.NewListener(m.ctx, httpRedirectListenerName, &lb.ListenerArgs{
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
	}, m.opts...)
	if err != nil {
		return nil, nil, nil, err
	}

	return loadBalancer, targetGroup, listener, nil
}

func (m *MetabaseResourceConstructor) NewMetabaseService(
	metabaseVersion pulumi.StringInput, regionName pulumi.StringOutput,
	metabaseContainerDef pulumi.StringOutput, ecsSubnetIDs pulumi.StringArrayInput,
	metabaseSecurityGroupID pulumi.IDOutput, metabasePort int, targetGroupARN pulumi.StringOutput,
	lbListener *lb.Listener,
) error {

	metabaseCluster, err := ecs.NewCluster(m.ctx, m.baseResourceName, &ecs.ClusterArgs{}, m.opts...)
	if err != nil {
		return err
	}

	metabaseExecutionRole, err := iam.GetRole(m.ctx, m.baseResourceName, pulumi.ID("ecsTaskExecutionRole"), nil, m.opts...)
	if err != nil {
		return err
	}

	metabaseTaskDefinition, err := ecs.NewTaskDefinition(m.ctx, m.baseResourceName, &ecs.TaskDefinitionArgs{
		Family:                  pulumi.String("metabase"),
		Cpu:                     pulumi.String("2048"),
		Memory:                  pulumi.String("8192"),
		RequiresCompatibilities: pulumi.ToStringArray([]string{"FARGATE"}),
		NetworkMode:             pulumi.StringPtr("awsvpc"),
		ExecutionRoleArn:        metabaseExecutionRole.Arn,
		ContainerDefinitions:    metabaseContainerDef,
	}, m.opts...)
	if err != nil {
		return err
	}

	serviceOpts := append(m.opts, pulumi.DependsOn([]pulumi.Resource{lbListener}))
	_, err = ecs.NewService(m.ctx, m.baseResourceName, &ecs.ServiceArgs{
		Cluster:                         metabaseCluster.Arn,
		TaskDefinition:                  metabaseTaskDefinition.Arn,
		DesiredCount:                    pulumi.Int(1),
		DeploymentMaximumPercent:        pulumi.IntPtr(100),
		DeploymentMinimumHealthyPercent: pulumi.IntPtr(0),
		LaunchType:                      pulumi.String("FARGATE"),
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			AssignPublicIp: pulumi.BoolPtr(false),
			Subnets:        ecsSubnetIDs,
			SecurityGroups: pulumi.ToStringArrayOutput([]pulumi.StringOutput{metabaseSecurityGroupID.ToStringOutput()}),
		},
		LoadBalancers: ecs.ServiceLoadBalancerArray{
			ecs.ServiceLoadBalancerArgs{
				ContainerName:  pulumi.String("metabase"),
				ContainerPort:  pulumi.Int(metabasePort),
				TargetGroupArn: targetGroupARN,
			},
		},
	}, serviceOpts...)
	return err
}
