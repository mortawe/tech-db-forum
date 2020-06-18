package models

type Params struct {
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
	Sort  string `json:"sort"`
}

