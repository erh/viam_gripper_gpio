package viam_gripper_gpio

import (
	"context"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/gripper"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
	"go.viam.com/utils"
)

var GripperModel = resource.ModelNamespace("erh").WithFamily("viam_gripper_gpio").WithModel("gripper")

type Config struct {
	Board    string
	Pin      string
	OpenHigh bool `json:"open_high"`
}

func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.Board == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "board")
	}

	if cfg.Pin == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "pin")
	}

	return []string{cfg.Board}, nil

}

func init() {
	resource.RegisterComponent(
		gripper.API,
		GripperModel,
		resource.Registration[gripper.Gripper, *Config]{
			Constructor: newGripper,
		})
}

func newGripper(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (gripper.Gripper, error) {
	newConf, err := resource.NativeConfig[*Config](config)
	if err != nil {
		return nil, err
	}

	g := &myGripper{
		name: config.ResourceName(),
		mf:   referenceframe.NewSimpleModel("foo"),
		conf: newConf,
	}

	b, err := board.FromDependencies(deps, newConf.Board)
	if err != nil {
		return nil, err
	}

	g.pin, err = b.GPIOPinByName(newConf.Pin)
	if err != nil {
		return nil, err
	}

	return g, nil
}

type myGripper struct {
	resource.AlwaysRebuild

	name resource.Name
	mf   referenceframe.Model

	conf *Config

	pin board.GPIOPin
}

func (g *myGripper) Grab(ctx context.Context, extra map[string]interface{}) (bool, error) {
	return false, g.pin.Set(ctx, !g.conf.OpenHigh, extra)
}

func (g *myGripper) Open(ctx context.Context, extra map[string]interface{}) error {
	return g.pin.Set(ctx, g.conf.OpenHigh, extra)
}

func (g *myGripper) Name() resource.Name {
	return g.name
}

func (g *myGripper) Close(ctx context.Context) error {
	return nil
}

func (g *myGripper) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (g *myGripper) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (g *myGripper) Stop(context.Context, map[string]interface{}) error {
	return nil
}

func (g *myGripper) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return []spatialmath.Geometry{}, nil
}

func (g *myGripper) ModelFrame() referenceframe.Model {
	return g.mf
}
