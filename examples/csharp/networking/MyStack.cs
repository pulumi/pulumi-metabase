using Pulumi;
using Metabase = Pulumi.Metabase;

class MyStack : Stack
{
    public MyStack()
    {
        var metabaseService = new Metabase.Metabase("metabaseService", new Metabase.MetabaseArgs
        {
            VpcId = "vpc-123",
            Networking = new Metabase.Inputs.NetworkingArgs
            {
                EcsSubnetIds = 
                {
                    "subnet-123",
                    "subnet-456",
                },
                DbSubnetIds = 
                {
                    "subnet-789",
                    "subnet-abc",
                },
                LbSubnetIds = 
                {
                    "subnet-def",
                    "subnet-ghi",
                },
            },
        });
        this.Url = metabaseService.DnsName;
    }

    [Output("url")]
    public Output<string> Url { get; set; }
}
