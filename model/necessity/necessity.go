package necessityModel

import "github.com/marvindeckmyn/drankspelletjes-server/uuid"

type Necessity struct {
	ID   *uuid.UUID         `json:"id"`
	Name *map[string]string `json:"name"`
}
