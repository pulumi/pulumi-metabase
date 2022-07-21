using Pulumi;
using Metabase = Pulumi.Metabase;

class MyStack : Stack
{
    public MyStack()
    {
        var metabaseService = new Metabase.Metabase("metabaseService", new Metabase.MetabaseArgs
        {
        });
        this.Url = metabaseService.DnsName;
    }

    [Output("url")]
    public Output<string> Url { get; set; }
}
