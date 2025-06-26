# viam-gripper-gpio

https://app.viam.com/module/erh/gripper-gpio

# gripper
GPIO Controlled Gripper where open close is high or low
```
{
  "board": "local",
  "pin": "37",
  "open_high" : <bool> // optional, default-false; false means open is low
}
```

# gripper-press
GPIO Controlled Gripper where it holds down gpio to open or close

For single pin support:
```
{
  "board": "local",
  "pin": "37",
  "seconds" : 3// optional
}
```

For multi pin support:
```
{
  "board": "local",
  "grab_pins": {
    "16": "high"
  },
  "open_pins": {
    "18": "low""
  },
  "wait_pins": {
    "11": "high"
  },
  "grab_time_ms": 0, // optional; 0 means no timeout
  "open_time_ms": 500, // optional; 0 means no timeout
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
