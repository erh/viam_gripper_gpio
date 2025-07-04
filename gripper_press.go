package viam_gripper_gpio

import (
	"context"
	"fmt"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/gripper"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
	"go.viam.com/utils"
)

var GripperPressModel = family.WithModel("gripper-press")

type ConfigPress struct {
	Board   string
	Pin     string
	Seconds int
}

func (cfg *ConfigPress) Validate(path string) ([]string, []string, error) {
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
		GripperPressModel,
		resource.Registration[gripper.Gripper, *ConfigPress]{
			Constructor: newGripperPress,
		})
}

func newGripperPress(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (gripper.Gripper, error) {
	newConf, err := resource.NativeConfig[*ConfigPress](config)
	if err != nil {
		return nil, err
	}

	g := &myGripperPress{
		name: config.ResourceName(),
		mf:   referenceframe.NewSimpleModel("foo"),
		conf: newConf,
	}

	if g.conf.Seconds <= 0 {
		g.conf.Seconds = 3
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

type myGripperPress struct {
	resource.AlwaysRebuild

	name resource.Name
	mf   referenceframe.Model

	conf *ConfigPress

	pin board.GPIOPin

	open bool
}

func force(extra map[string]interface{}) bool {
	if extra == nil {
		return false
	}
	return extra["force"] == true
}

func (g *myGripperPress) Grab(ctx context.Context, extra map[string]interface{}) (bool, error) {
	if !force(extra) && !g.open {
		return false, nil
	}
	g.open = false
	return false, g.press(ctx, extra)
}

func (g *myGripperPress) press(ctx context.Context, extra map[string]interface{}) error {
	err := g.pin.Set(ctx, true, extra)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * time.Duration(g.conf.Seconds))
	return g.pin.Set(ctx, false, extra)
}

func (g *myGripperPress) Open(ctx context.Context, extra map[string]interface{}) error {
	if !force(extra) && g.open {
		return nil
	}
	g.open = true
	return g.press(ctx, extra)
}

func (g *myGripperPress) Name() resource.Name {
	return g.name
}

func (g *myGripperPress) Close(ctx context.Context) error {
	return nil
}

func (g *myGripperPress) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (g *myGripperPress) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (g *myGripperPress) Stop(context.Context, map[string]interface{}) error {
	return nil
}

func (g *myGripperPress) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return []spatialmath.Geometry{}, nil
}

func (g *myGripperPress) ModelFrame() referenceframe.Model {
	return g.mf
}

func (g *myGripperPress) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	return []referenceframe.Input{}, nil
}

func (g *myGripperPress) GoToInputs(ctx context.Context, inputs ...[]referenceframe.Input) error {
	return fmt.Errorf("GoToInputs not implemented")
}

func (g *myGripperPress) IsHoldingSomething(ctx context.Context, extra map[string]interface{}) (gripper.HoldingStatus, error) {
	return gripper.HoldingStatus{}, fmt.Errorf("IsHoldingSomething not implemented")
}

func (g *myGripperPress) Kinematics(ctx context.Context) (referenceframe.Model, error) {
	return g.mf, nil
}
