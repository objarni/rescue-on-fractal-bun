package entities

import (
	"fmt"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal/events"
	internal "objarni/rescue-on-fractal-bun/internal/printers"
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
	for _, eventBox := range canvas.EventBoxes {
		for _, entityHitBox := range canvas.EntityHitBoxes {
			if eventBox.Box.Intersect(entityHitBox.HitBox).Area() > 0 {
				handler(eventBox, entityHitBox)
			}
		}
	}
}

func (canvas *EntityCanvas) String() string {
	entityText := ""
	for ix := range canvas.EntityHitBoxes {
		ehb := canvas.EntityHitBoxes[ix]
		entityText += fmt.Sprintf(
			"Entity %d is at %v\n",
			ehb.Entity,
			internal.PrintRect(ehb.HitBox),
		)
	}
	eventBoxText := ""
	for ix := range canvas.EventBoxes {
		eb := canvas.EventBoxes[ix]
		eventBoxText += fmt.Sprintf("Event: %v\n", eb)
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
	Event events.Event
	Box   pixel.Rect
}

func (eb EventBox) String() string {
	return fmt.Sprintf("%v %v", eb.Event, internal.PrintRect(eb.Box))
}

func MakeEntityCanvas() EntityCanvas {
	return EntityCanvas{}
}
