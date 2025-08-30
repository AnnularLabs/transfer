package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"transfer/internal/domain"
)

type NeteaseService interface {
	GetPlaylist(ctx context.Context, id int64) (*domain.MusicList, error)
}

const (
	TargetPattern = "https://music.163.com/api/v6/playlist/detail?id=%d"
)

type neteaseService struct {
	client *http.Client
}

func NewNeteaseService() NeteaseService {
	return &neteaseService{
		client: http.DefaultClient,
	}
}

func (n *neteaseService) GetPlaylist(ctx context.Context, nid int64) (*domain.MusicList, error) {
	if nid <= 0 {
		return nil, fmt.Errorf("invalid playlist ID: %d", nid)
	}

	target := fmt.Sprintf(TargetPattern, nid)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlist: %w", err)
	}
	defer resp.Body.Close()

	var apiResp PlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Code != 200 {
		return nil, fmt.Errorf("API returned error code: %d", apiResp.Code)
	}

	return n.convertToMusicList(&apiResp), nil
}

// convertToMusicList 将 API 响应转换为领域对象
// 好品味：数据转换逻辑独立，可测试
func (n *neteaseService) convertToMusicList(resp *PlaylistResponse) *domain.MusicList {
	tracks := make([]domain.Track, 0, len(resp.Playlist.Tracks))

	for _, track := range resp.Playlist.Tracks {
		// 构建艺术家名称
		artists := make([]string, 0, len(track.Ar))
		for _, artist := range track.Ar {
			if artist.Name != "" {
				artists = append(artists, artist.Name)
			}
		}

		artistName := strings.Join(artists, ", ")

		domainTrack := domain.Track{
			Title:    track.Name,
			Artist:   artistName,
			Album:    track.Al.Name,
			MatchKey: buildMatchKey(track.Name, artistName),
		}

		tracks = append(tracks, domainTrack)
	}

	return &domain.MusicList{
		Name:   resp.Playlist.Name,
		ID:     fmt.Sprintf("%d", resp.Playlist.Id),
		Tracks: tracks,
	}
}

// buildMatchKey 构建用于匹配的键
func buildMatchKey(title, artist string) string {
	// 简单的标准化：去除空格，转小写
	key := strings.ToLower(strings.TrimSpace(title))
	if artist != "" {
		key += "|" + strings.ToLower(strings.TrimSpace(artist))
	}
	return key
}

type PlaylistResponse struct {
	Code int `json:"code"`

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
	Al struct {
		Name string `json:"name"`
	} `json:"al"`
}
