package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"transfer/utils"
)

const (
	netEasyUrlV6 = "https://music.163.com/api/v6/playlist/detail" // 获取歌单详细数据
	chunkSize    = 500
)

type musicList struct {
	Name     string   `json:"name"`
	Songs    []string `json:"songs"`
	SongsCnt int      `json:"songs_cnt"`
}

type playlistResponse struct {
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

type Songs struct {
	Songs []struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
		Ar   []struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"ar"`
	} `json:"songs"`
}

func GetPlayListMusic(id string) (*musicList, error) {
	resp, err := utils.Post(netEasyUrlV6, strings.NewReader("id="+id))
	if err != nil {
		utils.Errorf("fail result: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	playlistResponse := &playlistResponse{}
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
	return &musicList{
		Name:     playListName,
		Songs:    getTracks(tracks),
		SongsCnt: tracksCount,
	}, nil
}

func getTracks(tracks []*track) []string {
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
