package utils
import _"fmt"

// months
var months = map[string]int{
	"January"   : 1,
	"February"  : 2,
	"March"     : 3,
	"April"     : 4,
	"May"       : 5,
	"June"      : 6,
	"July"      : 7,
	"August"    : 8,
	"September" : 9,
	"October"   : 10,
	"November"  : 11,
	"December"  : 12,
}


func MonthMapping(month string) int {
	return months[month]
}

func Swap (a,b *int) {
	*a, *b = *b, *a
}