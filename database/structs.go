package database

type DatabaseInfo struct {
	Initialized              bool
	UninitializedCollections []struct {
		Name     string
		Recovery string
	}
}

type Pending struct {
	SID           string
	DateRetreived string
}

type Record struct {
	Pending
	DateReviewed string
}

type Theme struct {
	Name      string `xml:"name,attr"`
	SubThemes []struct {
		Name        string   `xml:"name,attr"`
		Description string   `xml:"description"`
		Keywords    []string `xml:"keywords>kw"`
	} `xml:"subtheme"`
}

type Faculty struct {
	GivenName string
	Surname   string
	Title     string
	SID       []string
}

type Department struct {
	Name string
	SIDs []string
}

type Citation struct {
	Title        string
	PubType      string
	PubName      string
	SubType      string
	Volume       string
	Pages        string
	Date         string
	ISODate      string
	DOI          string
	Abstract     string
	CitedByCount int
	Keywords     []string
	SID          string
	Authors      []Author
	Affiliations []Affiliation
	Status       Status
}

type Author struct {
	Faculty
	Local       bool
	Affiliation string
}

type Affiliation struct {
	SID     string
	Name    string
	City    string
	Country string
}

type Status int

const (
	STATUS_UNCONFIRMED Status = iota
	STATUS_CONFIRMED
	STATUS_EXCLUDED
)
