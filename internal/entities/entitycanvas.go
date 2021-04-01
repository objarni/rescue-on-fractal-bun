package entities

import (
	"fmt"
	"github.com/faiface/pixel"
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

func (canvas *EntityCanvas) Consequences(handler func(eb EventBox, ehb EntityHitBox)) {
	for _, eb := range canvas.EventBoxes {
		for _, ehb := range canvas.EntityHitBoxes {
			if eb.Box.Intersect(ehb.HitBox).Area() > 0 {
				handler(eb, ehb)
			}
		}
	}
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
		eventBoxText += fmt.Sprintf("%v happened at %v\n", eb.Event, eb.Box)
	}
	consequences := ""
	canvas.Consequences(func(eb EventBox, ehb EntityHitBox) {
		consequences += fmt.Sprintf("- Entity %d needs to handle %v\n", ehb.Entity, eb.Event)
	})
	consequenceTitle := "Consequences:"
	if consequences == "" {
		consequenceTitle = "No hits detected"
	}
	return fmt.Sprintf("%v%v%v\n%v", entityText, eventBoxText, consequenceTitle, consequences)
}

type EventBox struct {
	Event string
	Box   pixel.Rect
}

func MakeEntityCanvas() EntityCanvas {
	return EntityCanvas{}
}
