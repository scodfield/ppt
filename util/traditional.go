package util

import (
	"fmt"
	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"time"
)

// GongWei: 小六壬宫位
type GongWei struct {
	ID   int
	Name string
	Desc string
}

func (gw GongWei) String() string {
	return fmt.Sprintf("ID: %d, Name: %s, Desc: %s", gw.ID, gw.Name, gw.Desc)
}

// ShiChen: 十二时辰
type ShiChen struct {
	ID   int
	Name string
	Hour int
}

var (
	xlr     map[int]GongWei
	shiChen map[int]ShiChen
)

func init() {
	initXLR()
	initShiChen()
}

func initXLR() {
	xlr = make(map[int]GongWei, 6)
	xlr[1] = GongWei{
		ID:   1,
		Name: "大安",
		Desc: "进展顺利、宜慢不宜快",
	}
	xlr[2] = GongWei{
		ID:   2,
		Name: "留连",
		Desc: "运气平平、事情反复",
	}
	xlr[3] = GongWei{
		ID:   3,
		Name: "速喜",
		Desc: "好运马上降临",
	}
	xlr[4] = GongWei{
		ID:   4,
		Name: "赤口",
		Desc: "提防小人",
	}
	xlr[5] = GongWei{
		ID:   5,
		Name: "小吉",
		Desc: "好事发生、耐心等待",
	}
	xlr[0] = GongWei{
		ID:   6,
		Name: "空亡",
		Desc: "诸事不顺",
	}

}

func initShiChen() {
	shiChen = make(map[int]ShiChen, 24)
	for i := 0; i < 24; i++ {
		switch i {
		case 0, 23:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "子时",
				Hour: 1,
			}
		case 1, 2:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "丑时",
				Hour: 2,
			}
		case 3, 4:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "寅时",
				Hour: 3,
			}
		case 5, 6:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "卯时",
				Hour: 4,
			}
		case 7, 8:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "辰时",
				Hour: 5,
			}
		case 9, 10:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "巳时",
				Hour: 6,
			}
		case 11, 12:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "午时",
				Hour: 7,
			}
		case 13, 14:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "未时",
				Hour: 8,
			}
		case 15, 16:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "申时",
				Hour: 9,
			}
		case 17, 18:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "酉时",
				Hour: 10,
			}
		case 19, 20:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "戌时",
				Hour: 11,
			}
		case 21, 22:
			shiChen[i] = ShiChen{
				ID:   i,
				Name: "亥时",
				Hour: 12,
			}
		}
	}
}

// QZYS: 掐指一算--小六壬
func QZYS(now time.Time) []string {
	lunar := calendar.ByTimestamp(now.Unix()).Lunar
	month, day, hour := lunar.GetMonth(), lunar.GetDay(), now.Hour()
	sh := shiChen[hour]
	if sh.Hour == 1 {
		return []string{fmt.Sprintf("%s不占", sh.Name)}
	}
	first := int(month) % 6
	second := (first + int(day)) % 6
	third := (second + sh.Hour) % 6
	return []string{xlr[first].String(), xlr[second].String(), xlr[third].String()}
}
