package grpcserver

import (
	"context"
	"strconv"

	pb "gittea.kittel.dev/marco/go-plz/internal/proto/pb"
	"gittea.kittel.dev/marco/go-plz/internal/staedte"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Stadtfinder staedte.Stadtfinder
	pb.UnimplementedStadteInformationenServer
}

func (s *Server) SucheNachAnteilEinerPostleitzahl(suche *pb.PostleitzahlenSuchbegriff, stream pb.StadteInformationen_SucheNachAnteilEinerPostleitzahlServer) error {
	if len(suche.GetPlzprefix()) < 2 {
		return status.Errorf(codes.Canceled, "Bitte mindestens zwei Ziffern des Postleitzahlenbereichs eingeben.")
	}
	stadtgebiete := s.Stadtfinder.SearchByPostalCodePrefix(suche.GetPlzprefix())
	for _, stadtteil := range stadtgebiete {
		stream.Send(&pb.Postleitzahl{Postleitzahl: stadtteil})
	}
	return nil
}
func (s *Server) SucheNachAnteilEinesStadtnamen(suche *pb.StadtSuchbegriff, stream pb.StadteInformationen_SucheNachAnteilEinesStadtnamenServer) error {
	if len(suche.GetStadtname()) < 2 {
		return status.Errorf(codes.Canceled, "Bitte mindestens zwei Buchstaben des Stadtnamen eingeben.")
	}
	stadtgebiete := s.Stadtfinder.SearchByCityNamePrefix(suche.GetStadtname())
	for _, stadtname := range stadtgebiete {
		stream.Send(&pb.Stadtname{Stadtname: stadtname})
	}
	return nil
}
func (s *Server) SucheNachExakenStadtnamen(suche *pb.StadtSuchbegriff, stream pb.StadteInformationen_SucheNachExakenStadtnamenServer) error {
	stadtgebiete := s.Stadtfinder.FindCitiesByName(suche.GetStadtname())
	if len(stadtgebiete) == 0 {
		return status.Errorf(codes.Canceled, "Es wurde keine Stadt unter dem Suchbegriff %s gefunden.", suche.GetStadtname())
	}
	for _, stadtteil := range stadtgebiete {
		//Es sind keine Fehler zu erwarten.
		lat, err1 := strconv.ParseFloat(stadtteil.Lat, 64)
		lon, err2 := strconv.ParseFloat(stadtteil.Lon, 64)
		si := pb.StadtInfo{LocationId: stadtteil.Loc_id, Langengrad: lon, Breitengrad: lat, Postleitzahl: stadtteil.Plz, Stadtname: stadtteil.Ort}
		if err1 == nil && err2 == nil && stadtteil.Plz != "" {
			stream.Send(&si)
		} else {
			return status.Errorf(codes.DataLoss, "unbekannter Fehler.")
		}
	}
	return nil
}
func (s *Server) SucheNachExtakerPostleitzahl(ctx context.Context, suche *pb.Postleitzahl) (*pb.StadtInfo, error) {
	stadt := s.Stadtfinder.FindCityByPostalCode(suche.GetPostleitzahl())
	if stadt.Plz == "" {
		return nil, status.Errorf(codes.Canceled, "Es wurde keine Stadt unter der Postleitzahl %s gefunden.", suche.GetPostleitzahl())
	}
	//Es sind keine Fehler zu erwarten.
	lat, err1 := strconv.ParseFloat(stadt.Lat, 64)
	lon, err2 := strconv.ParseFloat(stadt.Lon, 64)
	si := pb.StadtInfo{LocationId: stadt.Loc_id, Langengrad: lon, Breitengrad: lat, Postleitzahl: stadt.Plz, Stadtname: stadt.Ort}
	if err1 == nil && err2 == nil && stadt.Plz != "" {
		return &si, nil
	}
	return nil, status.Errorf(codes.NotFound, "Stadt wurde nicht gefunden.")
}
