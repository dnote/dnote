package infra

import (
	"time"
)

var (
	ActionAddNote    = "add_note"
	ActionRemoveNote = "remove_note"
	ActionEditNote   = "edit_note"
	ActionAddBook    = "add_book"
	ActionRemoveBook = "remove_book"
)

type Action struct {
	Type      string
	Data      map[string]interface{}
	Timestamp int64
}

func NewActionAddNote(uuid, content string) Action {
	return Action{
		Type: ActionAddNote,
		Data: map[string]interface{}{
			"UUID":    uuid,
			"Content": content,
		},
		Timestamp: time.Now().Unix(),
	}
}

func NewActionRemoveNote(uuid string) Action {
	return Action{
		Type: ActionRemoveNote,
		Data: map[string]interface{}{
			"UUID": uuid,
		},
		Timestamp: time.Now().Unix(),
	}
}

func NewActionEditNote(uuid, content string) Action {
	return Action{
		Type: ActionEditNote,
		Data: map[string]interface{}{
			"UUID":    uuid,
			"Content": content,
		},
		Timestamp: time.Now().Unix(),
	}
}

func NewActionAddBook(uuid, name string) Action {
	return Action{
		Type: ActionAddBook,
		Data: map[string]interface{}{
			"UUID": uuid,
			"Name": name,
		},
		Timestamp: time.Now().Unix(),
	}
}

func NewActionRemoveBook(uuid string) Action {
	return Action{
		Type: ActionRemoveBook,
		Data: map[string]interface{}{
			"UUID": uuid,
		},
		Timestamp: time.Now().Unix(),
	}
}
