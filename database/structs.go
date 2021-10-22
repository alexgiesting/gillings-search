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
	Abbr      string `xml:"abbr,attr"`
	SubThemes []struct {
		Name        string   `xml:"name,attr"`
		Abbr        string   `xml:"abbr,attr"`
		Description string   `xml:"description"`
		Keywords    []string `xml:"keywords>kw"`
	} `xml:"subtheme"`
}

type Faculty struct {
	Name      string
	Title     string
	SID       []string
	Email     string
	Strengths []Strength
}

type Strength struct { // TODO
	Theme    string
	SubTheme string
}

type Department struct { // TODO
	Name string
	SIDs []string
}

type Citation struct {
	Title        string
	PubType      string
	PubName      string
	SubType      string
	Volume       string
	Issue        string
	Pages        string
	Date         string
	ISODate      string
	DOI          string
	Abstract     string
	CitedByCount int
	Keywords     []string
	EID          string
	Authors      []Author
	Affiliations []Affiliation
	Status       Status
}

type Author struct {
	Name      string
	GivenName string
	Surname   string
	Initials  string
	SID       string
	AffilIDs  []string
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
