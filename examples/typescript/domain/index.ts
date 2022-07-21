import * as pulumi from "@pulumi/pulumi";
import * as metabase from "@pulumi/metabase";

const metabaseService = new metabase.Metabase("metabaseService", {
    vpcId: "vpc-123",
    domain: {
        hostedZoneName: "example.com",
        domainName: "metabase.example.com",
    },
});
export const url = metabaseService.dnsName;
