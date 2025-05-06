package staedte

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewStaedte(t *testing.T) {
	staedte := NewStaedte("postleitzahlen.txt")
	if staedte == nil {
		t.Errorf("NewStaedte() returned nil")
	}
}

func TestStaedte_PrefixPostleitzahl(t *testing.T) {
	staedte := NewStaedte("postleitzahlen.txt")
	tests := []struct {
		name      string
		plzprefix string
		want      []string
	}{
		{
			name:      "Test with prefix 01",
			plzprefix: "01",
			want:      []string{"01067", "01069", "01097", "01099", "01109", "01127", "01129", "01139", "01157", "01159", "01169", "01187", "01189", "01217", "01219", "01237", "01239", "01257", "01259", "01277", "01279", "01307", "01309", "01324", "01326", "01445", "01454", "01458", "01462", "01465", "01468", "01471", "01474", "01477", "01478", "01558", "01561", "01587", "01589", "01591", "01594", "01609", "01612", "01616", "01619", "01623", "01640", "01662", "01665", "01683", "01689", "01705", "01723", "01728", "01731", "01734", "01737", "01738", "01744", "01762", "01768", "01773", "01774", "01776", "01778", "01796", "01809", "01814", "01816", "01819", "01824", "01825", "01827", "01829", "01833", "01844", "01847", "01848", "01855", "01877", "01896", "01900", "01904", "01906", "01909", "01917", "01920", "01936", "01945", "01968", "01979", "01983", "01987", "01990", "01993", "01994", "01996", "01998"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := staedte.SearchByPostalCodePrefix(tt.plzprefix)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Staedte.PrefixPostleitzahl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaedte_PrefixStadtname(t *testing.T) {
	staedte := NewStaedte("postleitzahlen.txt")
	tests := []struct {
		name        string
		stadtprefix string
		want        []string
	}{
		{
			name:        "Test with prefix Dresden",
			stadtprefix: "Dre",
			want:        []string{"Dreisen", "Dreieich", "Drensteinfurt", "Drelsdorf", "Drebber", "Drebach", "Drebkau", "Dresden"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := staedte.SearchByCityNamePrefix(tt.stadtprefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Staedte.PrefixStadtname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaedte_GetAirDistance(t *testing.T) {
	staedte := NewStaedte("postleitzahlen.txt")
	tests := []struct {
		name   string
		stadt1 Stadt
		stadt2 Stadt
		want   float64
	}{
		{
			name:   "Berlin to Dresden",
			stadt1: Stadt{Loc_id: "12910", Plz: "10115", Lon: "13.387223860255", Lat: "52.5337069545604", Ort: "Berlin"},
			stadt2: Stadt{Loc_id: "5078", Plz: "01067", Lon: "13.7210676148814", Lat: "51.0600336463379", Ort: "Dresden"},
			want:   173.0, // Approximate distance in km
		},
		{
			name:   "Same city",
			stadt1: Stadt{Loc_id: "5078", Plz: "01067", Lon: "13.7210676148814", Lat: "51.0600336463379", Ort: "Dresden"},
			stadt2: Stadt{Loc_id: "5078", Plz: "01067", Lon: "13.7210676148814", Lat: "51.0600336463379", Ort: "Dresden"},
			want:   0.0,
		},
		{
			name:   "Invalid coordinates",
			stadt1: Stadt{Loc_id: "1", Plz: "1", Lon: "invalid", Lat: "51.0600336463379", Ort: "Invalid"},
			stadt2: Stadt{Loc_id: "5078", Plz: "01067", Lon: "13.7210676148814", Lat: "51.0600336463379", Ort: "Dresden"},
			want:   -1.0,
		},
		{
			name:   "Hamburg to München",
			stadt1: Stadt{Loc_id: "1", Plz: "20095", Lon: "9.993682", Lat: "53.551086", Ort: "Hamburg"},
			stadt2: Stadt{Loc_id: "2", Plz: "80331", Lon: "11.581981", Lat: "48.135125", Ort: "München"},
			want:   613.0, // Approximate distance in km
		},
		{
			name:   "Lübeck to Bönen",
			stadt1: Stadt{Loc_id: "3", Plz: "23552", Lon: "10.686559", Lat: "53.865467", Ort: "Lübeck"},
			stadt2: Stadt{Loc_id: "4", Plz: "59199", Lon: "7.755704", Lat: "51.598809", Ort: "Bönen"},
			want:   320.0, // Approximate distance in km
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := staedte.GetAirDistance(tt.stadt1, tt.stadt2)
			if tt.want == 0.0 && got != 0.0 {
				t.Errorf("Staedte.GetAirDistance() = %v, want %v", got, tt.want)
			} else if tt.want == -1.0 && got != -1.0 {
				t.Errorf("Staedte.GetAirDistance() = %v, want %v", got, tt.want)
			} else if tt.want > 0 && (got < tt.want*0.95 || got > tt.want*1.05) {
				t.Errorf("Staedte.GetAirDistance() = %v, want approx %v (±5%%)", got, tt.want)
			}
		})
	}
}

func TestStaedte_Stadtname(t *testing.T) {
	staedte := NewStaedte("postleitzahlen.txt")
	tests := []struct {
		name  string
		stadt string
		want  []Stadt
	}{
		{
			name:  "Test with Darmstadt",
			stadt: "Darmstadt",
			want:  []Stadt{{Loc_id: "9182", Plz: "64283", Lon: "8.65330818990971", Lat: "49.8715173852154", Ort: "Darmstadt"}},
		},
		{
			name:  "Test with Leipzig",
			stadt: "Leipzig",
			want:  []Stadt{{Loc_id: "5245", Plz: "04105", Lon: "12.3615440257456", Lat: "51.3532158481417", Ort: "Leipzig"}},
		},
		{
			name:  "Test with Berlin",
			stadt: "Berlin",
			want:  []Stadt{{Loc_id: "12910", Plz: "10115", Lon: "13.387223860255", Lat: "52.5337069545604", Ort: "Berlin"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := staedte.FindCitiesByName(tt.stadt)
			found := false
			for _, s := range got {
				if s.Loc_id == tt.want[0].Loc_id && s.Plz == tt.want[0].Plz && s.Lon == tt.want[0].Lon && s.Lat == tt.want[0].Lat && s.Ort == tt.want[0].Ort {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Staedte.Stadtname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaedte_Postleitzahl(t *testing.T) {
	staedte := NewStaedte("postleitzahlen.txt")
	tests := []struct {
		name string
		plz  string
		want Stadt
	}{
		{
			name: "Test with 01067",
			plz:  "01067",
			want: Stadt{Loc_id: "5078", Plz: "01067", Lon: "13.7210676148814", Lat: "51.0600336463379", Ort: "Dresden"},
		},
		{
			name: "Test with 01445 (shared by multiple cities)",
			plz:  "01445",
			want: Stadt{Loc_id: "5103", Plz: "01445", Lon: "13.6432842540266", Lat: "51.113631945926", Ort: "Radebeul"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := staedte.FindCityByPostalCode(tt.plz)
			found := false

			if got.Loc_id == tt.want.Loc_id && got.Plz == tt.want.Plz && got.Lon == tt.want.Lon && got.Lat == tt.want.Lat && got.Ort == tt.want.Ort {
				found = true
			}

			if !found {
				t.Errorf("Staedte.Postleitzahl() = %v, want %v", got, tt.want)
			}
		})
	}
}
