package tests

import (
	"fmt"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
)

func Example_initialRender() {
	cfg := scenes.TryReadCfgFrom("../"+internal.ConfigFile, scenes.Config{})
	res := internal.Resources{
		ImageMap: map[internal.Image]*pixel.Sprite{},
	}
	mapScene := scenes.MakeMapScene(&cfg, &res, "Skogen")
	op := mapScene.MapSceneWinOp()
	fmt.Print(op.String())
	// Output:
	/*
	WinOp Sequence:
	  Moved 400 pixels right 300 pixels up:
	    Image "IMap"
	  WinOp from ImdOp:
	    Sequence:
	      Sequence:
	  Color 72, 61, 139:
	  Circle radius 5 center <246, 109> thickness 3
	  Color 72, 61, 139:
	  Circle radius 5 center <355, 235> thickness 3
	  Color 72, 61, 139:
	  Circle radius 5 center <299, 375> thickness 3
	  Color 0, 128, 0:
	  Circle radius 15 center <299, 375> thickness 3
	  Color 255, 0, 0:
	  Circle radius 11 center <299, 375> thickness 3
	      Color 255, 0, 0:
	  Sequence:
	    Line from <299, 0> to <299, 600> thickness 2
	    Line from <0, 375> to <800, 375> thickness 2
	*/
}
