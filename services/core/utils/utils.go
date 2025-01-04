package utils


import (
	"time"

	"github.com/haoel/downsampling/core"
	"github.com/opengovern/opencomply/services/inventory/api"
)

func resourceTypeTrendDataPointsToPoints(trendDataPoints []api.ResourceTypeTrendDatapoint) []core.Point {
	points := make([]core.Point, len(trendDataPoints))
	for i, trendDataPoint := range trendDataPoints {
		points[i] = core.Point{
			X: float64(trendDataPoint.Date.UnixMilli()),
			Y: float64(trendDataPoint.Count),
		}
	}
	return points
}

func pointsToResourceTypeTrendDataPoints(points []core.Point) []api.ResourceTypeTrendDatapoint {
	trendDataPoints := make([]api.ResourceTypeTrendDatapoint, len(points))
	for i, point := range points {
		trendDataPoints[i] = api.ResourceTypeTrendDatapoint{
			Date:  time.UnixMilli(int64(point.X)),
			Count: int(point.Y),
		}
	}
	return trendDataPoints
}

func DownSampleResourceTypeTrendDatapoints(trendDataPoints []api.ResourceTypeTrendDatapoint, maxDataPoints int) []api.ResourceTypeTrendDatapoint {
	if len(trendDataPoints) <= maxDataPoints {
		return trendDataPoints
	}
	downSampledResourceCounts := core.LTTB(resourceTypeTrendDataPointsToPoints(trendDataPoints), maxDataPoints)
	return pointsToResourceTypeTrendDataPoints(downSampledResourceCounts)
}

func costTrendDataPointsToPoints(trendDataPoints []api.CostTrendDatapoint) []core.Point {
	points := make([]core.Point, len(trendDataPoints))
	for i, trendDataPoint := range trendDataPoints {
		points[i] = core.Point{
			X: float64(trendDataPoint.Date.UnixMilli()),
			Y: trendDataPoint.Cost,
		}
	}
	return points
}

func pointsToCostTrendDataPoints(points []core.Point) []api.CostTrendDatapoint {
	trendDataPoints := make([]api.CostTrendDatapoint, len(points))
	for i, point := range points {
		trendDataPoints[i] = api.CostTrendDatapoint{
			Date: time.UnixMilli(int64(point.X)),
			Cost: point.Y,
		}
	}
	return trendDataPoints
}

func DownSampleCostTrendDatapoints(trendDataPoints []api.CostTrendDatapoint, maxDataPoints int) []api.CostTrendDatapoint {
	if len(trendDataPoints) <= maxDataPoints {
		return trendDataPoints
	}
	downSampledResourceCounts := core.LTTB(costTrendDataPointsToPoints(trendDataPoints), maxDataPoints)
	return pointsToCostTrendDataPoints(downSampledResourceCounts)
}
