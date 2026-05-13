package utils

import "math"

// WilsonScore 计算 Wilson 区间下界
// 公式: score = (p + z²/2n - z*sqrt((p*(1-p)+z²/4n)/n)) / (1+z²/n)
// p = 正点赞率 (like/view), n = 总数 (view), z = 置信度参数
func WilsonScore(likes, views int64, z float64) float64 {
	if views == 0 {
		return 0
	}
	p := float64(likes) / float64(views)
	n := float64(views)
	z2 := z * z

	denominator := 1 + z2/n
	numerator := p + z2/(2*n) - z*math.Sqrt((p*(1-p)+z2/(4*n))/n)

	score := numerator / denominator
	if score < 0 {
		return 0
	}
	return score
}
