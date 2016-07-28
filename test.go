package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

//Geo struct fk u go
type Geo struct {
	AddressComponents []struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	} `json:"address_components"`
	FormattedAddress string `json:"formatted_address"`
	Geometry         struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
		LocationType string      `json:"location_type"`
		Types        interface{} `json:"types"`
		Viewport     struct {
			Northeast struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"northeast"`
			Southwest struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"southwest"`
		} `json:"viewport"`
	} `json:"geometry"`
	PlaceID string   `json:"place_id"`
	Types   []string `json:"types"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		var loc maps.LatLng
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		if r.Method == "GET" {
			fmt.Fprintln(w, "GET if")
		}
		if r.Method == "POST" {
			err := json.NewDecoder(r.Body).Decode(&loc)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			lati := strconv.FormatFloat(loc.Lat, 'f', -1, 64)
			long := strconv.FormatFloat(loc.Lng, 'f', -1, 64)
			fmt.Println(lati + "," + long)
			//fmt.Fprintln(w, "lat:", lati, "\nlong:", long)
			geoResult, err1 := geocode(loc.Lat, loc.Lng)

			if err1 != nil {
				http.Error(w, err1.Error(), 400)
				return
			}
			output, err2 := json.MarshalIndent(geoResult, "", "     ")

			if err2 != nil {
				http.Error(w, err2.Error(), 400)
				return
			}
			//s := string(output[:])
			var g []Geo
			err3 := json.Unmarshal(output, &g)
			if err3 != nil {
				http.Error(w, err3.Error(), 400)
				return
			}

			fmt.Fprintln(w, "FormattedAddress =  ", g[0].FormattedAddress)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func geocode(x float64, y float64) ([]maps.GeocodingResult, error) {
	c, err := maps.NewClient(maps.WithAPIKey("%API_KEY%"))
	if err != nil {
		return nil, err
	}
	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: x, Lng: y},
	}
	resp, err := c.Geocode(context.Background(), r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
