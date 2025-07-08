package viam_gripper_gpio

import (
	"context"
	"errors"

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
	GrabPins map[string]string `json:"grab_pins,omitempty"`
	OpenPins map[string]string `json:"open_pins,omitempty"`
}

func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.Board == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "board")
	}

	if cfg.Pin == "" && (cfg.GrabPins == nil || cfg.OpenPins == nil) {
		return nil, utils.NewConfigValidationError(path, errors.New("either pin or grab_pins and open_pins must be specified"))
	}

	if cfg.Pin != "" && (len(cfg.GrabPins) > 0 || len(cfg.OpenPins) > 0) {
		return nil, utils.NewConfigValidationError(path, errors.New("pin cannot be used with grab_pins, open_pins, or wait_pins"))
	}

	if cfg.Pin == "" && len(cfg.GrabPins) == 0 {
		return nil, utils.NewConfigValidationError(path, errors.New("grab_pins must not be empty"))
	}

	if cfg.Pin == "" && len(cfg.OpenPins) == 0 {
		return nil, utils.NewConfigValidationError(path, errors.New("open_pins must not be empty"))
	}

	for _, state := range cfg.GrabPins {
		if state != "high" && state != "low" {
			return nil, utils.NewConfigValidationError(path, errors.New("grab_pins must be 'high' or 'low'"))
		}
	}

	for _, state := range cfg.OpenPins {
		if state != "high" && state != "low" {
			return nil, utils.NewConfigValidationError(path, errors.New("open_pins must be 'high' or 'low'"))
		}
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
		mf:   referenceframe.NewSimpleModel(config.ResourceName().String()),
		conf: newConf,
	}

	g.board, err = board.FromDependencies(deps, newConf.Board)
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
	board board.Board
}

func (g *myGripper) Grab(ctx context.Context, extra map[string]interface{}) (bool, error) {
	if g.conf.Pin != "" {
		return false, g.pin.Set(ctx, !g.conf.OpenHigh, extra)
	}

	for pinName, level := range g.conf.GrabPins {
		pin, err := g.board.GPIOPinByName(pinName)
		if err != nil {
			return false, err
		}
		state := level == "high"
		err = pin.Set(ctx, state, extra)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (g *myGripper) Open(ctx context.Context, extra map[string]interface{}) error {
	if g.conf.Pin != "" {
		return g.pin.Set(ctx, g.conf.OpenHigh, extra)
	}

	for pinName, level := range g.conf.OpenPins {
		pin, err := g.board.GPIOPinByName(pinName)
		if err != nil {
			return err
		}
		state := level == "high"
		err = pin.Set(ctx, state, extra)
		if err != nil {
			return err
		}
	}

	return nil
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
