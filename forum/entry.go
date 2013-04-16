package forum

import (
	"time"
)

// Put ModifiedBy, ModifiedAuthor in a separate table. A post can only be 
// created once but modified an infinite number of times.
type Entry struct {
	Id       int64     "The ID of the post"
	Title    string    "Title of the post. Will be empty for entries that are really intended to be comments."
	Body     string    "Contents of the post. Will be empty for entries that are intended to be links."
	Created  time.Time "Time at which the post was created."
	AuthorId int64     "ID of the author of the post"
}

//Note: why not just manage what is real and what is not through methods?
//The 'get' methods would populate the view-type methods automatically
//or leave them blank if irrelevant. The 'set' methods would store the 
//essential / non-derived fields and ignore the others. That way you don't
//have to juggle two view types
type EntryView struct {
	Entry
	Points   int64
	HasVoted bool
}
