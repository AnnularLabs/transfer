package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"transfer/internal/domain"

	"github.com/zmb3/spotify"
)

const (
	SpotifyBatchLimit = 100
)

var (
	searchLimit = 1
)

type SpotifyService interface {
	GetUserInfo(ctx context.Context, userID string) (string, error)
	GetPlaylistsForUser(ctx context.Context, userID string) ([]*PlaylistInfo, error)
	TransferTracksWithUserClient(ctx context.Context, client spotify.Client, playlistID string, tracks []domain.Track) (*domain.TransferResult, error)
}

type PlaylistInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type spotifyService struct {
	client spotify.Client
}

func NewSpotifyService(client spotify.Client) SpotifyService {
	return &spotifyService{
		client: client,
	}
}

func (s *spotifyService) GetUserInfo(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID cannot be empty")
	}

	resp, err := s.client.GetUsersPublicProfile(spotify.ID(userID))
	if err != nil {
		return "", fmt.Errorf("failed to get user profile: %w", err)
	}

	return resp.DisplayName, nil
}

func (s *spotifyService) GetPlaylistsForUser(ctx context.Context, userID string) ([]*PlaylistInfo, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	resp, err := s.client.GetPlaylistsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	result := make([]*PlaylistInfo, 0, len(resp.Playlists))
	for _, playlist := range resp.Playlists {
		result = append(result, &PlaylistInfo{
			Name: playlist.Name,
			ID:   playlist.ID.String(),
		})
	}

	return result, nil
}

func (s *spotifyService) TransferTracksWithUserClient(ctx context.Context, client spotify.Client, playlistID string, tracks []domain.Track) (*domain.TransferResult, error) {
	if playlistID == "" {
		return nil, errors.New("playlist ID cannot be empty")
	}

	result := &domain.TransferResult{
		TotalTracks:   len(tracks),
		SuccessCount:  0,
		FailedTracks:  make([]domain.FailedTrack, 0),
		SuccessTracks: make([]string, 0),
	}

	// 批量处理
	for i := 0; i < len(tracks); i += SpotifyBatchLimit {
		end := i + SpotifyBatchLimit
		if end > len(tracks) {
			end = len(tracks)
		}

		batchTracks := tracks[i:end]
		s.processBatch(ctx, client, playlistID, batchTracks, result)
	}

	return result, nil
}

func (s *spotifyService) processBatch(ctx context.Context, client spotify.Client, playlistID string, tracks []domain.Track, result *domain.TransferResult) {
	trackIDs := make([]spotify.ID, 0, len(tracks))

	for _, track := range tracks {
		spotifyID, err := s.searchTrack(ctx, track)
		if err != nil {
			result.FailedTracks = append(result.FailedTracks, domain.FailedTrack{
				Track: track,
				Error: err.Error(),
			})
			continue
		}

		trackIDs = append(trackIDs, spotifyID)
		result.SuccessTracks = append(result.SuccessTracks, string(spotifyID))
	}

	if len(trackIDs) == 0 {
		return
	}

	// 添加到歌单
	_, err := client.AddTracksToPlaylist(spotify.ID(playlistID), trackIDs...)
	if err != nil {
		// 如果批量添加失败，将所有歌曲标记为失败
		for i := range trackIDs {
			result.FailedTracks = append(result.FailedTracks, domain.FailedTrack{
				Track: tracks[i],
				Error: fmt.Sprintf("failed to add to playlist: %s", err.Error()),
			})
		}
		return
	}

	result.SuccessCount += len(trackIDs)
}

// searchTrack 搜索单首歌曲
func (s *spotifyService) searchTrack(ctx context.Context, track domain.Track) (spotify.ID, error) {
	// 构建搜索查询，优先使用 artist + title 的组合
	query := buildSearchQuery(track)

	resp, err := s.client.SearchOpt(query, spotify.SearchTypeTrack, &spotify.Options{
		Limit: &searchLimit,
	})
	if err != nil {
		return "", fmt.Errorf("search failed for track %s: %w", track.Title, err)
	}

	if resp.Tracks == nil || len(resp.Tracks.Tracks) == 0 {
		return "", fmt.Errorf("no results found for track: %s by %s", track.Title, track.Artist)
	}

	return resp.Tracks.Tracks[0].ID, nil
}

// buildSearchQuery 构建搜索查询字符串
func buildSearchQuery(track domain.Track) string {
	var parts []string

	if track.Title != "" {
		parts = append(parts, fmt.Sprintf("track:\"%s\"", track.Title))
	}

	if track.Artist != "" {
		parts = append(parts, fmt.Sprintf("artist:\"%s\"", track.Artist))
	}

	return strings.Join(parts, " ")
}
