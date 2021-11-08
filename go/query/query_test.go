package query

import (
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func s(str string) []string {
	return strings.Split(str, ",")
}

func TestQuery(t *testing.T) {
	// ( OR
	//   ( AND KEYWORD(circulat)
	//         ( OR KEYWORD(pulmon) KEYWORD(cardi) )
	//         ( OR AUTHOR(hayes) AUTHOR(beth) )
	//   )
	//   ( OR
	//     ( AND AUTHOR(kur) AUTHOR(for) KEYWORD(birth) )
	//     ( AND AUTHOR(feng) KEYWORD(aging) )
	//     ( AND AUTHOR(mayer) KEYWORD(grow) KEYWORD(matur) )
	//   )
	// )
	q := bson.M{
		"or": []bson.M{
			{
				"keyword": s("circulat"),
				"or": []bson.M{
					{"keyword": "pulmon"},
					{"keyword": "cardi"},
				},
			},
			{
				"or": []bson.M{
					{"faculty": s("kur,for"), "keyword": s("birth")},
					{"faculty": s("feng"), "keyword": s("aging")},
					{"faculty": s("mayer"), "keyword": s("grow,matur")},
				},
			},
		},
	}
	t.Log(q)
}
