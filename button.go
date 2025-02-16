package viam_gripper_gpio

import (
	"context"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/button"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
)

var ButtonModel = family.WithModel("button")

type ConfigButton struct {
	Board   string
	Pin     string
	Seconds int
}

func (cfg *ConfigButton) Validate(path string) ([]string, error) {
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
		button.API,
		ButtonModel,
		resource.Registration[button.Button, *ConfigButton]{
			Constructor: newButton,
		})
}

func newButton(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (button.Button, error) {
	newConf, err := resource.NativeConfig[*ConfigButton](config)
	if err != nil {
		return nil, err
	}

	g := &buttonData{
		name: config.ResourceName(),
		conf: newConf,
	}

	if g.conf.Seconds <= 0 {
		g.conf.Seconds = 1
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

type buttonData struct {
	resource.AlwaysRebuild

	name resource.Name

	conf *ConfigButton

	pin board.GPIOPin
}

func (g *buttonData) Push(ctx context.Context, extra map[string]interface{}) error {
	err := g.pin.Set(ctx, true, extra)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * time.Duration(g.conf.Seconds))
	return g.pin.Set(ctx, false, extra)
}

func (g *buttonData) Name() resource.Name {
	return g.name
}

func (g *buttonData) Close(ctx context.Context) error {
	return nil
}

func (g *buttonData) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
