package main

import (
	"go.viam.com/rdk/components/button"
	"go.viam.com/rdk/components/gripper"
	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"

	"github.com/erh/viam_gripper_gpio"
)

func main() {
	module.ModularMain(
		resource.APIModel{gripper.API, viam_gripper_gpio.GripperModel},
		resource.APIModel{gripper.API, viam_gripper_gpio.GripperPressModel},
		resource.APIModel{button.API, viam_gripper_gpio.ButtonModel},
		resource.APIModel{toggleswitch.API, viam_gripper_gpio.SwitchModel},
		resource.APIModel{toggleswitch.API, viam_gripper_gpio.SwitchModelOneOf},
	)
}
