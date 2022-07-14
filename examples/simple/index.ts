import * as pulumi from "@pulumi/pulumi";
import * as metabase from "@pulumi/metabase";

const zackDevStack = new pulumi.StackReference("pulumi/pulumi-service/zachary");

// The pulumi-service production VPC.
const vpcId = zackDevStack.getOutput("vpcId");
// The database subnets. These subnets are private and have no internet connectivity.
const dbSubnets = zackDevStack.getOutput("dbSubnetIds");
// The subnets that house our ECS clusters. These subnets are private and have outbound internet connectivity.
const ecsSubnets = zackDevStack.getOutput("ecsSubnetIds");
// These subnets are public. These subnets contain infrastructure that needs to be publicy exposed.
const publicSubnets = zackDevStack.getOutput("publicSubnetIds");

const instance = new metabase.Metabase("zack-demo", {
    vpcId,
    dbSubnetIds: dbSubnets,
    ecsSubnetIds: ecsSubnets,
    lbSubnetIds: publicSubnets,
});

export const url = instance.dnsName;
