package utils


import (
	"time"
	"github.com/haoel/downsampling/core"
	"github.com/opengovern/opensecurity/services/core/api"
	"github.com/opengovern/opensecurity/services/core/db"
	"github.com/opengovern/opensecurity/services/core/db/models"

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



const (
	ConfigMetadataKeyPrefix = "config_metadata:"
)

func GetConfigMetadata(db db.Database, key string) (models.IConfigMetadata, error) {
	//value, err := rdb.Get(ConfigMetadataKeyPrefix + key)
	//if err == nil {
	//	var cm models.ConfigMetadata
	//	err := json.Unmarshal([]byte(value), &cm)
	//	if err != nil {
	//		return nil, err
	//	}
	//	typedCm, err := cm.ParseToType()
	//	if err != nil {
	//		return nil, err
	//	}
	//	return typedCm, nil
	//} else if err != redis.Nil {
	//	fmt.Printf("error getting config metadata from redis: %v\n", err)
	//}
	//
	typedCm, err := db.GetConfigMetadata(key)
	if err != nil {
		return nil, err
	}
	//jsonCm, err := json.Marshal(typedCm.GetCore())
	//if err != nil {
	//	fmt.Printf("error marshalling config metadata: %v\n", err)
	//	return typedCm, nil
	//}
	//
	//err = rdb.Set(ConfigMetadataKeyPrefix+key, string(jsonCm))
	//if err != nil {
	//	fmt.Printf("error setting config metadata in redis: %v\n", err)
	//	return typedCm, nil
	//}
	//
	return typedCm, nil
}

func SetConfigMetadata(db db.Database, key models.MetadataKey, value any) error {
	valueStr, err := key.GetConfigMetadataType().SerializeValue(value)
	if err != nil {
		return err
	}
	err = db.SetConfigMetadata(models.ConfigMetadata{
		Key:   key,
		Type:  key.GetConfigMetadataType(),
		Value: valueStr,
	})
	if err != nil {
		return err
	}
	return nil
}


var categoryMap = map[string][]string{
	"Identity & Access": []string{
		"aws::iam::user",
		"aws::iam::group",
		"aws::iam::policy",
		"aws::iam::policyattachment",
		"aws::iam::role",
		"aws::iam::accessadvisor",
		"aws::iam::accountpasswordpolicy",
		"aws::identitystore::groupmembership",
		"aws::identitystore::user",
		"aws::identitystore::group",
		"aws::ssoadmin::accountassignment",
		"aws::ssoadmin::permissionset",
		"aws::ssoadmin::attachedmanagedpolicy",
		"aws::ssoadmin::instance",
		"microsoft.authorization/roleassignment",
		"microsoft.authorization/policyassignments",
		"microsoft.authorization/roledefinitions",
	},
	"Entra ID Directory": []string{
		"microsoft.entra/users",
		"microsoft.entra/directoryauditreport",
		"microsoft.entra/userregistrationdetails",
		"microsoft.entra/groups",
	},
}


