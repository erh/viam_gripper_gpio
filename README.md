# viam-gripper-gpio

https://app.viam.com/module/erh/gripper-gpio

# gripper
GPIO Controlled Gripper where open close is high or low

For single pin support:
```
{
  "board": "local",
  "pin": "37",
  "open_high" : <bool> // optional, default-false; false means open is low
  "geometries" : [ { "type" : "box", "x" : 100, "y": 100, "z" : 100 } ] <optional>

}
```

For multi pin support:
```
{
  "board": "local",
  "grab_pins": {
    "11": "high",
    "16": "high",
    "18": "low"
  },
  "open_pins": {
    "11": "high",
    "16": "low",
    "18": "high"
  }
}
```

# gripper-press
GPIO Controlled Gripper where it holds down gpio to open or close
```
{
  "board": "local",
  "pin": "37",
  "seconds" : 3// optional
  "geometries" : [ { "type" : "box", "x" : 100, "y": 100, "z" : 100 } ] <optional>

}
```

# button
Push turns gpio for seconds
```
{
  "board": "local",
  "pin": "37",
  "seconds" : 1 // optional
}
```

# switch
Switch for gpio with 2 positions
```
{
  "board": "local",
  "pin": "37",
}
```
