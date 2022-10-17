package gameModel

import (
	"time"

	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
)

type Game struct {
	ID           *uuid.UUID         `json:"id"`
	GameCategory *uuid.UUID         `json:"game_category"`
	Name         *map[string]string `json:"name"`
	Alias        *map[string]string `json:"alias"`
	PlayerCount  *int32             `json:"player_count"`
	Img          *string            `json:"image"`
	Description  *map[string]string `json:"description"`
	Highlight    *bool              `json:"highlight"`
	Order        *int32             `json:"order"`
	CreatedAt    *time.Time         `json:"created_at"`
}
