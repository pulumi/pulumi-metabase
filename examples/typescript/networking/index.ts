import * as pulumi from "@pulumi/pulumi";
import * as metabase from "@pulumi/metabase";

const metabaseService = new metabase.Metabase("metabaseService", {
    vpcId: "vpc-123",
    networking: {
        ecsSubnetIds: [
            "subnet-123",
            "subnet-456",
        ],
        dbSubnetIds: [
            "subnet-789",
            "subnet-abc",
        ],
        lbSubnetIds: [
            "subnet-def",
            "subnet-ghi",
        ],
    },
});
export const url = metabaseService.dnsName;
