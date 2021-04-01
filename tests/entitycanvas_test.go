package tests

import (
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal/entities"
	"testing"
)

func Test_oneEntityOverlappingOneEventBox(t *testing.T) {
	entityCanvas := entities.MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(entities.EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{10, 10},
	}})
	entityCanvas.AddEventBox(entities.EventBox{Event: "ACTION", Box: pixel.Rect{
		Min: pixel.Vec{5, 5},
		Max: pixel.Vec{6, 6},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_oneEntityNoOverlapping(t *testing.T) {
	entityCanvas := entities.MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(entities.EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{10, 10},
	}})
	entityCanvas.AddEventBox(entities.EventBox{Event: "ACTION", Box: pixel.Rect{
		Min: pixel.Vec{55, 55},
		Max: pixel.Vec{76, 76},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_twoEntityOneOverlap(t *testing.T) {
	entityCanvas := entities.MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(entities.EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{10, 10},
	}})
	entityCanvas.AddEntityHitBox(entities.EntityHitBox{Entity: 2, HitBox: pixel.Rect{
		Min: pixel.Vec{20, 20},
		Max: pixel.Vec{30, 30},
	}})
	entityCanvas.AddEventBox(entities.EventBox{Event: "BOOM", Box: pixel.Rect{
		Min: pixel.Vec{25, 25},
		Max: pixel.Vec{26, 26},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_twoEntitiesOneEventBox(t *testing.T) {
	entityCanvas := entities.MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(entities.EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{20, 20},
		Max: pixel.Vec{30, 30},
	}})
	entityCanvas.AddEntityHitBox(entities.EntityHitBox{Entity: 2, HitBox: pixel.Rect{
		Min: pixel.Vec{30, 20},
		Max: pixel.Vec{40, 30},
	}})
	entityCanvas.AddEventBox(entities.EventBox{Event: "BOOM", Box: pixel.Rect{
		Min: pixel.Vec{0, 0},
		Max: pixel.Vec{100, 100},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

func Test_oneEntityTwoEventBoxes(t *testing.T) {
	entityCanvas := entities.MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(entities.EntityHitBox{Entity: 1, HitBox: pixel.Rect{
		Min: pixel.Vec{20, 20},
		Max: pixel.Vec{30, 30},
	}})
	entityCanvas.AddEventBox(entities.EventBox{Event: "BOOM1", Box: pixel.Rect{
		Min: pixel.Vec{21, 21},
		Max: pixel.Vec{22, 22},
	}})
	entityCanvas.AddEventBox(entities.EventBox{Event: "BOOM2", Box: pixel.Rect{
		Min: pixel.Vec{29, 29},
		Max: pixel.Vec{30, 30},
	}})
	approvals.VerifyString(
		t,
		entityCanvas.String(),
	)
}

/* notes / test list - general behaviour
#   adding an entity with a hitbox HB, and an event box that overlaps
#   adding an entity, and a event box that does not overlap the entity
#   adding two entities, and a event box, that overlaps only one entity
#   adding two entities, and a event box that overlaps both entities
   adding an entity, and two hitboxes HB1, HB2, both overlapping the entity
*/
