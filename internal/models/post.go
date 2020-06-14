package models

import "time"

type Post struct {
	Author   string `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string `json:"forum"`
	ID       int    `json:"id"`
	IsEdited bool   `json:"isEdited"`
	Message  string `json:"message"`
	Parent   int    `json:"parent"`
	Thread   int    `json:"thread"`
}

type PostDetails struct {
	*User `json:"author"`
	*Post	`json:"post"`
	*Thread	`json:"thread"`
	*Forum `json:"forum"`
}