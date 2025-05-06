package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type jwtservice struct {
	geheimnis []byte
}

var ErrUngultigesToken = errors.New("Token ungültig")

type (
	AuthInterceptor struct {
		validateToken func(ctx context.Context, token string) (string, error)
	}
)

func NewAuthInterceptor(validator func(ctx context.Context, token string) (string, error)) (*AuthInterceptor, error) {
	if validator == nil {
		return nil, errors.New("Der Validator darf nicht null sein.")
	}
	return &AuthInterceptor{validateToken: validator}, nil
}

func (i *AuthInterceptor) UnaryJWTAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Metadaten fehlen")
	}

	token := md["authorization"]

	userID, err := i.validateToken(ctx, token[0])

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	ctx = context.WithValue(ctx, "user_id", userID)
	return handler(ctx, req)
}

// Der gewrappte Stream besitzt einen Context. In diesem wird der Token gesetzt.
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Der WrappedServerStream hilft beim Senden des JWT Token
func (w *wrappedServerStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}

func (i *AuthInterceptor) JWTAuthStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Skip authentication for public methods
	if !info.IsClientStream && strings.HasPrefix(info.FullMethod, "/api.public.") {
		return handler(srv, ss)
	}
	// Require authentication for everything else
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "method %s requires jwt token authentication", info.FullMethod)
	}

	token := md.Get("authorization")

	if len(token) == 0 {
		return status.Errorf(codes.Unauthenticated, "method %s error in decoding", info.FullMethod)
	}
	userID, err := i.validateToken(ss.Context(), token[0])

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Der Token ist ungültig: %v", token)
	}

	ctx := context.WithValue(ss.Context(), "user_id", userID)

	wrappedStream := &wrappedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	}
	return handler(srv, wrappedStream)
}

func NewService(geheimnis string) (*jwtservice, error) {
	if geheimnis == "" {
		return nil, errors.New("Das Geheimnis darf nicht leer sein.")
	}
	return &jwtservice{geheimnis: []byte(geheimnis)}, nil
}

// Erzeugt einen Token mit einer Benutzerbeziehung über eine Benutzerkennung
func (s *jwtservice) IssueToken(_ context.Context, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iss": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}, nil)

	signed, err := token.SignedString(s.geheimnis)
	if err != nil {
		return "", fmt.Errorf("Das Signiern mit JWT schlug fehl.: %w", err)
	}
	return signed, nil
}

func (s *jwtservice) ValidateToken(_ context.Context, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unerwartede Signiermethode: %v", token.Header["alg"])
		}
		return s.geheimnis, nil
	})

	if err != nil {
		return "", errors.Join(ErrUngultigesToken)
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		id, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("%w: failed to extract id from claims", ErrUngultigesToken)
		}
		return id, nil
	}

	return "", ErrUngultigesToken
}

type JWTAuth struct {
	Token string
}

func (j JWTAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.Token,
	}, nil
}
func (j JWTAuth) RequireTransportSecurity() bool {
	return true
}
