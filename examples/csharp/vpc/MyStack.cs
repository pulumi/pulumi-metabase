using Pulumi;
using Metabase = Pulumi.Metabase;

class MyStack : Stack
{
    public MyStack()
    {
        var metabaseService = new Metabase.Metabase("metabaseService", new Metabase.MetabaseArgs
        {
            VpcId = "vpc-123",
        });
        this.Url = metabaseService.DnsName;
    }

    [Output("url")]
    public Output<string> Url { get; set; }
}
