package query

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/alexgiesting/gillings-search/go/database"
	"go.mongodb.org/mongo-driver/bson"
)

func TestQuery(t *testing.T) {
	var q bson.M
	// ( AND AUTHOR(chen)
	//       ( OR
	//         ( AND KEYWORD(circulat)
	//               ( OR KEYWORD(pulmon) KEYWORD(cardi) )
	//               ( OR AUTHOR(hayes) AUTHOR(beth) )
	//         )
	//         ( OR
	//           ( AND AUTHOR(kur) AUTHOR(for) KEYWORD(birth) )
	//           ( AND AUTHOR(feng) KEYWORD(aging) )
	//           ( AND AUTHOR(mayer) KEYWORD(grow) KEYWORD(matur) )
	//         )
	//       )
	// )
	err := json.Unmarshal([]byte(`{
		"faculty": ["chen"],
		"or": [
			[
				{
					"keyword": ["circulat"],
					"or": [
						[ { "keyword": ["pulmon"] }, { "keyword": ["cardi"]  } ],
						[ { "faculty": ["hayes"] }, { "faculty": ["beth"]  } ]
					]
				},
				{
					"or": [
						[
							{ "faculty": ["kur", "for"], "keyword": ["birth"]         },
							{ "faculty": ["feng"],       "keyword": ["aging"]         },
							{ "faculty": ["mayer"],      "keyword": ["grow", "matur"] }
						]
					]
				}
			]
		]
	}`), &q)
	if err != nil {
		t.Logf("Failed to parse test json: %v", err)
		t.FailNow()
	}

	db := database.Connect()
	db.Citations.Drop(context.TODO())

	s := func(strings ...string) []string {
		return strings
	}
	c := func(eid string, names []string, title string, abstract string, keywords []string) database.Citation {
		authors := make([]database.Author, len(names))
		for i, name := range names {
			authors[i].Name = name
		}
		return database.Citation{
			EID:      eid,
			Authors:  authors,
			Title:    title,
			Abstract: abstract,
			Keywords: keywords,
		}
	}
	citations := []database.Citation{
		c("a", s("Chen, J.", "Beth, M.N."), "Study on Circulation", "We studied the cardiological symptoms of ...", s("sleep")),
		c("b", s("XXXX, J.", "Beth, M.N."), "Study on Circulation", "We studied the cardiological symptoms of ...", s("sleep")),
		c("c", s("Chen, J.", "XXXX, M.N."), "Study on Circulation", "We studied the cardiological symptoms of ...", s("sleep")),
		c("d", s("Chen, J.", "Beth, M.N."), "Study on XXXXXXXXXXX", "We studied the cardiological symptoms of ...", s("sleep")),
		c("e", s("Chen, J.", "Beth, M.N."), "Study on Circulation", "We studied the XXXXXXXXXXXXX symptoms of ...", s("sleep")),

		c("f", s("chen, j.", "hayes, a.f."), "How to limit symptoms", "The pulmonary arteries play a major role in circulating ...", s()),
		c("g", s("XXXX, j.", "hayes, a.f."), "How to limit symptoms", "The pulmonary arteries play a major role in circulating ...", s()),
		c("h", s("chen, j.", "XXXXX, a.f."), "How to limit symptoms", "The pulmonary arteries play a major role in circulating ...", s()),
		c("i", s("chen, j.", "hayes, a.f."), "How to limit symptoms", "The XXXXXXXXX arteries play a major role in circulating ...", s()),
		c("j", s("chen, j.", "hayes, a.f."), "How to limit symptoms", "The pulmonary arteries play a major role in XXXXXXXXXXX ...", s()),

		c("k", s("CHEN, J.", "BETHAM, E.R."), "Circulatory conditions", "In this study we observed ...", s("cardiac arrest")),
		c("l", s("XXXX, J.", "BETHAM, E.R."), "Circulatory conditions", "In this study we observed ...", s("cardiac arrest")),
		c("m", s("CHEN, J.", "XXXXXX, E.R."), "Circulatory conditions", "In this study we observed ...", s("cardiac arrest")),
		c("n", s("CHEN, J.", "BETHAM, E.R."), "XXXXXXXXXXX conditions", "In this study we observed ...", s("cardiac arrest")),
		c("o", s("CHEN, J.", "BETHAM, E.R."), "Circulatory conditions", "In this study we observed ...", s("XXXXXXX arrest")),

		c("p", s("Yi Chen, K.", "Kurtiss, M.K.N.", "Forthe, M."), "Neonates: Development", "Medical technology has improved ...", s("infant care", "premature births")),
		c("q", s("Hao, R.W.", "Feng, E.", "Chen, M.", "Kai, S."), "Healthcare facilities", "In regions with aging hospitals ...", s("infrastructure", "regional cost")),
		c("r", s("Mayer, Jonathan", "Ma, Chengdu", "Argent, J."), "Maturation and Growth", "The development of the skeletal ...", s("osteology", "puberty", "density")),

		c("s", s("Yi Chen, K.", "Kurtiss, M.K.N.", "Forthe, M."), "Neonates: Development", "Medical technology has improved ...", s("infrastructure", "regional cost")),
		c("t", s("Yi Chen, K.", "Kurtiss, M.K.N.", "Forthe, M."), "Neonates: Development", "Medical technology has improved ...", s("osteology", "puberty", "density")),
		c("u", s("Hao, R.W.", "Feng, E.", "Chen, M.", "Kai, S."), "Healthcare facilities", "Medical technology has improved ...", s("infrastructure", "regional cost")),
		c("v", s("Hao, R.W.", "Feng, E.", "Chen, M.", "Kai, S."), "Healthcare facilities", "The development of the skeletal ...", s("infrastructure", "regional cost")),
		c("w", s("Mayer, Jonathan", "Ma, Chengdu", "Argent, J."), "Neonates: Development", "The development of the skeletal ...", s("osteology", "puberty", "density")),
		c("x", s("Mayer, Jonathan", "Ma, Chengdu", "Argent, J."), "Healthcare facilities", "The development of the skeletal ...", s("osteology", "puberty", "density")),
	}
	for _, citation := range citations {
		db.Citations.Insert(citation)
	}

	var results []database.Citation
	db.Citations.Filter(makeSearch(q)).Decode(&results)
	expected := []string{"a", "f", "k", "p", "q", "r"}
	found := make([]bool, len(expected))
	for _, result := range results {
		inExpected := false
		for i, eid := range expected {
			if result.EID == eid {
				found[i] = true
				inExpected = true
				break
			}
		}
		if !inExpected {
			t.Logf("Did not expect to find %v in results:\n%v", result, results)
			t.FailNow()
		}
	}
	for i, inResults := range found {
		if !inResults {
			t.Logf("Expected to find %v in results:\n%v", citations[i], results)
			t.FailNow()
		}
	}
}
