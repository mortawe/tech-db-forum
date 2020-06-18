package models

type Status struct {
	Forums  int `json:"forum"`
	Threads int `json:"thread"`
	Posts   int `json:"post"`
	Users   int `json:"user"`
}

