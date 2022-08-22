package gameModel

import "github.com/marvindeckmyn/drankspelletjes-server/uuid"

type GameNecessity struct {
	ID   *uuid.UUID         `json:"id"`
	Name *map[string]string `json:"name"`
}
