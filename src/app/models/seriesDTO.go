package models

type SeriesModel struct {
    ID int `json:"id"`
    Title string `json:"title"`
    Status string `json:"status"`
	Episodes int `json:"episodes"`
    LastEpisode int `json:"last_episode"`
    Ranking int `json:"ranking"`
}
