package netease

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"transfer/utils"
)

type neteaseService struct{}

var NeteaseService = neteaseService{}

// TODO env控制
const (
	NETEASY_URL_V6 = "https://music.163.com/api/v6/playlist/detail" // 获取歌单详细数据
	// chunkSize    = 500
)

func (n neteaseService) GetPlayListMusic(ctx context.Context, id string) (*MusicList, error) {

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	resp, err := utils.Post(NETEASY_URL_V6, strings.NewReader("id="+id))
	if err != nil {
		utils.Errorf("fail result: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	playlistResponse := &PlaylistResponse{}
	err = json.Unmarshal(body, playlistResponse)

	if err != nil {
		utils.Errorf("fail unmarshal: %v", err)
		return nil, err
	}

	if playlistResponse.Code == 401 {
		utils.Errorf("fail, 无权访问: %v", err)
		return nil, err
	}

	playListName := playlistResponse.Playlist.Name
	tracks := playlistResponse.Playlist.Tracks
	tracksCount := playlistResponse.Playlist.TrackCount
	return &MusicList{
		Name:     playListName,
		Songs:    n.getTracks(tracks),
		SongsCnt: tracksCount,
	}, nil

}

func (n neteaseService) getTracks(tracks []*track) []string {
	strings := make([]string, 0, len(tracks))
	for _, v := range tracks {
		var res = ""
		if len(v.Ar) == 1 {
			res = fmt.Sprintf("%s - %s", v.Name, v.Ar[0].Name)
		} else {
			res = fmt.Sprintf("%s - %s", v.Name, "未知")
		}
		strings = append(strings, res)
	}
	return strings
}
