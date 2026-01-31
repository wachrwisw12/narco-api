package models

type Menu struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
	Icon  string `json:"icon"`
	Path  string `json:"path"`
	Sort  int    `json:"sort"`
}
