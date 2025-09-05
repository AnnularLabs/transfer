package domain

// Track 代表一首歌曲的完整信息
type Track struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album,omitempty"`
	// 用于匹配的唯一标识，组合 title + artist
	MatchKey string `json:"match_key"`
}

// MusicList 歌单的完整表示
type MusicList struct {
	Name   string  `json:"name"`
	ID     string  `json:"id"`
	Tracks []Track `json:"tracks"`
}

// TransferResult 传输结果，明确记录成功/失败
type TransferResult struct {
	TotalTracks   int           `json:"total_tracks"`
	SuccessCount  int           `json:"success_count"`
	FailedTracks  []FailedTrack `json:"failed_tracks"`
	SuccessTracks []string      `json:"success_tracks"` // Spotify track IDs
}

// FailedTrack 失败的歌曲，不静默忽略
type FailedTrack struct {
	Track Track  `json:"track"`
	Error string `json:"error"`
}
