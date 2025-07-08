package viam_gripper_gpio

import (
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

var family = resource.ModelNamespace("erh").WithFamily("viam_gripper_gpio")

func ParseGeometries(gcs []spatialmath.GeometryConfig) ([]spatialmath.Geometry, error) {
	gs := []spatialmath.Geometry{}

	for _, gc := range gcs {
		g, err := gc.ParseConfig()
		if err != nil {
			return nil, err
		}
		gs = append(gs, g)
	}

	return gs, nil
}
