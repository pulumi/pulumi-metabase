import pulumi
import pulumi_metabase as metabase

metabase_service = metabase.Metabase("metabaseService",
    vpc_id="vpc-123",
    domain=metabase.CustomDomainArgs(
        hosted_zone_name="example.com",
        domain_name="metabase.example.com",
    ))
pulumi.export("url", metabase_service.dns_name)
