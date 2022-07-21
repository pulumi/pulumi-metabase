import pulumi
import pulumi_metabase as metabase

metabase_service = metabase.Metabase("metabaseService")
pulumi.export("url", metabase_service.dns_name)
