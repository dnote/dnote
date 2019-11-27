package handlers

import (
	"fmt"
)

func (a *App) getNoteURL(uuid string) string {
	return fmt.Sprintf("%s/notes/%s", a.WebURL, uuid)
}
