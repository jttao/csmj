package card_test

import (
	"fmt"
	"testing"

	"game/mahjong/card"
)

const (
	test1Card  = "test1Card"
	test11Card = "test11Card"
	test12Card = "test12Card"
	test13Card = "test13Card"

	test111Card = "test111Card"
	test112Card = "test112Card"

	test122Card = "test122Card"
	test123Card = "test123Card"
	test124Card = "test124Card"

	test133Card = "test133Card"
	test134Card = "test134Card"

	test1111Card = "test1111Card"
	test1112Card = "test1112Card"
	test1113Card = "test1113Card"
	test1122Card = "test1122Card"
	test1123Card = "test1123Card"
	test1124Card = "test1124Card"
	test1222Card = "test1222Card"
	test1223Card = "test1223Card"
	test1224Card = "test1224Card"
	test1233Card = "test1233Card"
	test1234Card = "test1234Card"
	test1235Card = "test1235Card"
	test1244Card = "test1244Card"
	test1245Card = "test1245Card"
	test1333Card = "test1333Card"
	test1334Card = "test1334Card"
	test1344Card = "test1344Card"
	test1345Card = "test1345Card"
	test1346Card = "test1346Card"

	//test = "test"

	// test11112Card = "test11112Card"
	// test11113Card = "test11113Card"
	// test11122Card = "test11122Card"
	// test11123Card = "test11123Card"
	// test11124Card = "test11124Card"
	// test11134Card = "test11134Card"
	// test11135Card = "test11135Card"
	// test11222Card = "test11222Card"
	// test11223Card = "test11223Card"
	// test11224Card = "test11224Card"
	// test12222Card = "test12222Card"
	// test12223Card = "test12223Card"
	// test12224Card = "test12224Card"
	// test12233Card = "test12233Card"
	// test12234Card = "test12234Card"
	// test12235Card = "test12235Card"
	// test12244Card = "test12244Card"
	// test12245Card = "test12245Card"
	// test12333Card = "test12333Card"
	// test12334Card = "test12334Card"
	// test12335Card = "test12335Card"
	// test12344Card = "test12344Card"
	// test12345Card = "test12345Card"
	// test12346Card = "test12346Card"
	// test12355Card = "test12355Card"
	// test12356Card = "test12356Card"
	// test13333Card = "test13333Card"
	// test13334Card = "test13334Card"
	// test13345Card = "test13345Card"
	// test13444Card = "test13444Card"
	// test13445Card = "test13445Card"
	// test13455Card = "test13455Card"
	// test13456Card = "test13456Card"
	// test13457Card = "test13457Card"
)

var (
	testSuit = map[string]card.CardList{
		test1Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
		},
		test11Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
		},
		test12Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
		},
		test13Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test111Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
		},
		test112Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
		},
		test122Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
		},
		test123Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test124Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test133Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test134Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test1111Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
		},
		test1112Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
		},
		test1113Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test1122Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
		},
		test1123Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test1124Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test1222Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
		},
		test1223Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test1224Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test1233Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test1234Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test1235Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueFive),
		},
		test1244Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test1245Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueTwo),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
			card.NewCard(card.CardTypeTiao, card.CardValueFive),
		},
		test1333Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
		},
		test1334Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test1344Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
		},
		test1345Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
			card.NewCard(card.CardTypeTiao, card.CardValueFive),
		},
		test1346Card: card.CardList{
			card.NewCard(card.CardTypeTiao, card.CardValueOne),
			card.NewCard(card.CardTypeTiao, card.CardValueThree),
			card.NewCard(card.CardTypeTiao, card.CardValueFour),
			card.NewCard(card.CardTypeTiao, card.CardValueSix),
		},
		// test: card.CardList{
		// 	card.NewCard(card.CardTypeTiao, card.CardValueTwo),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueFour),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueSix),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueSix),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueSix),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueSeven),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueEight),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueEight),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueEight),
		// 	card.NewCard(card.CardTypeTiao, card.CardValueNine),
		// },
	}
)

var (
	testSuitResult = map[string][5]card.CardList{
		test1Card:  [5]card.CardList{nil, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}},
		test11Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil},
		test12Card: [5]card.CardList{nil, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueTwo)}},
		test13Card: [5]card.CardList{nil, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueThree)}},

		test111Card: [5]card.CardList{nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, nil},
		test112Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo)}},
		test122Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}},
		test123Card: [5]card.CardList{card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, nil, nil, nil},

		test124Card: [5]card.CardList{nil, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueTwo), card.NewCard(card.CardTypeTiao, card.CardValueFour)}},
		test133Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueThree)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}},
		test134Card: [5]card.CardList{nil, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueThree), card.NewCard(card.CardTypeTiao, card.CardValueFour)}},

		test1111Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil},

		test1112Card: [5]card.CardList{nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo)}},

		test1113Card: [5]card.CardList{nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueThree)}},

		test1122Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueTwo)}, nil},
		test1123Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo), card.NewCard(card.CardTypeTiao, card.CardValueThree)}},
		test1124Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo), card.NewCard(card.CardTypeTiao, card.CardValueFour)}},
		test1222Card: [5]card.CardList{nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo)}, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}},
		test1223Card: [5]card.CardList{card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo)}},
		test1224Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueFour)}},
		test1233Card: [5]card.CardList{card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueThree)}},
		test1234Card: [5]card.CardList{card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueFour)}},
		test1235Card: [5]card.CardList{card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueFive)}},
		test1244Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueFour)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueTwo)}},
		test1245Card: [5]card.CardList{nil, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueTwo), card.NewCard(card.CardTypeTiao, card.CardValueFour), card.NewCard(card.CardTypeTiao, card.CardValueFive)}},
		test1333Card: [5]card.CardList{nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueThree)}, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}},
		test1334Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueThree)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueFour)}},
		test1344Card: [5]card.CardList{nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueFour)}, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueThree)}},
		test1345Card: [5]card.CardList{card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueThree)}, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne)}},
		test1346Card: [5]card.CardList{nil, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueThree), card.NewCard(card.CardTypeTiao, card.CardValueFour), card.NewCard(card.CardTypeTiao, card.CardValueSix)}},
		//	test:         [5]card.CardList{card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueTwo)}, nil, nil, nil, card.CardList{card.NewCard(card.CardTypeTiao, card.CardValueOne), card.NewCard(card.CardTypeTiao, card.CardValueThree), card.NewCard(card.CardTypeTiao, card.CardValueFour), card.NewCard(card.CardTypeTiao, card.CardValueSix)}},
	}
)

func TestPartJudge(t *testing.T) {

	for key, testCards := range testSuit {
		shunziCards, gangCards, keziCards, duiziCards, remainCards := card.Combine(testCards)

		testResult := testSuitResult[key]
		if !equalCards(shunziCards, testResult[0]) {
			t.Errorf("test %s failed", key)
			t.FailNow()
		}
		if !equalCards(gangCards, testResult[1]) {
			t.Errorf("test %s failed", key)
			t.FailNow()
		}
		if !equalCards(keziCards, testResult[2]) {
			fmt.Println("sdsd", testResult[2])
			t.Errorf("test %s failed", key)
			t.FailNow()
		}
		if !equalCards(duiziCards, testResult[3]) {
			t.Errorf("test %s failed", key)
			t.FailNow()
		}

		if !equalCards(remainCards, testResult[4]) {

			t.Errorf("test %s failed", key)
			t.FailNow()
		}
	}
}

func equalCards(a card.CardList, b card.CardList) bool {

	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}

	for i, c := range a {
		if !b[i].Equal(c) {
			return false
		}
	}
	return true
}
