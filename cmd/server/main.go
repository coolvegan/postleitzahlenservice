package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	grpcserver "gittea.kittel.dev/marco/go-plz/internal/grpcmisc"
	pb "gittea.kittel.dev/marco/go-plz/internal/proto/pb"
	"gittea.kittel.dev/marco/go-plz/internal/staedte"
	"google.golang.org/grpc"
)

var (
	port = ":50051"
)

func main() {
	arbeitsverzeichnis, _ := os.Getwd()

	//Hier wird die Verbindungssicherheit festgestellt.
	//Umgebungsvariablen für die Modi Unverschlüsselt, TLS und mTLS werden gesetzt.
	//Dazu wird hier auch die Nachrichtensicherheit eingestellt. Aktuell nur Basic Auth.
	srv := grpc.NewServer(grpcserver.InitSecurity(arbeitsverzeichnis)...)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starte auf Port: %s \n", port)
	postleitzahlen := filepath.Join(arbeitsverzeichnis, "internal", "staedte", "postleitzahlen.txt")
	pb.RegisterStadteInformationenServer(srv, &grpcserver.Server{Stadtfinder: staedte.NewStaedte(postleitzahlen)})
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Betriebsfehler: v", err)
	}
}
