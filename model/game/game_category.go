package gameModel

import "github.com/marvindeckmyn/drankspelletjes-server/uuid"

type GameCategory struct {
	ID          *uuid.UUID         `json:"id"`
	Name        *map[string]string `json:"name"`
	Description *map[string]string `json:"description"`
	Image       *string            `json:"image"`
	Order       *int32             `json:"order"`
}
