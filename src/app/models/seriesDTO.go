package models

type SeriesModel struct {
    ID int `json:"id"`
    Title string `json:"title"`
    Status string `json:"status"`
	Episodes int `json:"totalEpisodes"`
    LastEpisode int `json:"lastEpisodeWatched"`
    Ranking int `json:"ranking"`
}
