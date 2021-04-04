package entities

import (
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/tests"
	"testing"
)

func Test_1_Entity_1_EventBox_1_overlap(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{10, 10},
	}})
	entityCanvas.AddEventBox(EventBox{Event: "ACTION", Box: pixel.Rect{
		Min: pixel.Vec{5, 5},
		Max: pixel.Vec{6, 6},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

func Test_1_Entity_1_EventBox_no_overlapping(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{10, 10},
	}})
	entityCanvas.AddEventBox(EventBox{Event: "ACTION", Box: pixel.Rect{
		Min: pixel.Vec{55, 55},
		Max: pixel.Vec{76, 76},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_2_Entities_1_EventBox_1_overlap(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{10, 10},
	}})
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 2, HitBox: pixel.Rect{
		Min: pixel.Vec{20, 20},
		Max: pixel.Vec{30, 30},
	}})
	entityCanvas.AddEventBox(EventBox{Event: "BOOM", Box: pixel.Rect{
		Min: pixel.Vec{25, 25},
		Max: pixel.Vec{26, 26},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_2_Entities_1_EventBox_2_overlaps(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{20, 20},
		Max: pixel.Vec{30, 30},
	}})
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 2, HitBox: pixel.Rect{
		Min: pixel.Vec{30, 20},
		Max: pixel.Vec{40, 30},
	}})
	entityCanvas.AddEventBox(EventBox{Event: "BOOM", Box: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{100, 100},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_1_Entity_2_EventBoxes_2_overlaps(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{20, 20},
		Max: pixel.Vec{30, 30},
	}})
	entityCanvas.AddEventBox(EventBox{Event: "BOOM1", Box: pixel.Rect{
		Min: pixel.Vec{21, 21},
		Max: pixel.Vec{22, 22},
	}})
	entityCanvas.AddEventBox(EventBox{Event: "BOOM2", Box: pixel.Rect{
		Min: pixel.Vec{29, 29},
		Max: pixel.Vec{30, 30},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}
