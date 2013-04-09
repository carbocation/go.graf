package forum

import (
	"time"
)

type Delta struct {
	Id         int64     "Unique identifier of this delta"
	PostId     int64     "ID of the modified post"
	TitleDelta string    "Changes made to the title, if any"
	BodyDelta  string    "Changes made to the body, if any"
	Modified   time.Time "Time at which the changes were made"
	ModifierId int64     "ID of the user who modified the post"
}
