package viam_gripper_gpio

import (
	"context"
	"fmt"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/gripper"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
	"go.viam.com/utils"
)

var GripperModel = resource.ModelNamespace("erh").WithFamily("viam_gripper_gpio").WithModel("gripper")

type GripperConfig struct {
	Board    string
	Pin      string
	OpenHigh bool `json:"open_high"`
}

func (cfg *GripperConfig) Validate(path string) ([]string, []string, error) {
	if cfg.Board == "" {
		return nil, nil, utils.NewConfigValidationFieldRequiredError(path, "board")
	}

	if cfg.Pin == "" {
		return nil, nil, utils.NewConfigValidationFieldRequiredError(path, "pin")
	}

	return []string{cfg.Board}, nil, nil

}

func init() {
	resource.RegisterComponent(
		gripper.API,
		GripperModel,
		resource.Registration[gripper.Gripper, *GripperConfig]{
			Constructor: newGripper,
		})
}

func newGripper(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (gripper.Gripper, error) {
	newConf, err := resource.NativeConfig[*GripperConfig](config)
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

	conf *GripperConfig

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

func (g *myGripper) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	return []referenceframe.Input{}, nil
}

func (g *myGripper) GoToInputs(ctx context.Context, inputs ...[]referenceframe.Input) error {
	return fmt.Errorf("GoToInputs not implemented")
}

func (g *myGripper) IsHoldingSomething(ctx context.Context, extra map[string]interface{}) (gripper.HoldingStatus, error) {
	return gripper.HoldingStatus{}, fmt.Errorf("IsHoldingSomething not implemented")
}

func (g *myGripper) Kinematics(ctx context.Context) (referenceframe.Model, error) {
	return g.mf, nil
}
