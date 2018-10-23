package parag

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

/*
GlobalDB is based on session memory
Can be exchanged with other types of databases
*/
var GlobalDB TrackStorage

/*
RegTrack object has all information about the track records.
RegTrack is the key in the database which keep Track objects.
*/
type RegTrack struct {
	Track
	TrURL string `json:"url"`
	TrID  int    `json:"id"` // Unique /igcinfo/api/<id>/
}

/*
Track object stores all flight track information.
Track is embedded in the RegTrack object
*/
type Track struct {
	HDate       string  `json:"H_date"`       // Date Header H-record
	Pilot       string  `json:"pilot"`        // Pilots name
	Glider      string  `json:"glider"`       // Glider Type
	GliderID    string  `json:"glider_id"`    // Glider ID
	TrackLength float64 `json:"track_length"` // Calculated length (km)
}

// Use the igc Points struct.
type Points struct {

}

/*
NB DOES NOT USE THIS
ServiceIGC is an object containing package information and
uptime for the web-service.
*/
type ServiceIGC struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

/*
TrackStorage is an interface serving all access
to Track data
*/
type TrackStorage interface {
	Init()
	Add(reg RegTrack) error
	Count() int
	FindRegTrack(reg RegTrack) (RegTrack, bool)
	GetAll() []RegTrack
	TransformIGC(u string) (RegTrack, bool)
	ApiFields(ap string) bool
	//NewID() (string)
}
/*
TrackDB takes the RegTrack object as key for maps
keeping easy look up
TODO maybe this is redundant
*/
type TrackDB struct {
	// Database for metadata for flight track records
	tracks  map[RegTrack]Track
	// Look-up map for simplified checking of existing
	// urls in database and allocating of new ids
	urlkeys map[string]int
	// Map of API elements for checking incoming text
	fields  map[string]int
}
/*
Init() initiates an empty database with
session memory storage as a map with RegTrack as key
*/
func (db *TrackDB) Init() {
	// db.tracks is the database for track data
	db.tracks  = make(map[RegTrack]Track)
	// db.urlkeys store urls as keys for easy look-up of ids
	db.urlkeys = make(map[string]int)
	// db.fields HARDCODED of valid API words (used instead of regexp check)
	db.fields  = map[string]int{"pilot":1,"glider":2,"glider_id":3,"track_length":4,"h_date":5}

/*
HARDCODED TEST input to database
	s1 := RegTrack{TrID: 1, TrURL: "http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc", Track: Track{HDate: "2016-10-05", Pilot: "Siv Toppers", Glider: "Mypmyp", GliderID: "AIKK-3", TrackLength: 764}}
	s2 := RegTrack{TrID: 2, TrURL: "http://example.com/igcinfo/api/igc/track2", Track: Track{HDate: "2015-11-10", Pilot: "Vanja Falck", Glider: "Ompa", GliderID: "AIKK-5", TrackLength: 223}}
	s3 := RegTrack{TrID: 3, TrURL: "http://example.com/igc/track3", Track: Track{HDate: "2017-04-09", Pilot: "Marius Muller", Glider: "Theodor", GliderID: "AIKK-12", TrackLength: 346}}

	db.tracks[s1] = s1.Track
	db.tracks[s2] = s2.Track
	db.tracks[s3] = s3.Track

	db.urlkeys[s1.TrURL] = 1
	db.urlkeys[s2.TrURL] = 2
	db.urlkeys[s3.TrURL] = 3
*/

}

func (db *TrackDB) ApiFields(ap string) bool {
  _, ok := db.fields[ap]
	 if !ok {
		 return false
	 }
return true
}

/*
TransformIGC takes an url as a string and check if
it is a valid .igc flight record file. It it is, the content
is extracted. The track is stored as a RegTrack object
(url and trackID) with Tracks (metadata). Using Goigc/igc.
*/
func (db *TrackDB) TransformIGC(u string) (RegTrack, bool) {

	// TODO This is a clumsy way of doing this, have to restructure
	// this function to give better error messages...
	// Check if url is already registered. If it is, returns empty
	// RegTrack object and false.
	// NB this is blocking for cheking if parsing is not working as there
	// is no good debug handling here. (Should skip this when testing parsing)
	// An example: Get an error that "Track already registered - if parsing fails"

	/*
	_, ok := db.urlkeys[u]
	if ok { return RegTrack{}, false }

	*/

/*
HARDCODED for testing:
u = "http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc"
u = "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
*/

	// Goigc/igc package parsing the flight track data from .igc file
	track, err := igc.ParseLocation(u)
	if err != nil {
/*
TODO change this to an internal errorhandling.
Current function only returns a RegTrack object - could be
an empty one
Treat the ERROR  - do not use err (no formatting directives).
Should include an error as return value instead of bool!
*/
		fmt.Errorf("Problem reading the received track file. Are you sure it is an .igc?", err)
		return RegTrack{}, false
	}
	var dist float64
	// TODO implement function from received points.
	// Distance() is not working has to calculate from points
	dist = track.Task.Distance()


	// Check if url is already registered.
	// Parsing into new RegTrack object if not yet registered.
	_, ok := db.urlkeys[u]
	if !ok {
	// GET A NEW ID into newId (quick-fix)
	// TODO This will not work when deleting tracks!
	// Not concurrency proof
	  newNumber := len(db.urlkeys) + 1
	  // New id is assigned to the db.urlkeys map
	  db.urlkeys[u] = newNumber

		// Formatting the igc-data into RegTrack object according to JSON
		// Returning as a RegTrack-object reg
		reg := RegTrack{TrID: newNumber, TrURL: u, Track: Track{HDate: track.Date.String(), Pilot: track.Pilot, Glider: track.GliderType, GliderID: track.GliderID, TrackLength: dist}}

		return reg, true
	}
	// Empty RegTrack
	return RegTrack{}, false
}

// Add RegTrack objects to the database
func (db *TrackDB) Add(reg RegTrack) error {
	db.tracks[reg] = reg.Track
	return nil
}

// Count the number of objects in the database
func (db *TrackDB) Count() int {
	return len(db.tracks)
}

/*
FindRegTrack takes a RegTrack (key in database) if
match on either the url or the id of the object, a
complete RegTrack object with flight records is returned.
If RegTrack object is not complete - false, else true.
*/
func (db *TrackDB) FindRegTrack(reg RegTrack) (RegTrack, bool) {
  // TODO better handling of errors - may be restructure return values
  // Maybe make a function isInDB() to check if a reg is present and
	// make error handling based on this (as this checking is done several times)
	// Fetch Track object with key == tr
	// If track exist:
	// Return the RegTrack object with id/url and track
switch {
	//Get track data directly from DATABASE
  case reg.TrID > 0 && reg.TrURL != "":
	// Checks if reg is in database:
	_ , ok := db.tracks[reg]
	if !ok {
		return reg, false
	}
	return reg, true

	// Get track data from GET api/<track-00xx> by id
	// To display metadata and separate fields
  case reg.TrURL == "":
	// Loops the database to check based on the track id
	if reg.TrID > 0 {
		var key RegTrack
		for k, trc := range db.tracks {
			if trc == db.tracks[k] {
				key = k
				if key.TrID == reg.TrID {
					return key, true
				}
			}
		 }
	 }
	// Get track data from POST api/igc
	// To transform by goigc/igc package and
	// stored in database and displayed
  case reg.TrURL != "":
	if reg.TrID == 0 {
	// Get the track id from db.urlkeys
	numID    := db.urlkeys[reg.TrURL]
	reg.TrID  = numID
	reg.Track = db.tracks[reg]
	return reg, true
  }
 }
 // IF no option in switch - the checked reg cannot be confirmed
return reg, false
}

/*
GetAll() returns a slice of all RegTrack objects (the keys
in the database).
*/
func (db *TrackDB) GetAll() []RegTrack {
	// Makes a slice at the size of db
	gAll := make([]RegTrack, 0, db.Count())

	// Returns all Track-objects:
	//for _, reg := range db.tracks {

	// Returns all keys as RegTrack objects
	for k := range db.tracks {
		gAll = append(gAll, k)
	}
	return gAll
}

// DISPLAY api/igc ---> array of all Tracks of RegTracks (incl url, id)
func displayAllRegTracks(w http.ResponseWriter, db TrackStorage) {
	if db.Count() == 0 {
		json.NewEncoder(w).Encode([]RegTrack{})
	} else {
		add := make([]RegTrack, 0, db.Count())
		// Retreive all Track-objects by RegTrack key
		for _, trc := range db.GetAll() {
			add = append(add, trc)
		}
		json.NewEncoder(w).Encode(add)
	}
}

// DISPLAY api/igc ---> STRING array of all IDs
func displayAllRegTracksID(w http.ResponseWriter, db TrackStorage) {
	if db.Count() == 0 {
		json.NewEncoder(w).Encode([]string{})
	} else {
		add := make([]string, 0, db.Count())
		// Retreive all RegTracks (=key for Tracks)
		for _, reg := range db.GetAll() {
			add = append(add, strconv.Itoa(reg.TrID))
		}
		json.NewEncoder(w).Encode(add)
	}
}

// DISPLAY api/igc/<id> ---> one Track object (unnested) displayed
// The query for <id> identified by the split url response "u"
func displayOneRegTrack(w http.ResponseWriter, db TrackStorage, reg RegTrack, field string) {
	//http.Header.Add(w.Header(), "date", time.RFC3339)
	// Find the reg in db with missing either RegTrack.Id or u (string)
	trc, ok := db.FindRegTrack(reg)
	// If the reg cannot be confirmed as belonging to the database
	// it returns !ok (false)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch {
	case field == "":
		http.Header.Add(w.Header(), "Content-Type", "application/json")
		json.NewEncoder(w).Encode(trc)
	case field == "pilot":
		//http.Header.Add(w.Header(), "Content-Type", "text/plain")
		json.NewEncoder(w).Encode(trc.Pilot)
	case field == "glider":
	//	http.Header.Add(w.Header(), "Content-Type", "text/plain")
		json.NewEncoder(w).Encode(trc.Glider)
	case field == "glider_id":
		//http.Header.Add(w.Header(), "Content-Type", "text/plain")
		json.NewEncoder(w).Encode(trc.GliderID)
	case field == "h_date":
		//http.Header.Add(w.Header(), "content-type", "text/plain")
		json.NewEncoder(w).Encode(trc.HDate)
	case field == "track_length":
		//http.Header.Add(w.Header(), "Content-Type", "text/plain")
		json.NewEncoder(w).Encode(trc.TrackLength)
	default:
		//http.Header.Add(w.Header(), "Content-Type", "text/plain")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

// TODO: change this to fit the router Gin framework:
// Serves api/igc POST and GET

/*
HandlerRegTrack serves the /igcinfo/api/igc/ path
*/
func HandlerRegTrack(w http.ResponseWriter, r *http.Request) {
  http.Header.Add(w.Header(), "Date", time.RFC3339)
	//http.Header.Add(w.Header(), "content-type", "application/json")
	// Can also use:
	w.Header().Add("content-type", "application/json")

	switch r.Method {
	// Receiving url as json {"url":"<url>"} in r.Body
	case "POST":
		// Placeholder for incomming url
		var urlIGC string
		// New RegTrack object to recieve the content of r.Body
		var reg RegTrack
		err := json.NewDecoder(r.Body).Decode(&reg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
    urlIGC = reg.TrURL

/*
REGEX checking of appropriate url-format:
Tested regexp in http://regexr.com
Possibly it has to be a http and not https?
if both http/https: http?
pattern := `^https?:\/\/[^\s$.?#].[^\s]*igc$`
If http only: http:
*/
		pattern := `^http?:\/\/[^\s$.?#].[^\s]*igc$`
		rex := regexp.MustCompile(pattern)
		match := rex.MatchString(urlIGC)
		if !match {
			http.Error(w, "The url is not correctly formatted as .igc, please try again.", http.StatusBadRequest)
			return
		}

		// Retrieve the correctly formatted url as string from received RegTrack object


/*
IF Deploying with GIN:
ok, err2 := regexp.MatchString(pattern, urlIGC)
if err2 != nil {
http.Error(w, err2.Error(), http.StatusBadRequest)
In GIN:
c.Error(err)
}
*/

		// Chek if track-object already exists
		// Returns a filled RegTrack object and retreive the
		// track record linked to the sendt urlIGC
		regTrans, ok := GlobalDB.TransformIGC(urlIGC)
		// If track url is already in database:
		if !ok {
			http.Error(w, "This track id is already registered.", http.StatusBadRequest)
			return
		}

		// If url is valid and not yet registered - the data Will
		// be added to the database.
		// This one is working:
		// This should never happen!!!
		if regTrans.TrURL == "" {
			http.Error(w, "This track is either EMPTY but not REJECTED!", http.StatusBadRequest)
			return
		}

			GlobalDB.Add(regTrans)
			// Formatted according to specifications {"id":"<id>"}
			fmt.Fprintf(w, "{\"id\":\"%d\"}", regTrans.TrID)
			return

		// GET for api/igc
		// Returns/display an array of all tracks ids as JSON:
	case "GET":
		//http.Header.Add(w.Header(), "content-type", "application/json")
		// Can also use:
		// w.Header().Add("content-type", "application/json")

		// ID as string array [] = empty
		displayAllRegTracksID(w, GlobalDB)

		// Track-objects as an array {} when empty
		//displayAllRegTracks(w, GlobalDB)

	default:
		http.Error(w, "This is not an option in this API.", http.StatusBadRequest)
		return
	}
}

/* TODO INCLUDE to safeguard concurrency writing to and from maps
SafeTrack sync.Mutex 		// Protects session
Maxlifetime int64				// Session duration
*/

/*
HandlerRegSingleTrack serves single track views
api/igc/track + number ONLY GET requests.
*/
func HandlerRegSingleTrack(w http.ResponseWriter, r *http.Request) {
 http.Header.Add(w.Header(), "Date", time.RFC3339)
	//http.Header.Add(w.Header(), "content-type", "application/json")
	// Can also use:
	//w.Header().Add("content-type", "application/json")

	switch r.Method {
	case "GET":
		// Getting the right path and checing for correct format/names
		// Making a RegTrack object to receive track id (TrID) from URL
		// part api/igc/<track-xxxx>
		var reg RegTrack
		// Make a slice of url path components
		// root empty/empty <api>/empty<igc>/empty<track-(id)>/empty/<field>
		// parts[2] = api, parts[3] = igc, parts[4] = track, part[5] = field
		parts := strings.Split(r.URL.Path, "/")
		// Pattern of <track-xxxx> with exactly 4 digits 0-9 (starting at 0001)
		// Changing from {3} to {2} decrease digits to 3
		var le int
		le = len(parts)

		if le <= 1 || le >= 8 {
			// TODO change this to a not found or bad request response!
			// Currently used for checking
			http.Error(w, "(OUT OF RANGE) The api-url is not correctly formatted.\n Use api/igc/<track-number>.\n Number is always 4 digits with zeros in front like 0022.", http.StatusBadRequest)
			return
		}

		//patternTrack := `^track-[0-9]{3}[0-9]$`
		//rex   := regexp.MustCompile(patternTrack)
		rex := regexp.MustCompile("^track-[0-9]{3}[0-9]$")
		match := rex.MatchString(parts[4])
		if !match {
			// TODO - give an appropriate error -> StatusBadRequest
			// CURRENTLY FOR TESTING ONLY:
			//fmt.Fprintf(w, "Parts: 1%s, 2%s,\n 3%s, 4%s,\n 5%s, 6%s, 7%s", parts[0],parts[1], parts[2], parts[3], parts[4], parts[5], parts[6])
			http.Error(w, "(DONT FIT REGEX) The url you ask for is not correct in this api. Correct format: api/igc/<track-xxx>/<field>.", http.StatusBadRequest)
			return
		}
		// Extracting TrID as an integer from the API url part <track-xxxx>
		numString := strings.Split(parts[4], "-")
		numInt, err := strconv.Atoi(numString[1])
		if err != nil {
			// TODO give this a StatusBadRequest
			// TESTING ONLY
			http.Error(w, "(CANT MAKE INT FROM STRING) The tracknumber is not formatted correctly: api/igc/<track-xxx> with 0´s in front of number.", http.StatusBadRequest)
			return
		}
		// Pass the trackid part of the url to a RegTrack object
		//reg := RegTrack{}
		reg.TrID = numInt
/*
NB This is handled in the displayOneRegTrack()
TODO Find the best way of poulating a RegTrack object from TrID (here or in display)

		// Check if the trackid is in database
		_, ok := GlobalDB.FindRegTrack(reg)
		if !ok {
			http.Error(w, "The track you asked for is not found.", http.StatusNotFound)
			return
		}
*/

		// ALL GET REQUESTS for api/igc/<track-xxxx>
		// Returns/display an array a single track by id as JSON
     if le == 5 {
			//http.Header.Add(w.Header(), "Content-Type", "application/json")
			//w.Header().Add("content-type", "application/json")
			// Return and display the trackrecord data as JSON in body response
			displayOneRegTrack(w, GlobalDB, reg, "")

		}
			// Printing confirm that the number conversion is ok and that parts 4 = track-00xx
			//fmt.Fprintf(w, "Did the string number convert? String:%s, Int:%d\n Parts[4]: %s.", numString[1], numInt, parts[4])



			// Get the /api/igc/track-xxxx/<field>
		  if le == 6 || le == 7 {
			http.Header.Add(w.Header(), "Content-Type", "text/plain")
			//var field string
			field := strings.ToLower(parts[5])
			// Checing if field name is valid (present in a map)
      if !GlobalDB.ApiFields(field){
				http.Error(w, "The field you asked for is not in this API.", http.StatusBadRequest)
				return
				}
				// ONLY FOR TEST:
				//fmt.Fprintf(w,"Field: Tekst:%s, Type:%t,\n Variable:%v, Parts[5]:%s, Parts[4]:%s", field, field, field, parts[5], parts[4])
				//fmt.Fprintf(w, "Parts: 1%s, 2%s,\n 3%s, 4%s,\n 5%s, 6%s, 7%s", parts[0],parts[1], parts[2], parts[3], parts[4], parts[5], parts[6])
      fmt.Fprintf(w, "track-000%s \npilot: ", strconv.Itoa(reg.TrID))
			displayOneRegTrack(w, GlobalDB, reg, field)
      }
	default:
		http.Error(w, "This is not an option in this API.", http.StatusBadRequest)
		return
	}
}

func HandlerApiInfo(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("content-type", "application/json")
  http.Header.Add(w.Header(), "Date", time.RFC3339)

// NB uptime is not working:
var start time.Duration
var step time.Time
//step  = `2006-01-02T15:04:05Z07:00`
start = 0
step = time.Now()
http.Header.Add(w.Header(), "Uptime", start.String())
//t =: NewTimeCounter()
// returns a *Timer object - will be called on by its own goroutine
//uptime1 := time.AfterFunc(1 * time.Year, f func())
//start := 0
//step := time.Now()
//t := NewTimeCounter()

t := NewTimeCounter()
t(&start, &step)
http.Header.Add(w.Header(), "Uptime", start.String())
	fmt.Fprintf(w, "{\"uptime\":\"%v\" og %v\n", &start, &step)
	fmt.Fprintf(w, " \"info\":%s\n", "\"Service for IGC tracks\"")
  fmt.Fprintf(w, " \"version\":%s\n}", "\"v1\"")
}

/*
TimeCounter register the uptime of the web service.
*/
func NewTimeCounter() func(start *time.Duration, step *time.Time) time.Duration {
	return func(start *time.Duration, step *time.Time) time.Duration {
		*start = *start + step.Sub(*step)
		return *start + step.Sub(*step)
	}
}

/*
TimeCounter register the uptime of the web service.
*/
func TimeCounter() func(start *time.Duration, step *time.Time) time.Duration {
	return func(start *time.Duration, step *time.Time) time.Duration {
		*start = *start + time.Since(*step)
		return *start + time.Since(*step)
	}
}
