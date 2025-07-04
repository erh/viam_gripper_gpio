# viam-gripper-gpio

https://app.viam.com/module/erh/gripper-gpio

# gripper
GPIO Controlled Gripper where open close is high or low
```
{
  "board": "local",
  "pin": "37",
  "open_high" : <bool> // optional, default-false; false means open is low
  "geometries" : [ { "type" : "box", "x" : 100, "y": 100, "z" : 100 } ] <optional>

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
