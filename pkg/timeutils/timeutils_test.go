package timeutils_test

import (
	"fmt"
	. "mahjong/pkg/timeutils"
	"testing"
)

func TestBeginOfNow(t *testing.T) {
	begin := BeginOfNow()
	fmt.Println(begin)
}
