package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gittea.kittel.dev/marco/go-plz/internal/auth"
	stadtservice "gittea.kittel.dev/marco/go-plz/internal/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address = "localhost:50051"
	//Für TLS Wichtig
	certDirectory = "certs/myca_certs"
	srvCrtFile    = "server.crt"

	//Für mTLS werden noch zusätzlich eine client.key, client.crt und eine ca gebraucht.
	keyFile = "client.key"
	crtFile = "client.crt"
	caFile  = "myca.crt"

	//Unverschlüsselt  | TLS | mTLS
	clientType = ""
	hostname   = "localhost"

	//Auth Optionen
	authOptionen []grpc.DialOption
)

func main() {
	wd, _ := os.Getwd()
	clientType = os.Getenv("RPCTYPE")

	if os.Getenv("RPCCRTDIRECTORY") != "" {
		certDirectory = os.Getenv("RPCCRTDIRECTORY")
	}

	if os.Args[1] != "" {
		address = os.Args[1]
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

	if os.Getenv("RPCSRVCRTFILE") == "" {
		srvCrtFile = filepath.Join(wd, certDirectory, srvCrtFile)
	} else {
		crtFile = os.Getenv("RPCSRVCRTFILE")
	}

	if os.Getenv("RPCCRTKEY") == "" {
		keyFile = filepath.Join(wd, certDirectory, keyFile)
	} else {
		keyFile = os.Getenv("RPCCRTKEY")
	}

	if os.Getenv("BASICAUTH") != "" {
		data := strings.Split(os.Getenv("BASICAUTH"), ":")
		if len(data) != 2 {
			log.Fatalln("Error: BASICAUTH=<username>:<password> lautet der Syntax.")
		}

		auth := auth.BasicAuth{
			Username: data[0],
			Password: data[1],
		}
		authOptionen = append(authOptionen, grpc.WithPerRPCCredentials(auth))
	}

	if os.Getenv("JWT") != "" {
		auth := auth.JWTAuth{Token: os.Getenv("JWT")}
		authOptionen = append(authOptionen, grpc.WithPerRPCCredentials(auth))

	}

	var conn *grpc.ClientConn
	var err error

	switch clientType {
	case "TLS":

		creds, err := credentials.NewClientTLSFromFile(srvCrtFile, hostname)
		if err != nil {
			log.Fatalf("Fehler beim Laden des Schlüsselpaars: %s", err)
		}
		opts := []grpc.DialOption{

			grpc.WithTransportCredentials(creds),
		}

		opts = append(opts, authOptionen...)

		conn, err = grpc.NewClient(address, opts...)
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

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				ServerName:   hostname, // NOTE: this is required!
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
			})),
		}
		//Hier setzen wir unsere Auth Optionen rein, falls es diese gibt.
		opts = append(opts, authOptionen...)
		conn, err = grpc.NewClient(address, opts...)
		log.Println("mTLS Modus ausgewählt.")

	default:
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials())}
		//Unverschlüsselt
		conn, err = grpc.NewClient(address, opts...)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		log.Println("Unverschlüsselter Modus ausgewählt.")
		log.Println("Achtung: Authentifizierung werden ignoriert.")
	}

	defer conn.Close()
	c := stadtservice.NewStadteInformationenClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	//a, err := c.SucheNachExtakerPostleitzahl(ctx, &stadtservice.Postleitzahl{Postleitzahl: "59192"})
	s, err := c.SucheNachExakenStadtnamen(ctx, &stadtservice.StadtSuchbegriff{Stadtname: os.Args[2]})
	if err != nil {
		log.Fatal(err)
	}
	for {
		data, err := s.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(data)
	}
	fmt.Println("\n\n")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(a)

}

//https://www.bytesizego.com/blog/securing-grpc-golang
