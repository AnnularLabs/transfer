package netease

type MusicList struct {
	Name     string   `json:"name"`
	Songs    []string `json:"songs"`
	SongsCnt int      `json:"songs_cnt"`
}

type PlaylistResponse struct {
	Code     int `json:"code"`
	Playlist struct {
		Id         int64    `json:"id"`
		Name       string   `json:"name"`
		Tracks     []*track `json:"tracks"`
		TrackCount int      `json:"trackCount"`
	} `json:"playlist"`
}

type track struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Ar   []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"ar"`
}
