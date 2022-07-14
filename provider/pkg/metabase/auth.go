package metabase

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AuthenticationStrategy string

const (
	GoogleAuthentication AuthenticationStrategy = "google"
)

type AuthenticationProvider interface {
	DefaultActions(targetGroupARN pulumi.StringOutput) lb.ListenerDefaultActionArrayInput
}

type Authentication struct {
	Type               AuthenticationStrategy `pulumi:"authenticationStrategy"`
	AuthenticationArgs AuthenticationProvider `pulumi:"authenticationArgs"`
}

// Google Authentication
type GoogleOIDCAuthentication struct {
	OIDCClientId     pulumi.StringInput
	OIDCClientSecret pulumi.StringInput
}

func (g GoogleOIDCAuthentication) DefaultActions(targetGroupARN pulumi.StringOutput) lb.ListenerDefaultActionArrayInput {
	return lb.ListenerDefaultActionArray{
		lb.ListenerDefaultActionArgs{
			Type:  pulumi.String("authenticate-oidc"),
			Order: pulumi.IntPtr(1),
			AuthenticateOidc: &lb.ListenerDefaultActionAuthenticateOidcArgs{
				OnUnauthenticatedRequest: pulumi.String("authenticate"),
				Issuer:                   pulumi.String("https://accounts.google.com"),
				AuthorizationEndpoint:    pulumi.String("https://accounts.google.com/o/oauth2/v2/auth"),
				TokenEndpoint:            pulumi.String("https://oauth2.googleapis.com/token"),
				UserInfoEndpoint:         pulumi.String("https://openidconnect.googleapis.com/v1/userinfo"),
				ClientId:                 g.OIDCClientId,
				ClientSecret:             g.OIDCClientSecret,
			},
		},
		lb.ListenerDefaultActionArgs{
			Type:           pulumi.String("forward"),
			Order:          pulumi.IntPtr(2),
			TargetGroupArn: targetGroupARN,
		},
	}
}
