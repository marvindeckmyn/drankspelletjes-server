package gameModel

import (
	"time"

	"github.com/marvindeckmyn/drankspelletjes-server/uuid"
)

type Game struct {
	ID            *uuid.UUID         `json:"id"`
	GameCategory  *uuid.UUID         `json:"game_category"`
	GameNecessity *uuid.UUID         `json:"game_necessity"`
	Name          *map[string]string `json:"name"`
	Alias         *map[string]string `json:"alias"`
	PlayerCount   *int32             `json:"player_count"`
	Image         *string            `json:"image"`
	Credits       *string            `json:"credits"`
	Description   *map[string]string `json:"description"`
	Highlight     *bool              `json:"highlight"`
	Views         *int32             `json:"views"`
	Order         *int32             `json:"order"`
	CreatedAt     *time.Time         `json:"created_at"`
}
