package gameModel

import "github.com/marvindeckmyn/drankspelletjes-server/uuid"

type GameCategory struct {
	ID    *uuid.UUID         `json:"id"`
	Name  *map[string]string `json:"name"`
	Order *int32             `json:"order"`
}
