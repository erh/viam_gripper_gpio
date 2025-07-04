package viam_gripper_gpio

import (
	"context"
	"fmt"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
)

var SwitchModel = family.WithModel("switch")

type ConfigSwitch struct {
	Board string
	Pin   string
}

func (cfg *ConfigSwitch) Validate(path string) ([]string, []string, error) {
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
		toggleswitch.API,
		SwitchModel,
		resource.Registration[toggleswitch.Switch, *ConfigSwitch]{
			Constructor: newSwitch,
		})
}

func newSwitch(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (toggleswitch.Switch, error) {
	newConf, err := resource.NativeConfig[*ConfigSwitch](config)
	if err != nil {
		return nil, err
	}

	g := &switchData{
		name: config.ResourceName(),
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

type switchData struct {
	resource.AlwaysRebuild

	name resource.Name

	conf *ConfigSwitch

	pin board.GPIOPin

	position uint32
}

func (g *switchData) Name() resource.Name {
	return g.name
}

func (g *switchData) Close(ctx context.Context) error {
	return nil
}

func (g *switchData) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (g *switchData) SetPosition(ctx context.Context, position uint32, extra map[string]interface{}) error {
	if position > 1 {
		return fmt.Errorf("gpio SetPosition only support 0 and 1, not %d", position)
	}

	g.position = position
	if position == 0 {
		return g.pin.Set(ctx, false, extra)
	}
	return g.pin.Set(ctx, true, extra)
}

func (g *switchData) GetPosition(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return g.position, nil
}

func (g *switchData) GetNumberOfPositions(ctx context.Context, extra map[string]interface{}) (uint32, []string, error) {
	return 2, []string{"off", "on"}, nil
}
