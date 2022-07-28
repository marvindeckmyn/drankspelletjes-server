package accountModel

import "github.com/marvindeckmyn/drankspelletjes-server/uuid"

type Account struct {
	ID       *uuid.UUID `json:"id"`
	Name     *string    `json:"name"`
	Email    *string    `json:"email"`
	Password *string    `json:"password"`
}
