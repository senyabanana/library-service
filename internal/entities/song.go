package entities

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
