package main

import (
	"context"
	"log"
	"strings"

	"net/http"

	"io/ioutil"

	"encoding/json"

	"os"

	"fmt"

	"googlemaps.github.io/maps"
)

var key = os.Getenv("WAW_GOOGLE_API_KEY")
var icon = os.Getenv("WAW_CUSTOM_ICON")

func addressesToLocation(c maps.Client, addresses []string) []maps.LatLng {
	r := []maps.LatLng{}
	for _, a := range addresses {
		gcr, err := c.Geocode(context.Background(), &maps.GeocodingRequest{Address: a})

		if err != nil {
			log.Fatal(err)
		}

		if len(gcr) == 0 {
			log.Fatal(err)
		}

		r = append(r, gcr[0].Geometry.Location)
	}

	return r
}

func generateMap(c *maps.Client, addrs []maps.LatLng) []byte {
	r := &maps.StaticMapRequest{
		Size:    "1024x768",
		MapType: maps.RoadMap,
		Markers: []maps.Marker{
			maps.Marker{
				Location: addrs,
			},
		},
	}

	if icon != "" {
		r.Markers[0].CustomIcon.IconURL = icon
	}
	resp, err := c.StaticMap(context.Background(), r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

type addrs struct {
	Addrs []string
}

func main() {

	if key == "" {
		fmt.Fprintln(os.Stderr, "Key is not provided")
		os.Exit(1)
	}

	c, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	sm := http.NewServeMux()
	sm.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			rw.Write(generateMap(c, addressesToLocation(*c, strings.Split(req.URL.Query().Get("addr"), "|"))))
		case http.MethodPost:
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			a := addrs{}
			json.Unmarshal(b, &a)
			rw.Write(generateMap(c, addressesToLocation(*c, a.Addrs)))
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}

	})
	if err := http.ListenAndServe(":8080", sm); err != nil {
		panic(err)
	}

}
