package viam_gripper_gpio

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
)

var SwitchModelOneOf = family.WithModel("switch-one-of")

type ConfigSwitchOneOf struct {
	Board string

	Pins  []string
	Names []string
}

func (cfg *ConfigSwitchOneOf) Validate(path string) ([]string, []string, error) {
	if cfg.Board == "" {
		return nil, nil, utils.NewConfigValidationFieldRequiredError(path, "board")
	}

	if len(cfg.Pins) == 0 {
		return nil, nil, utils.NewConfigValidationFieldRequiredError(path, "pins")
	}

	if len(cfg.Names) != len(cfg.Pins) {
		return nil, nil, fmt.Errorf("pins and names have to be the same length")
	}

	return []string{cfg.Board}, nil, nil
}

func init() {
	resource.RegisterComponent(
		toggleswitch.API,
		SwitchModelOneOf,
		resource.Registration[toggleswitch.Switch, *ConfigSwitchOneOf]{
			Constructor: newSwitchOneOf,
		})
}

func newSwitchOneOf(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (toggleswitch.Switch, error) {
	newConf, err := resource.NativeConfig[*ConfigSwitchOneOf](config)
	if err != nil {
		return nil, err
	}

	g := &switchDataOneOf{
		name: config.ResourceName(),
		conf: newConf,
	}

	b, err := board.FromDependencies(deps, newConf.Board)
	if err != nil {
		return nil, err
	}

	for _, p := range newConf.Pins {
		pp, err := b.GPIOPinByName(p)
		if err != nil {
			return nil, err
		}
		g.pins = append(g.pins, pp)
	}

	return g, nil
}

type switchDataOneOf struct {
	resource.AlwaysRebuild

	name resource.Name
	conf *ConfigSwitchOneOf

	pins []board.GPIOPin

	position uint32
}

func (g *switchDataOneOf) Name() resource.Name {
	return g.name
}

func (g *switchDataOneOf) Close(ctx context.Context) error {
	return nil
}

func (g *switchDataOneOf) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (g *switchDataOneOf) SetPosition(ctx context.Context, position uint32, extra map[string]interface{}) error {
	if int(position) < 0 || int(position) > len(g.pins) {
		return fmt.Errorf("SetPosition wrong %d", position)
	}

	g.position = position

	var err error
	for idx, p := range g.pins {
		if (idx + 1) == int(position) {
			err = multierr.Combine(err, p.Set(ctx, true, extra))
		} else {
			err = multierr.Combine(err, p.Set(ctx, false, extra))
		}
	}
	return err
}

func (g *switchDataOneOf) GetPosition(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return g.position, nil
}

func (g *switchDataOneOf) GetNumberOfPositions(ctx context.Context, extra map[string]interface{}) (uint32, []string, error) {
	x := []string{"off"}
	for _, n := range g.conf.Names {
		x = append(x, n)
	}
	return uint32(len(x)), x, nil
}
