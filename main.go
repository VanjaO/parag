/* IGC url browser
   Building on marni/goigc/igc IGC parser
   Author Vanja Falck

Sources used to develop the code:

Testing for API in go-router:
https://github.com/appleboy/gofight

REST-API testing:
https://github.com/gavv/httpexpect

Security for memory ect:
https://github.com/awnumar/memguard

About GET, POST etc with gin router in golang:
https://github.com/gin-gonic/gin#using-get-post-put-patch-delete-and-options

IGC format and files: https://aerofiles.readthedocs.io/en/latest/guide/igc-writing.html
http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc

*/

package main

import (
	"fmt"
	"net/http"
  "strconv"
	"strings"
	"encoding/json"
	"regexp"
	"time"
	//"github.com/gin-gonic/gin"
	//	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/marni/goigc"
  "github.com/vanjao/parag"
)

func main() {

	GlobalDB = &TrackDB{}

	/*
		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("$PORT must be set")
		}
	*/

	GlobalDB.Init()

	/* Initiate with dummies to test:
	   s1 := RegTrack{TrID: "track1", TrURL: "/igcinfo/api/igc/track1", Track: Track{HDate: "2016-10-05", Pilot: "Siv Toppers", Glider: "Mypmyp", GliderId: "AIKK-3", TrackLength: 764}}
	   s2 := RegTrack{TrID: "track2", TrURL: "/igcinfo/api/igc/track2", Track: Track{HDate: "2015-11-10", Pilot: "Vanja Falck", Glider: "Ompa", GliderId: "AIKK-5", TrackLength: 223}}
	   s3 := RegTrack{TrID: "track3", TrURL: "/igcinfo/api/igc/track3", Track: Track{HDate: "2017-04-09", Pilot: "Marius Muller", Glider: "Theodor", GliderId: "AIKK-12", TrackLength: 346}}

	   GlobalDB.Add(s1)
	   GlobalDB.Add(s2)
	   GlobalDB.Add(s3)

	   fmt.Println("GlobalDB content: ", GlobalDB)
	*/
	//port := os.Getenv("PORT")
	http.HandleFunc("/igcinfo/api", HandlerApiInfo)
	http.HandleFunc("/igcinfo/api/igc", HandlerRegTrack)
	http.HandleFunc("/igcinfo/api/igc/", HandlerRegSingleTrack)
	//http.ListenAndServe(":"+port, nil)
	http.ListenAndServe("127.0.0.1:8809", nil)
}
