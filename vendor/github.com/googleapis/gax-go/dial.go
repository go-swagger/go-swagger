package gax

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func DialGRPC(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	settings := &ClientSettings{}
	clientOptions(opts).Resolve(settings)

	if settings.Connection != nil {
		return settings.Connection, nil
	}
	var dialOpts = settings.DialOptions
	if len(dialOpts) == 0 {
		tokenSource, err := google.DefaultTokenSource(ctx, settings.Scopes...)
		if err != nil {
			return nil, fmt.Errorf("google.DefaultTokenSource: %v", err)
		}
		dialOpts = []grpc.DialOption{
			grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: tokenSource}),
			grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		}
	}
	return grpc.Dial(settings.Endpoint, dialOpts...)
}
