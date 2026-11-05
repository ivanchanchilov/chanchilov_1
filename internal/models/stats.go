package models

type CategoryStat struct {
	Name  string  `json:"name"`
	Total float64 `json:"total"`
}

type MonthStat struct {
	Month int     `json:"month"`
	Total float64 `json:"total"`
}

type StatsPageData struct {
	From string
	To   string
	Year int
}
