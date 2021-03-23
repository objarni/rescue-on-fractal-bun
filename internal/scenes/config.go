package scenes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	MapSceneBlinkSpeed             int
	MapSceneCrosshairSpeed         float64
	MapSceneLocationTextX          int
	MapSceneLocationTextY          int
	MapSceneLocCircleRadius        int
	MapSceneCurrentLocCircleRadius int
	MapSceneTargetLocCircleRadius  int
	MapSceneTargetLocMaxDistance   int
	LevelSceneMoveSpeed            float64
	LevelSceneEliseFPS             int
}

func TryReadCfgFrom(filename string, defaultCfg Config) Config {
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
