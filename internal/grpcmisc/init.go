package grpcserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
	"gittea.kittel.dev/marco/go-plz/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	//Für TLS Wichtig
	certDirectory = "certs/myca_certs"
	crtFile       = "server.crt"
	keyFile       = "server.key"
	//Für mTLS wird noch zusätzlich eine CA gebraucht.
	caFile = "myca.crt"
	//Unverschlüsselt  | TLS | mTLS
)

func InitSecurity(arbeitsverzeichnis string) []grpc.ServerOption {
	wd := arbeitsverzeichnis
	srvType := os.Getenv("RPCTYPE")
	var authOptionen []grpc.ServerOption
	if os.Getenv("AUTH") != "" && os.Getenv("AUTH") == "BASIC" {
		if os.Getenv("USERNAME") == "" || os.Getenv("PASSWORD") == "" && (os.Getenv("RPCTYPE") != "TLS" || os.Getenv("RPCTYPE") != "mTLS") {
			log.Fatalf("BasicAuth erfordert die Umgebungsvariablen USERNAME und PASSWORD.")
		}
		log.Printf("Security BasicAuth.")
		authOptionen = append(authOptionen, grpc.UnaryInterceptor(auth.BasicAuthUnaryInterceptor), grpc.StreamInterceptor(auth.BasicAuthInterceptor))
	}
	if os.Getenv("AUTH") != "" && os.Getenv("AUTH") == "JWT" {
		if os.Getenv("JWTSECRET") == "" {
			log.Fatal("Environment Variable JWTSECRET fehlt.\n")
		}

		validaterFunc := func(ctx context.Context, token string) (string, error) {
			jwttoken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWTSECRET")), nil // 🔑 _Derselbe Secret wie beim Generieren_
			})

			if err != nil {
				return "", fmt.Errorf("Token ungültig.")
			}

			// Gültig → Claims extrahieren
			if claims, ok := jwttoken.Claims.(jwt.MapClaims); !ok {
				fmt.Println("Name:", claims["name"]) // "Dev Tester"
				fmt.Println("Exp:", claims["exp"])   // 1746272830 (Unix timestamp)
			}

			return "", nil
		}

		var authI, err = auth.NewAuthInterceptor(validaterFunc)
		if err != nil {
			log.Fatalln("Es gibt Probleme mit der Funktionsdefintion des Validators.")
		}

		log.Printf("Security JWTAuth.")
		authOptionen = append(authOptionen, grpc.UnaryInterceptor(authI.UnaryJWTAuthInterceptor), grpc.StreamInterceptor(authI.JWTAuthStreamInterceptor))
	}

	if os.Getenv("RPCCRTDIRECTORY") != "" {
		certDirectory = os.Getenv("RPCCRTDIRECTORY")
	}

	if os.Getenv("RPCCAFILE") == "" {
		caFile = filepath.Join(wd, certDirectory, caFile)
	} else {
		caFile = os.Getenv("RPCCAFILE")
	}

	if os.Getenv("RPCCRTFILE") == "" {
		crtFile = filepath.Join(wd, certDirectory, crtFile)
	} else {
		crtFile = os.Getenv("RPCCRTFILE")
	}

	if os.Getenv("RPCCRTKEY") == "" {
		keyFile = filepath.Join(wd, certDirectory, keyFile)
	} else {
		keyFile = os.Getenv("RPCCRTKEY")
	}

	var opts []grpc.ServerOption
	switch srvType {
	case "TLS":
		cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Fatalf("Fehler beim Laden des Schlüsselpaars: %s", err)
		}
		tlsOptionen := []grpc.ServerOption{
			grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		}
		opts = append(opts, tlsOptionen...)
		opts = append(opts, authOptionen...)

		log.Println("TLS Modus ausgewählt.")
	case "mTLS":
		cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Fatalf("Fehler beim Laden des Schlüsselpaars: %s", err)
		}
		certPool := x509.NewCertPool()
		ca, err := os.ReadFile(caFile)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)

		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append ca certificate")
		}

		mtlsOptionen := []grpc.ServerOption{
			// Enable TLS for all incoming connections.

			grpc.Creds(
				credentials.NewTLS(&tls.Config{
					ClientAuth:   tls.RequireAndVerifyClientCert,
					Certificates: []tls.Certificate{cert},
					ClientCAs:    certPool},
				)),
		}
		opts = append(opts, mtlsOptionen...)
		opts = append(opts, authOptionen...)
		log.Println("mTLS Modus ausgewählt.")
	default:
		opts := []grpc.ServerOption{}
		opts = append(opts, authOptionen...)

		//Unverschlüsselt

		log.Println("Unverschlüsselter Modus ausgewählt.")
		log.Println("Achtung: Authentifizierung werden ignoriert.")

	}
	return opts
}
