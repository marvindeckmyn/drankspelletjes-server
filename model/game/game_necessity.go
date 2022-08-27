package gameModel

import "github.com/marvindeckmyn/drankspelletjes-server/uuid"

type GameNecessity struct {
	ID   *uuid.UUID         `json:"id"`
	Game *uuid.UUID         `json:"game"`
	Name *map[string]string `json:"name"`
}
