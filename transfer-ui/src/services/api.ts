import { Playlist, SpotifyPlaylist, Track } from '../App'

const API_BASE = '/api'

export const api = {
  // 获取网易云歌单
  fetchNeteasePlaylist: async (playlistId: string): Promise<Playlist> => {
    const response = await fetch(`${API_BASE}/netease/playlist?id=${playlistId}`)
    if (!response.ok) {
      throw new Error('Failed to fetch NetEase playlist')
    }
    return response.json()
  },

  // Spotify OAuth 相关
  getSpotifyAuthUrl: (): string => {
    return `${API_BASE}/user/auth/spotify/login`
  },

  checkSpotifyAuthStatus: async (): Promise<{authenticated: boolean, user_id?: string, message?: string}> => {
    try {
      const response = await fetch(`${API_BASE}/user/auth/spotify/status`, {
        method: 'POST',
        credentials: 'include' // 包含 cookies
      })
      if (!response.ok) {
        return { authenticated: false }
      }
      return await response.json()
    } catch {
      return { authenticated: false }
    }
  },

  // 获取当前用户信息 (不再需要 userId 参数)
  getCurrentSpotifyUser: async (): Promise<any> => {
    const response = await fetch(`${API_BASE}/spotify/me`, {
      credentials: 'include' // 包含 cookies
    })
    if (!response.ok) {
      throw new Error('Failed to get current user')
    }
    return response.json()
  },

  // 获取 Spotify 用户歌单 (不再需要手动传递认证信息)
  fetchSpotifyPlaylists: async (): Promise<SpotifyPlaylist[]> => {
    const response = await fetch(`${API_BASE}/spotify/playlists`, {
      credentials: 'include' // 包含 cookies
    })
    if (!response.ok) {
      throw new Error('Failed to fetch Spotify playlists')
    }
    return response.json()
  },

  // 迁移歌曲到 Spotify 歌单 (不再需要手动传递认证信息)
  transferTracks: async (playlistId: string, tracks: Track[]): Promise<void> => {
    const trackNames = tracks.map(track => track.title)
    
    const response = await fetch(`${API_BASE}/spotify/playlists/${playlistId}/tracks`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include', // 包含 cookies
      body: JSON.stringify({
        track_names: trackNames
      })
    })
    
    if (!response.ok) {
      throw new Error('Failed to transfer tracks')
    }
  },

  // 登出
  logout: async (): Promise<void> => {
    const response = await fetch(`${API_BASE}/user/auth/spotify/logout`, {
      method: 'POST',
      credentials: 'include'
    })
    if (!response.ok) {
      throw new Error('Failed to logout')
    }
  }
}
