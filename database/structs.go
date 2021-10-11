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
	CitedByCount uint
	Keywords     []string
	SID          string
	Authors      []struct {
		Faculty
		Local       bool
		Affiliation string
	}
	Affiliations []struct {
		SID     string
		Name    string
		City    string
		Country string
	}
	Status Status
}

type Status uint

const (
	STATUS_UNCONFIRMED Status = iota
	STATUS_CONFIRMED
	STATUS_EXCLUDED
)
