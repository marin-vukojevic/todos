package todos

import (
	"github.com/google/uuid"
)

type Todo struct {
	Uuid        uuid.UUID
	Description string
	Done        bool
}
