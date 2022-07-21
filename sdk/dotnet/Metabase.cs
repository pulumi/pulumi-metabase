// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Metabase
{
    /// <summary>
    /// This resources provisions a container running Metabase on AWS ECS Fargate. By default
    /// the resource will run the service in the AWS Account's Default VPC unless a VPC is defined. This
    /// resource will also deploy the `latest` version of Metabase unless a version is supplied.
    /// 
    /// You can provide specific subnets to host the Load Balancer, Database, and ECS Service, as well
    /// as provide a custom domain name for the service.
    /// 
    /// ## Example Usage
    /// ### Default
    /// 
    /// ```csharp
    /// using Pulumi;
    /// using Metabase = Pulumi.Metabase;
    /// 
    /// class MyStack : Stack
    /// {
    ///     public MyStack()
    ///     {
    ///         var metabaseService = new Metabase.Metabase("metabaseService", new Metabase.MetabaseArgs
    ///         {
    ///         });
    ///         this.Url = metabaseService.DnsName;
    ///     }
    /// 
    ///     [Output("url")]
    ///     public Output&lt;string&gt; Url { get; set; }
    /// }
    /// ```
    /// {{ /example }}
    /// ### Custom Domain &amp; Networking
    /// 
    /// ```csharp
    /// using Pulumi;
    /// using Metabase = Pulumi.Metabase;
    /// 
    /// class MyStack : Stack
    /// {
    ///     public MyStack()
    ///     {
    ///         var metabaseService = new Metabase.Metabase("metabaseService", new Metabase.MetabaseArgs
    ///         {
    ///             VpcId = "vpc-123",
    ///             Networking = new Metabase.Inputs.NetworkingArgs
    ///             {
    ///                 EcsSubnetIds =
    ///                 {
    ///                     "subnet-123",
    ///                     "subnet-456",
    ///                 },
    ///                 DbSubnetIds =
    ///                 {
    ///                     "subnet-789",
    ///                     "subnet-abc",
    ///                 },
    ///                 LbSubnetIds =
    ///                 {
    ///                     "subnet-def",
    ///                     "subnet-ghi",
    ///                 },
    ///             },
    ///             Domain = new Metabase.Inputs.CustomDomainArgs
    ///             {
    ///                 HostedZoneName = "example.com",
    ///                 DomainName = "metabase.example.com",
    ///             },
    ///         });
    ///         this.Url = metabaseService.DnsName;
    ///     }
    /// 
    ///     [Output("url")]
    ///     public Output&lt;string&gt; Url { get; set; }
    /// }
    /// ```
    /// {{ /example }}
    /// </summary>
    [MetabaseResourceType("metabase:index:Metabase")]
    public partial class Metabase : Pulumi.ComponentResource
    {
        /// <summary>
        /// The DNS name for the Metabase instance.
        /// </summary>
        [Output("dnsName")]
        public Output<string> DnsName { get; private set; } = null!;

        /// <summary>
        /// The security group id for the Metabase instance.
        /// </summary>
        [Output("securityGroupId")]
        public Output<string> SecurityGroupId { get; private set; } = null!;


        /// <summary>
        /// Create a Metabase resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Metabase(string name, MetabaseArgs? args = null, ComponentResourceOptions? options = null)
            : base("metabase:index:Metabase", name, args ?? new MetabaseArgs(), MakeResourceOptions(options, ""), remote: true)
        {
        }

        private static ComponentResourceOptions MakeResourceOptions(ComponentResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new ComponentResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = ComponentResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
    }

    public sealed class MetabaseArgs : Pulumi.ResourceArgs
    {
        /// <summary>
        /// Optionally provide a hosted zone and domain name for the Metabase service.
        /// </summary>
        [Input("domain")]
        public Input<Inputs.CustomDomainArgs>? Domain { get; set; }

        /// <summary>
        /// The version of Metabase to run - used as a tag on the `metabase/metabase` Dockerhub image.
        /// </summary>
        [Input("metabaseVersion")]
        public Input<string>? MetabaseVersion { get; set; }

        /// <summary>
        /// Optionally provide specific subnet IDs to run the different resources of Metabase.
        /// </summary>
        [Input("networking")]
        public Input<Inputs.NetworkingArgs>? Networking { get; set; }

        /// <summary>
        /// The VPC to use for the Metabase service. If left blank then the default VPC will be used.
        /// </summary>
        [Input("vpcId")]
        public Input<string>? VpcId { get; set; }

        public MetabaseArgs()
        {
        }
    }
}
