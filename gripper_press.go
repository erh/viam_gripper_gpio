package viam_gripper_gpio

import (
	"context"
	"errors"
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
	Pin     string // The pin to use for the gripper, if not using grab_pins or open_pins
	Seconds *int
	GrabPins map[string]bool `json:"grab_pins"`
	OpenPins map[string]bool `json:"open_pins"`
	WaitPins map[string]bool `json:"wait_pins,omitempty"`
	OpenTime *int             `json:"open_time_ms,omitempty"`
	GrabTime *int             `json:"grab_time_ms,omitempty"`
}

func (cfg *ConfigPress) Validate(path string) ([]string, error) {
	if cfg.Board == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "board")
	}

	if cfg.Pin == "" && (cfg.GrabPins == nil && cfg.OpenPins == nil) {
		return nil, utils.NewConfigValidationError(path, errors.New("either pin or grab_pins and open_pins must be specified"))
	}

	if cfg.Pin != "" && (len(cfg.GrabPins) > 0 || len(cfg.OpenPins) > 0 || len(cfg.WaitPins) > 0) {
		return nil, utils.NewConfigValidationError(path, errors.New("pin cannot be used with grab_pins, open_pins, or wait_pins"))
	}

	if cfg.Pin == "" && len(cfg.GrabPins) == 0 {
		return nil, utils.NewConfigValidationError(path, errors.New("grab_pins must not be empty"))
	}

	if cfg.Pin == "" && len(cfg.OpenPins) == 0 {
		return nil, utils.NewConfigValidationError(path, errors.New("open_pins must not be empty"))
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
		mf:   referenceframe.NewSimpleModel(config.ResourceName().String()),
		conf: newConf,
	}

	if g.conf.Seconds == nil {
		defaultSeconds := 3
		g.conf.Seconds = &defaultSeconds
	}

	defaultSecondsMs := 3000
	if g.conf.GrabTime == nil {
		g.conf.GrabTime = &defaultSecondsMs
	}

	if g.conf.OpenTime == nil {
		g.conf.OpenTime = &defaultSecondsMs
	}

	if g.conf.Pin != "" {
		g.pins = make(map[string]bool)
		g.pins[g.conf.Pin] = true
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

	pins map[string]bool
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

	// If the "Pin" field is set, only use that pin to control grab/open
	if len(g.pins) != 0 {
		err := g.setPins(ctx, g.pins, true, extra)
		if err != nil {
			return false, err
		}
		g.open = false
		// Return early if no grab time is specified
		if *g.conf.Seconds == 0 {
			return false, nil
		}
		time.Sleep(time.Second * time.Duration(*g.conf.Seconds))
		err = g.setPins(ctx, g.pins, false, extra)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Set any wait pins to their active state
	if len(g.conf.WaitPins) > 0 {
		err := g.setPins(ctx, g.conf.WaitPins, true, extra)
		if err != nil {
			return false, err
		}
	}

	// Set open pins to their inactive state
	err := g.setPins(ctx, g.conf.OpenPins, false, extra)
	if err != nil {
		return false, err
	}

	// Set the grab pins to their active state
	err = g.setPins(ctx, g.conf.GrabPins, true, extra)
	if err != nil {
		return false, err
	}
	g.open = false

	// Return early if no grab time is specified
	if *g.conf.GrabTime == 0 || *g.conf.Seconds == 0 {
		return false, nil
	}
	time.Sleep(time.Millisecond * time.Duration(*g.conf.GrabTime))

	// Set the grab pins to their inactive state
	err = g.setPins(ctx, g.conf.GrabPins, false, extra)
	if err != nil {
		return false, err
	}

	// Set any wait pins to their inactive state
	if len(g.conf.WaitPins) > 0 {
		err = g.setPins(ctx, g.conf.WaitPins, false, extra)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (g *myGripperPress) setPins(ctx context.Context, pins map[string]bool, activate bool, extra map[string]interface{}) error {
	for pinName, state := range pins {
		pin, err := g.board.GPIOPinByName(pinName)
		if err != nil {
			return err
		}
		if !activate {
			state = !state
		}
		err = pin.Set(ctx, state, extra)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *myGripperPress) Open(ctx context.Context, extra map[string]interface{}) error {
	if !force(extra) && g.open {
		return nil
	}

	// If the "Pin" field is set, only use that pin to control grab/open
	if len(g.pins) != 0 {
		err := g.setPins(ctx, g.pins, true, extra)
		if err != nil {
			return err
		}
		g.open = false
		// Return early if no grab time is specified
		if *g.conf.Seconds == 0 {
			return nil
		}
		time.Sleep(time.Second * time.Duration(*g.conf.Seconds))
		err = g.setPins(ctx, g.pins, false, extra)
		if err != nil {
			return err
		}
		return nil
	}

	// Set any wait pins to their active state
	if len(g.conf.WaitPins) > 0 {
		err := g.setPins(ctx, g.conf.WaitPins, true, extra)
		if err != nil {
			return err
		}
	}

	// Set grab pins to their inactive state
	err := g.setPins(ctx, g.conf.GrabPins, false, extra)
	if err != nil {
		return err
	}

	// Set the open pins to their active state
	err = g.setPins(ctx, g.conf.OpenPins, true, extra)
	if err != nil {
		return err
	}
	g.open = true

	// Return early if no open time is specified
	if *g.conf.OpenTime == 0 || *g.conf.Seconds == 0 {
		return nil
	}
	time.Sleep(time.Millisecond * time.Duration(*g.conf.OpenTime))

	// Set the open pins to their inactive state
	err = g.setPins(ctx, g.conf.OpenPins, false, extra)
	if err != nil {
		return err
	}

	// Set any wait pins to their inactive state
	if len(g.conf.WaitPins) > 0 {
		err = g.setPins(ctx, g.conf.WaitPins, false, extra)
		if err != nil {
			return err
		}
	}

	return nil
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
