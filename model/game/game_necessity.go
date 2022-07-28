package gameModel

import "github.com/marvindeckmyn/drankspelletjes-server/uuid"

type GameNecessity struct {
	Game      *uuid.UUID `json:"game"`
	Necessity *uuid.UUID `json:"necessity"`
}
