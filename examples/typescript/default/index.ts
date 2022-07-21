import * as pulumi from "@pulumi/pulumi";
import * as metabase from "@pulumi/metabase";

const metabaseService = new metabase.Metabase("metabaseService", {});
export const url = metabaseService.dnsName;
