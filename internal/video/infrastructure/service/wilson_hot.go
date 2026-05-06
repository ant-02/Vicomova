package service

import (
	"context"
	"math"

	"vicomova/internal/video/domain/entity"
	"vicomova/internal/video/domain/repository"
)

type WilsonHotAlgorithm struct {
	repo repository.VideoRepository
	z    float64 // 置信度参数，默认1.96(95%)
}

func NewWilsonHotAlgorithm(repo repository.VideoRepository) *WilsonHotAlgorithm {
	return &WilsonHotAlgorithm{
		repo: repo,
		z:    1.96,
	}
}

func (w *WilsonHotAlgorithm) CalculateScore(ctx context.Context, videoID int64) (float64, error) {
	v, err := w.repo.GetByID(ctx, videoID)
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}

	return w.wilsonScore(v.LikeCount, v.ViewCount), nil
}

func (w *WilsonHotAlgorithm) IncrementView(ctx context.Context, videoID int64) error {
	return w.repo.IncrementView(ctx, videoID)
}

func (w *WilsonHotAlgorithm) GetHotVideos(ctx context.Context, limit int) ([]*entity.Video, error) {
	// 获取所有已发布的视频，按 view_count 降序
	videos, _, err := w.repo.ListPublished(ctx, 1, limit*3) // 多取一些用于排序
	if err != nil {
		return nil, err
	}

	// Wilson 排序
	sorted := make([]*entity.Video, len(videos))
	copy(sorted, videos)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			si := w.wilsonScore(sorted[i].LikeCount, sorted[i].ViewCount)
			sj := w.wilsonScore(sorted[j].LikeCount, sorted[j].ViewCount)
			if si < sj {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted, nil
}

// wilsonScore 计算 Wilson 区间下界
// 公式: score = (p + z²/2n - z*sqrt((p(1-p)+z²/4n)/n)) / (1+z²/n)
// p = 正点赞率 (like/view), n = 总数 (view)
func (w *WilsonHotAlgorithm) wilsonScore(likes, views int64) float64 {
	if views == 0 {
		return 0
	}
	p := float64(likes) / float64(views)
	n := float64(views)
	z := w.z
	z2 := z * z

	denominator := 1 + z2/n
	numerator := p + z2/(2*n) - z*math.Sqrt((p*(1-p)+z2/(4*n))/n)

	return numerator / denominator
}
