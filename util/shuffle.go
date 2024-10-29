package util

import "math/rand"

// ShuffleT ：Donald E. Knuth --- 歌单随机播放、洗牌、扫雷等需要随机打乱的场景
func ShuffleT[T any](s []T) {
	for i := len(s) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// ManualShuffle 模拟手洗 --- 偶数张牌（定时&全局随机 缓存若干玩家出手顺序）
func ManualShuffle[T any](s []T, nums int) {
	lens := len(s)
	middle := lens / 2
	temp := make([]T, lens)
	for i := 0; i < nums; i++ {
		for j, k := 0, 0; j < middle; j++ {
			temp[k], temp[k+1] = s[j], s[middle+j]
			k += 2
		}
		s, temp = temp, s
	}
}
