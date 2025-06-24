package viam_gripper_gpio

import (
	"context"
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
	Board    string
	GrabPins map[string]bool `json:"grab_pins"`
	OpenPins map[string]bool `json:"open_pins"`
	WaitPins map[string]bool `json:"wait_pins,omitempty"`
	OpenTime int             `json:"open_time_ms,omitempty"`
	GrabTime int             `json:"grab_time_ms,omitempty"`
}

func (cfg *ConfigPress) Validate(path string) ([]string, error) {
	if cfg.Board == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "board")
	}

	if len(cfg.GrabPins) == 0 {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "grab_pins")
	}

	if len(cfg.OpenPins) == 0 {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "open_pins")
	}

	return []string{cfg.Board}, nil
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

	if g.conf.GrabTime <= 0 {
		g.conf.GrabTime = 3000
	}

	if g.conf.OpenTime <= 0 {
		g.conf.OpenTime = 3000
	}

	g.board, err = board.FromDependencies(deps, newConf.Board)
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

	board board.Board

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
	return false, g.press(ctx, extra, g.conf.GrabPins, g.conf.GrabTime)
}

func (g *myGripperPress) press(ctx context.Context, extra map[string]interface{}, pins map[string]bool, duration int) error {
	if len(g.conf.WaitPins) > 0 {
		for waitPinName, waitState := range g.conf.WaitPins {
			waitPin, err := g.board.GPIOPinByName(waitPinName)
			if err != nil {
				return err
			}
			err = waitPin.Set(ctx, waitState, extra)
			if err != nil {
				return err
			}
		}
	}
	for pinName, state := range pins {
		pin, err := g.board.GPIOPinByName(pinName)
		if err != nil {
			return err
		}
		err = pin.Set(ctx, state, extra)
		if err != nil {
			return err
		}
	}
	if duration == 0 {
		return nil
	}
	time.Sleep(time.Millisecond * time.Duration(duration))
	for pinName, state := range pins {
		pin, err := g.board.GPIOPinByName(pinName)
		if err != nil {
			return err
		}
		err = pin.Set(ctx, !state, extra)
		if err != nil {
			return err
		}
	}
	if len(g.conf.WaitPins) > 0 {
		for waitPinName, waitState := range g.conf.WaitPins {
			waitPin, err := g.board.GPIOPinByName(waitPinName)
			if err != nil {
				return err
			}
			err = waitPin.Set(ctx, !waitState, extra)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *myGripperPress) Open(ctx context.Context, extra map[string]interface{}) error {
	if !force(extra) && g.open {
		return nil
	}
	g.open = true
	return g.press(ctx, extra, g.conf.OpenPins, g.conf.OpenTime)
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
