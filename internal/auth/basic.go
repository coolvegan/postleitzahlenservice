package auth

import (
	"context"
	"encoding/base64"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type BasicAuth struct {
	Username string
	Password string
}

func BasicAuthInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Skip authentication for public methods
	if !info.IsClientStream && strings.HasPrefix(info.FullMethod, "/api.public.") {
		return handler(srv, ss)
	}
	// Require authentication for everything else
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "method %s requires basic authentication", info.FullMethod)
	}

	authEncodec := md.Get("authorization")

	if len(authEncodec) == 0 {
		return status.Errorf(codes.Unauthenticated, "method %s error in decoding", info.FullMethod)
	}

	token := strings.TrimPrefix(authEncodec[0], "Basic ")
	uDec, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "method %s error in decoding", info.FullMethod)
	}
	authDecoded := strings.Split(string(uDec), ":")
	if len(authDecoded) != 2 {
		return status.Errorf(codes.Unauthenticated, "method %s requires basic authentication", info.FullMethod)
	}

	username := authDecoded[0]
	password := authDecoded[1]

	if username != os.Getenv("USERNAME") || password != os.Getenv("PASSWORD") {
		return status.Errorf(codes.Unauthenticated, "method %s unknown user or wrong password", info.FullMethod)
	}

	// Proceed with the RPC
	return handler(srv, ss)
}

func BasicAuthUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Require authentication for everything else
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "method %s requires basic authentication", info.FullMethod)
	}

	authEncodec := md.Get("authorization")

	if len(authEncodec) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "method %s error in decoding", info.FullMethod)
	}

	token := strings.TrimPrefix(authEncodec[0], "Basic ")
	uDec, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "method %s error in decoding", info.FullMethod)
	}
	authDecoded := strings.Split(string(uDec), ":")
	if len(authDecoded) != 2 {
		return nil, status.Errorf(codes.Unauthenticated, "method %s requires basic authentication", info.FullMethod)
	}

	username := authDecoded[0]
	password := authDecoded[1]

	if username != os.Getenv("USERNAME") || password != os.Getenv("PASSWORD") {
		return nil, status.Errorf(codes.Unauthenticated, "method %s unknown user or wrong password", info.FullMethod)
	}

	return handler(ctx, req)
}

//Das Interface muss für PerRPCCredentials überschrieben werden

func (b BasicAuth) GetRequestMetadata(ctx context.Context,
	in ...string) (map[string]string, error) {
	auth := b.Username + ":" + b.Password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}
func (b BasicAuth) RequireTransportSecurity() bool {
	return true
}
