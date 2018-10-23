package parag

import "testing"



/*
TrackStorage is an interface serving all access
to Track data
*/
type Test_TrackStorage interface {
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
type Test_TrackDB struct {
	// Database for metadata for flight track records
	tracks  map[RegTrack]Track
	// Look-up map for simplified checking of existing
	// urls in database and allocating of new ids
	urlkeys map[string]int
	// Map of API elements for checking incoming text
	fields  map[string]int
}



// TODO This is not functioning because of a lot of dependencies on structs
//to do a lot of other tests to get this to work

 func Test_Init(t *testing.T) {
   // db.tracks is the database for track data
   t.tracks  = make(map[RegTrack]Track)
   // db.urlkeys store urls as keys for easy look-up of ids
   t.urlkeys = make(map[string]int)
   // db.fields HARDCODED of valid API words (used instead of regexp check)
   t.fields  = map[string]int{"pilot":1,"glider":2,"glider_id":3,"track_length":4,"h_date":5}

    // Dummies to fill the databases (maps)
   	s1 := RegTrack{TrID: 1, TrURL: "http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc", Track: Track{HDate: "2016-10-05", Pilot: "Siv Toppers", Glider: "Mypmyp", GliderID: "AIKK-3", TrackLength: 64.1}}
   	s2 := RegTrack{TrID: 2, TrURL: "http://example.com/igcinfo/api/igc/track2", Track: Track{HDate: "2015-11-10", Pilot: "Vanja Falck", Glider: "Ompa", GliderID: "AIKK-5", TrackLength: 23.2}}
   	s3 := RegTrack{TrID: 3, TrURL: "http://example.com/igc/track3", Track: Track{HDate: "2017-04-09", Pilot: "Marius Muller", Glider: "Theodor", GliderID: "AIKK-12", TrackLength: 46.4}}

   	t.tracks[s1] = s1.Track
   	t.tracks[s2] = s2.Track
   	t.tracks[s3] = s3.Track

   	t.urlkeys[s1.TrURL] = 1
   	t.urlkeys[s2.TrURL] = 2
   	t.urlkeys[s3.TrURL] = 3

    if ((t.tracks[s1].TrURL != "http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc") || (t.tracks[s2].Track.H_date != "2015-11-10") || (t.tracks[s3].Track.Pilot  != "Marius Muller") ||  (t.tracks[s1].Track.Glider != "Mypmyp") ||  (t.tracks[s1].Track.Glider_id != "AIKK-3") || (t.trakcs[s3].Track.TrackLength != 46.4)){
      t.Error("The entry to main database (db.tracks) did not work")
      }
    if t.urlkeys["http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc"] != s1.TrURL {
      t.Error("The key-map does not work")
    }
 }
