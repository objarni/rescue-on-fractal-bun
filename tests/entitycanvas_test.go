package tests

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"testing"
)

type EntityHitBox struct {
	Entity int
	HitBox pixel.Rect
}

type EntityCanvas struct {
	EntityHitBoxes []EntityHitBox
	EventBoxes     []EventBox
}

func (canvas *EntityCanvas) AddEntityHitBox(entityHitBox EntityHitBox) {
	canvas.EntityHitBoxes = append(canvas.EntityHitBoxes, entityHitBox)
}

func (canvas *EntityCanvas) AddEventBox(eventBox EventBox) {
	canvas.EventBoxes = append(canvas.EventBoxes, eventBox)
}

func (canvas *EntityCanvas) String() string {
	entityText := ""
	for ix := range canvas.EntityHitBoxes {
		ehb := canvas.EntityHitBoxes[ix]
		entityText += fmt.Sprintf("Entity %d is at %v\n", ehb.Entity, ehb.HitBox)
	}
	eventBoxText := ""
	for ix := range canvas.EventBoxes {
		eb := canvas.EventBoxes[ix]
		eventBoxText += fmt.Sprintf("%v happened at %v", eb.Event, eb.Box)
	}
	consequences := ""
	for _, eb := range canvas.EventBoxes {
		for _, ehb := range canvas.EntityHitBoxes {
			if eb.Box.Intersect(ehb.HitBox).Area() > 0 {
				consequences = fmt.Sprintf("- Entity %d needs to handle %v\n", ehb.Entity, eb.Event)
			}
		}
	}
	consequenceTitle := "Consequences:"
	if consequences == "" {
		consequenceTitle = "No hits detected"
	}
	return fmt.Sprintf("%v%v\n%v\n%v", entityText, eventBoxText, consequenceTitle, consequences)
}

type EventBox struct {
	Event string
	Box   pixel.Rect
}

func Test_oneEntityOverlappingOneEventBox(t *testing.T) {
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

func Test_oneEntityNoOverlapping(t *testing.T) {
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

func Test_twoEntityOneOverlap(t *testing.T) {
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

func MakeEntityCanvas() EntityCanvas {
	return EntityCanvas{}
}

/* notes / test list - general behaviour
#   adding an entity with a hitbox HB, and an event box that overlaps
#   adding an entity, and a hitbox that does not overlap the entity
   adding two entities, and a hitbox HB1, that overlaps only one entity
   adding two entities, and a hitbox HB1 that overlaps both entities
   adding an entity, and two hitboxes HB1, HB2, both overlapping the entity
*/
