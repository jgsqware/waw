package main

import (
	"context"
	"log"

	"net/http"

	"io/ioutil"

	"googlemaps.github.io/maps"
)

const key = "AIzaSyCBaQfn-uQpppEynyr6Qvb9Yu_8Pahh14k"

func main() {
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		r := &maps.StaticMapRequest{
			Center:  "Brooklyn Bridge,New York,NY",
			Zoom:    13,
			Size:    "600x300",
			MapType: maps.RoadMap,
			Markers: []maps.Marker{
				maps.Marker{
					CustomIcon: maps.CustomIcon{
						IconURL: "https://pbs.twimg.com/profile_images/508895937817608193/g5fQgDRJ_400x400.png",
					},
					Location: []maps.LatLng{
						maps.LatLng{
							Lat: 40.702147,
							Lng: -74.015794,
						},
					},
				},
			},
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
		rw.Write(b)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
