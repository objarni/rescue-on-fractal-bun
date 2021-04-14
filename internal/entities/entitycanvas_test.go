package entities

import (
	approvals "github.com/approvals/go-approval-tests"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal/events"
	"objarni/rescue-on-fractal-bun/tests"
	"testing"
)

func Test_1_Entity_1_EventBox_1_overlap(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: px.Rect{
		Min: px.V(0, 0),
		Max: px.V(10, 10),
	}})
	entityCanvas.AddEventBox(EventBox{Event: events.Action, Box: px.Rect{
		Min: px.V(5, 5),
		Max: px.V(6, 6),
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
	entityCanvas.AddEntityHitBox(EntityHitBox{
		Entity: 1,
		HitBox: px.Rect{
			Min: px.ZV,
			Max: px.V(10, 10),
		},
	})
	entityCanvas.AddEventBox(EventBox{Event: events.ButtonPressed, Box: px.Rect{
		Min: px.V(55, 55),
		Max: px.V(76, 76),
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_2_Entities_1_EventBox_1_overlap(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: px.Rect{
		Min: px.ZV,
		Max: px.V(10, 10),
	}})
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 2, HitBox: px.Rect{
		Min: px.V(20, 20),
		Max: px.V(30, 30),
	}})
	entityCanvas.AddEventBox(EventBox{Event: events.Action, Box: px.Rect{
		Min: px.V(25, 25),
		Max: px.V(26, 26),
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_2_Entities_1_EventBox_2_overlaps(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: px.Rect{
		Min: px.V(20, 20),
		Max: px.V(30, 30),
	}})
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 2, HitBox: px.Rect{
		Min: px.V(30, 20),
		Max: px.V(40, 30),
	}})
	entityCanvas.AddEventBox(EventBox{Event: events.Action, Box: px.Rect{
		Min: px.V(0, 0),
		Max: px.V(100, 100),
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_1_Entity_2_EventBoxes_2_overlaps(t *testing.T) {
	entityCanvas := MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{Entity: 1, HitBox: px.Rect{
		Min: px.V(20, 20),
		Max: px.V(30, 30),
	}})
	entityCanvas.AddEventBox(EventBox{Event: events.Action, Box: px.Rect{
		Min: px.V(21, 21),
		Max: px.V(22, 22),
	}})
	entityCanvas.AddEventBox(EventBox{Event: events.ButtonPressed, Box: px.Rect{
		Min: px.V(29, 29),
		Max: px.V(30, 30),
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}
