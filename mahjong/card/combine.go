package card

type CardList []*Card

func (cl CardList) Len() int {
	return len(cl)
}

func (cl CardList) Less(i, j int) bool {
	if Value(cl[i]) < Value(cl[j]) {
		return true
	}
	return false
}

func (cl CardList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

//1->next
//11 ->next
//111 -> next
//1111 -> next
//12 ->next
//122 ->next
//1222 ->next
//12222 ->next

//13 ->take,reset
//123 ->take,reset
//124 ->take,reset
//112 ->take,reset
//1112 ->1.try duizi 2.try kezi
//11112 ->1.try kezi 3,try gangzi
//1223  ->take,reset
//1224 ->take,reset
//12223 ->take,reset
//12224 ->take,reset
//122223 ->take,reset
//122224 ->take,reset

//must sort before
func Combine(cards CardList) (shunziList CardList, gangList CardList, keziList CardList, duiziList CardList, remainList CardList) {
	if len(cards) < 2 {
		return nil, nil, nil, nil, append(remainList, cards...)
	}
	var (
		tempCard *Card
	)
	var (
		numEquals   = 0
		numAdjacent = 0
	)
	for _, c := range cards {
		if tempCard == nil {
			goto RESET
		}

		if tempCard.Equal(c) {
			switch numAdjacent {
			//1,11,111,1111,122,1222,12222
			case 1, 2:
				numEquals += 1
				goto Next
			default:
				panic("never reach here")
			}
		}

		if tempCard.Adjacent(c) {
			switch numEquals {
			//12
			case 1:
				numAdjacent += 1
				if numAdjacent == 3 {
					shunzi := cards[0:1]
					var remains CardList
					remains = append(remains, cards[3:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, shunzi...)
					shunziList = append(shunziList, tShunzi...)
					gangList = append(gangList, tGang...)
					keziList = append(keziList, tKezi...)
					duiziList = append(duiziList, tDuizi...)
					remainList = append(remainList, tRemains...)
					return
				}
				goto Next
			case 2, 3, 4:
				{
					//112,1112,11112
					if numAdjacent == 1 {
						//try duizi
						duizi := cards[0:1]
						var duiziRemains CardList
						duiziRemains = append(duiziRemains, cards[2:]...)

						tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(duiziRemains)
						if numEquals == 2 || len(tRemains) == 0 {
							shunziList = append(shunziList, tShunzi...)
							gangList = append(gangList, tGang...)
							keziList = append(keziList, tKezi...)
							duiziList = append(duiziList, duizi...)
							duiziList = append(duiziList, tDuizi...)
							remainList = append(remainList, tRemains...)
							return
						}

						//try kezi
						kezi := cards[0:1]
						var keziRemains CardList
						keziRemains = append(keziRemains, cards[3:]...)

						tShunzi, tGang, tKezi, tDuizi, tRemains = Combine(keziRemains)
						if numEquals == 3 || len(tRemains) == 0 {
							shunziList = append(shunziList, tShunzi...)
							gangList = append(gangList, tGang...)
							keziList = append(keziList, kezi...)
							keziList = append(keziList, tKezi...)

							duiziList = append(duiziList, tDuizi...)
							remainList = append(remainList, tRemains...)
							return
						}

						//try 2duizi
						gang := cards[0:1]
						var gangRemains CardList
						gangRemains = append(gangRemains, cards[4:]...)

						tShunzi, tGang, tKezi, tDuizi, tRemains = Combine(gangRemains)
						shunziList = append(shunziList, tShunzi...)
						gangList = append(gangList, gang...)
						gangList = append(gangList, tGang...)
						keziList = append(keziList, tKezi...)

						duiziList = append(duiziList, tDuizi...)
						remainList = append(remainList, tRemains...)
						return

					}

					//1223,12223,122223
					if numAdjacent == 2 {
						//try shunzi
						shunzi := append(cards[0:1])
						var remains CardList
						remains = append(remains, cards[2:numEquals+1]...)
						remains = append(remains, cards[numEquals+2:]...)

						tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)

						shunziList = append(shunziList, shunzi...)
						shunziList = append(shunziList, tShunzi...)
						gangList = append(gangList, tGang...)
						keziList = append(keziList, tKezi...)

						duiziList = append(duiziList, tDuizi...)
						remainList = append(remainList, tRemains...)
						return
					}
					panic("never reach here")
				}
			default:
				panic("never reach here")
			}

		}

		switch numEquals {
		case 1:
			{
				//124
				if numAdjacent == 2 {
					rem := cards[0:2]
					var remains CardList
					remains = append(remains, cards[2:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)

					shunziList = append(shunziList, tShunzi...)

					keziList = append(keziList, tKezi...)
					gangList = append(gangList, tGang...)
					duiziList = append(duiziList, tDuizi...)
					remainList = append(remainList, rem...)
					remainList = append(remainList, tRemains...)
					return
				}
				//13
				if numAdjacent == 1 {
					rem := cards[0:1]
					var remains CardList
					remains = append(remains, cards[1:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, tShunzi...)
					keziList = append(keziList, tKezi...)
					gangList = append(gangList, tGang...)
					duiziList = append(duiziList, tDuizi...)
					remainList = append(remainList, rem...)
					remainList = append(remainList, tRemains...)
					return
				}
				panic("never reache here")
			}

		case 2:
			{

				//1224
				if numAdjacent == 2 {
					rem := cards[0:1]
					duizi := cards[2:3]

					var remains CardList
					remains = append(remains, cards[3:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, tShunzi...)
					gangList = append(gangList, tGang...)
					keziList = append(keziList, tKezi...)

					duiziList = append(duiziList, duizi...)
					duiziList = append(duiziList, tDuizi...)
					remainList = append(remainList, rem...)
					remainList = append(remainList, tRemains...)
					return
				}
				//113
				if numAdjacent == 1 {
					duizi := cards[0:1]

					var remains CardList
					remains = append(remains, cards[2:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, tShunzi...)
					gangList = append(gangList, tGang...)
					keziList = append(keziList, tKezi...)
					duiziList = append(duiziList, duizi...)
					duiziList = append(duiziList, tDuizi...)

					remainList = append(remainList, tRemains...)
					return
				}
				panic("never reache here")
			}

		case 3:
			{

				//12224
				if numAdjacent == 2 {
					rem := cards[0:1]
					kezi := cards[1:2]

					var remains CardList
					remains = append(remains, cards[4:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, tShunzi...)
					gangList = append(gangList, tGang...)
					keziList = append(keziList, kezi...)
					keziList = append(keziList, tKezi...)

					duiziList = append(duiziList, tDuizi...)
					remainList = append(remainList, rem...)
					remainList = append(remainList, tRemains...)
					return
				}
				//1113
				if numAdjacent == 1 {
					kezi := cards[0:1]

					var remains CardList
					remains = append(remains, cards[3:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, tShunzi...)
					gangList = append(gangList, tGang...)
					keziList = append(keziList, kezi...)
					keziList = append(keziList, tKezi...)

					duiziList = append(duiziList, tDuizi...)

					remainList = append(remainList, tRemains...)
					return
				}
				panic("never reache here")
			}

		//todo 122224
		case 4:
			{
				//122224
				if numAdjacent == 2 {
					rem := cards[0:1]
					gang := cards[1:2]

					var remains CardList
					remains = append(remains, cards[5:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, tShunzi...)
					gangList = append(gangList, gang...)
					gangList = append(gangList, tGang...)
					keziList = append(keziList, tKezi...)

					duiziList = append(duiziList, tDuizi...)
					remainList = append(remainList, rem...)
					remainList = append(remainList, tRemains...)
					return
				}
				//11113
				if numAdjacent == 1 {
					gang := cards[0:1]

					var remains CardList
					remains = append(remains, cards[4:]...)

					tShunzi, tGang, tKezi, tDuizi, tRemains := Combine(remains)
					shunziList = append(shunziList, tShunzi...)
					gangList = append(gangList, gang...)
					gangList = append(gangList, tGang...)
					keziList = append(keziList, tKezi...)

					duiziList = append(duiziList, tDuizi...)

					remainList = append(remainList, tRemains...)
					return
				}
				panic("never reache here")
			}
		default:
			panic("never reach here")
		}

	Next:
		tempCard = c
		continue
	RESET:
		tempCard = c
		numAdjacent = 1
		numEquals = 1
	}

	switch numEquals {
	//1,12
	case 1:
		{
			remainList = append(remainList, cards...)
			return
		}
		//11,122
	case 2:
		{
			if numAdjacent == 2 {
				//TODO
				duizi := cards[1:2]
				remains := cards[0:1]
				duiziList = append(duiziList, duizi...)

				remainList = append(remainList, remains...)
				return
			}
			//11
			if numAdjacent == 1 {
				duizi := cards[0:1]
				duiziList = append(duiziList, duizi...)

				return
			}
		}
		//111,1222
	case 3:
		{
			if numAdjacent == 2 {
				//TODO
				kezi := cards[1:2]
				remains := cards[0:1]
				keziList = append(keziList, kezi...)
				remainList = append(remainList, remains...)
				return
			}
			//111
			if numAdjacent == 1 {
				kezi := cards[0:1]
				keziList = append(keziList, kezi...)
				return
			}
		}
		//1111,12222
	case 4:
		{
			if numAdjacent == 2 {
				//TODO
				duizi := cards[1:3]
				remains := cards[0:1]
				duiziList = append(duiziList, duizi...)

				remainList = append(remainList, remains...)
				return
			}
			//11
			if numAdjacent == 1 {
				duizi := cards[0:2]
				duiziList = append(duiziList, duizi...)
				return
			}
		}
	default:
		panic("never reach here")
	}
	return
}

func IsStraight(cards CardList) bool {
	var tc *Card
	for _, c := range cards {
		if tc == nil {
			tc = c
			continue
		}
		if !tc.Adjacent(c) {
			return false
		}
		tc = c
	}
	return true
}

func Count(cards CardList) (gangList CardList, keziList CardList, duiziList CardList, remainList CardList) {
	var tc *Card
	numEqual := 0
	for _, c := range cards {
		if tc == nil {
			goto Reset
		}

		if c.Equal(tc) {
			numEqual++
			continue
		}
		switch numEqual {
		case 1:
			remainList = append(remainList, tc)
		case 2:
			duiziList = append(remainList, tc)
		case 3:
			keziList = append(remainList, tc)
		case 4:
			gangList = append(gangList, tc)
		default:
			panic("never reach here")
		}
	Reset:
		tc = c
		numEqual = 1
	}

	switch numEqual {
	case 1:
		remainList = append(remainList, tc)
	case 2:
		duiziList = append(remainList, tc)
	case 3:
		keziList = append(remainList, tc)
	case 4:
		gangList = append(gangList, tc)
	}

	return
}
