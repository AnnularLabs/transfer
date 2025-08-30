import { useState, useEffect } from 'react'
import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  Alert,
  AlertIcon,
  useToast,
  Spinner,
  Select,
  FormLabel,
  FormControl,
} from '@chakra-ui/react'
import { api } from '../services/api'
import { SpotifyPlaylist } from '../App'

interface SpotifyUserConfirmProps {
  userInfo: any // 从授权步骤传递过来的用户信息
  playlists: SpotifyPlaylist[]
  setPlaylists: (playlists: SpotifyPlaylist[]) => void
  selectedPlaylist: string
  setSelectedPlaylist: (id: string) => void
  onNext: () => void
  onBack: () => void
}

function SpotifyUserConfirm({ 
  userInfo,
  playlists,
  setPlaylists,
  selectedPlaylist,
  setSelectedPlaylist,
  onNext, 
  onBack 
}: SpotifyUserConfirmProps) {
  const [isLoadingPlaylists, setIsLoadingPlaylists] = useState(false)
  const toast = useToast()

  useEffect(() => {
    // 组件加载时自动获取歌单
    if (userInfo) {
      loadPlaylists()
    }
  }, [userInfo])

  const loadPlaylists = async () => {
    setIsLoadingPlaylists(true)
    try {
      const userPlaylists = await api.fetchSpotifyPlaylists()
      setPlaylists(userPlaylists)
    } catch (error) {
      toast({
        title: '获取歌单失败',
        description: '无法获取您的 Spotify 歌单',
        status: 'error',
        duration: 3000,
        isClosable: true,
      })
    } finally {
      setIsLoadingPlaylists(false)
    }
  }

  const canProceed = selectedPlaylist

  return (
    <VStack spacing={6} align="stretch">
      <Text fontSize="2xl" fontWeight="bold" textAlign="center" color="#319795">
        确认 Spotify 用户
      </Text>

      <Box p={4} borderWidth={1} borderRadius="lg" borderColor="#319795" bg="gray.50">
        <VStack spacing={4}>
          <Alert status="success">
            <AlertIcon />
            <Box>
              <Text fontWeight="bold">授权成功！</Text>
              <Text>用户名: {userInfo?.display_name || userInfo?.id}</Text>
              <Text fontSize="sm" color="gray.600">
                ID: {userInfo?.id}
              </Text>
            </Box>
          </Alert>
        </VStack>
      </Box>

      <Box p={4} borderWidth={1} borderRadius="lg" borderColor="#319795" bg="gray.50">
        <FormControl>
          <FormLabel>选择目标歌单</FormLabel>
          {isLoadingPlaylists ? (
            <HStack justify="center" p={4}>
              <Spinner size="sm" />
              <Text>加载歌单中...</Text>
            </HStack>
          ) : (
            <Select
              placeholder="请选择要添加歌曲的歌单"
              value={selectedPlaylist}
              onChange={(e) => setSelectedPlaylist(e.target.value)}
              borderColor="#319795"
              _focus={{ borderColor: '#319795', boxShadow: '0 0 0 1px #319795' }}
            >
              {playlists.map((playlist) => (
                <option key={playlist.id} value={playlist.id}>
                  {playlist.name}
                </option>
              ))}
            </Select>
          )}
        </FormControl>
      </Box>

      <Text fontSize="sm" color="gray.500" textAlign="center">
        您已成功授权 Spotify 账户，请选择要添加歌曲的目标歌单
      </Text>

      <HStack justify="space-between" pt={4}>
        <Button
          variant="outline"
          colorScheme="teal"
          onClick={onBack}
          size="lg"
        >
          返回
        </Button>
        <Button
          colorScheme="teal"
          onClick={onNext}
          size="lg"
          isDisabled={!canProceed}
        >
          下一步：选择歌曲
        </Button>
      </HStack>
    </VStack>
  )
}

export default SpotifyUserConfirm
