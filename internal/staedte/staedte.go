package staedte

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gittea.kittel.dev/marco/go-plz/internal/misc"
)

type Stadtfinder interface {
	SearchByPostalCodePrefix(prefix string) []string
	SearchByCityNamePrefix(prefix string) []string
	FindCitiesByName(name string) []Stadt
	FindCityByPostalCode(plz string) Stadt
	GetAirDistance(stadt, stadt2 Stadt) float64
}

type Stadt struct {
	Loc_id string
	Plz    string
	Lon    string
	Lat    string
	Ort    string
	Next   *Stadt
	Head   *Stadt
}

type Staedte struct {
	stadtnamenPrefixe     *misc.Trie
	postleitzahlenPrefixe *misc.Trie
}

func NewStaedte(pfadZurPostleitzahlenTxt string) *Staedte {
	snp, err := erstelleStadtTrie(pfadZurPostleitzahlenTxt)
	if err != nil {
		log.Fatalf("Error creating Stadt Trie: %v", err)
	}
	plzp, err := erstellePlzTrie(pfadZurPostleitzahlenTxt)
	if err != nil {
		log.Fatalf("Error creating Plz Trie: %v", err)
	}
	s := Staedte{stadtnamenPrefixe: snp, postleitzahlenPrefixe: plzp}
	return &s
}

func (s *Staedte) GetAirDistance(stadt, stadt2 Stadt) float64 {
	hs := misc.NewHaversine()
	lat1, err := strconv.ParseFloat(stadt.Lat, 64)
	if err != nil {
		return -1.0
	}
	lat2, err := strconv.ParseFloat(stadt2.Lat, 64)
	if err != nil {
		return -1.0
	}
	lon1, err := strconv.ParseFloat(stadt.Lon, 64)
	if err != nil {
		return -1.0
	}
	lon2, err := strconv.ParseFloat(stadt2.Lon, 64)
	if err != nil {
		return -1.0
	}
	return hs.CalcDistanceFromDegree(lat1, lon1, lat2, lon2)
}

func (s *Staedte) SearchByPostalCodePrefix(prefix string) []string {
	for {
		wort := prefix
		x := s.postleitzahlenPrefixe.Wildcard(fmt.Sprintf("%s.....", wort))
		i := 0
		for {
			if i >= len(x) {
				break
			}
			if len(x[i]) >= len(wort) {
				break
			}
			i++
		}
		return x[i:len(x)]
	}
}

func (s *Staedte) SearchByCityNamePrefix(prefix string) []string {
	for {
		wort := prefix
		x := s.stadtnamenPrefixe.Wildcard(fmt.Sprintf("%s..................", wort))
		i := 0
		for {
			if i >= len(x) {
				break
			}
			if len(x[i]) >= len(wort) {
				break
			}
			i++
		}
		return x[i:len(x)]
	}
}
func (s *Staedte) FindCitiesByName(name string) []Stadt {
	var result []Stadt
	node := s.stadtnamenPrefixe.Get(name)
	switch node.(type) {
	case Stadt:
		{
			node := node.(Stadt).Head
			for node != nil {
				result = append(result, *node)
				node = node.Next
			}
		}
	}
	return result
}
func (s *Staedte) FindCityByPostalCode(plz string) Stadt {
	var result = Stadt{}
	node := s.postleitzahlenPrefixe.Get(plz)
	switch node.(type) {
	case Stadt:
		{
			node := node.(Stadt)
			return node
		}
	}
	return result
}

func erstelleStadtTrie(postleitzahlen string) (*misc.Trie, error) {
	var stadtnamenPrefixer misc.Trie
	data, err := os.Open(postleitzahlen)
	if err != nil {
		return nil, fmt.Errorf("Plz Liste fehlt.")
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		lc := scanner.Text()
		lcarray := strings.Split(lc, "\t")
		if len(lcarray) != 5 {
			return nil, fmt.Errorf("Fehler: Anzahl der Elemente pro Zeile muss genau 4 betragen.")
		}
		// stadtSammlung = append(stadtSammlung, Stadt{Loc_id: lcarray[0], Plz: lcarray[1], Lon: lcarray[2], Lat: lcarray[3], Ort: lcarray[4]})
		dornodata := stadtnamenPrefixer.Get(lcarray[4])
		if dornodata == nil {
			stadtUndErsterKnoten := Stadt{Loc_id: lcarray[0], Plz: lcarray[1], Lon: lcarray[2], Lat: lcarray[3], Ort: lcarray[4]}
			stadtUndErsterKnoten.Head = &stadtUndErsterKnoten
			stadtUndErsterKnoten.Next = nil
			stadtnamenPrefixer.Put(lcarray[4], stadtUndErsterKnoten)
		} else {
			node := dornodata.(Stadt).Head
			for node.Next != nil {
				node = node.Next
			}
			node.Next = &Stadt{Loc_id: lcarray[0], Plz: lcarray[1], Lon: lcarray[2], Lat: lcarray[3], Ort: lcarray[4]}
			node = node.Next
			node.Next = nil
		}
	}
	return &stadtnamenPrefixer, nil
}
func erstellePlzTrie(postleitzahlen string) (*misc.Trie, error) {
	var plzPrefixer misc.Trie
	data, err := os.Open(postleitzahlen)
	if err != nil {
		return nil, fmt.Errorf("Plz Liste fehlt.")
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		lc := scanner.Text()
		lcarray := strings.Split(lc, "\t")
		if len(lcarray) != 5 {
			return nil, fmt.Errorf("Fehler: Anzahl der Elemente pro Zeile muss genau 4 betragen.")
		}
		// stadtSammlung = append(stadtSammlung, Stadt{Loc_id: lcarray[0], Plz: lcarray[1], Lon: lcarray[2], Lat: lcarray[3], Ort: lcarray[4]})
		dornodata := plzPrefixer.Get(lcarray[1])
		if dornodata == nil {
			stadtUndErsterKnoten := Stadt{Loc_id: lcarray[0], Plz: lcarray[1], Lon: lcarray[2], Lat: lcarray[3], Ort: lcarray[4]}
			stadtUndErsterKnoten.Head = &stadtUndErsterKnoten
			stadtUndErsterKnoten.Next = nil
			plzPrefixer.Put(lcarray[1], stadtUndErsterKnoten)
		} else {
			node := dornodata.(Stadt).Head
			for node.Next != nil {
				node = node.Next
			}
			node.Next = &Stadt{Loc_id: lcarray[0], Plz: lcarray[1], Lon: lcarray[2], Lat: lcarray[3], Ort: lcarray[4]}
			node = node.Next
			node.Next = nil
		}
	}
	return &plzPrefixer, nil
}
