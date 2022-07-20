import * as metabase from "@pulumi/metabase";

const instance = new metabase.Metabase("demo", {});

export const url = instance.dnsName;
