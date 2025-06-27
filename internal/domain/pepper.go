package domain

type Pepper struct {
	Name        string `json:"name"`
	ShuLow      int    `json:"shu_low"`
	ShuMid      int    `json:"shu_mid"`
	ShuUp       int    `json:"shu_up"`
	Description string `json:"description"`
}
