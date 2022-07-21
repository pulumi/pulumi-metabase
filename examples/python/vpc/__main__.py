import pulumi
import pulumi_metabase as metabase

metabase_service = metabase.Metabase("metabaseService", vpc_id="vpc-123")
pulumi.export("url", metabase_service.dns_name)
