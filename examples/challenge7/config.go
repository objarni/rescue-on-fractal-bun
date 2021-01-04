package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Config struct {
	Gravity                 float64
	SpeedX                  float64
	StartX                  float64
	StartY                  float64
	GubbeMaxVelocity        float64
	GubbeKickTicks          int
	GubbeWalkAnimTickSwitch int
	GubbeAcceleration       float64
}

func TryReadCfgFrom(filename string, defaultCfg Config) Config {
	fmt.Println("Reading ", filename, " at ", time.Now().Format(time.Stamp))
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile error, defaulting")
		return defaultCfg
	}
	var cfg = defaultCfg
	err = json.Unmarshal(byteArray, &cfg)
	if err != nil {
		fmt.Println("JSON parse error, defaulting. JSON was: " + string(byteArray))
		return defaultCfg
	}
	return cfg
}
