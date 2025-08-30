import { useState, useEffect } from 'react'
import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  Progress,
  Alert,
  AlertIcon,
  AlertTitle,
  AlertDescription,
  List,
  ListItem,
  Badge,
  useToast,
} from '@chakra-ui/react'
import { api } from '../services/api'
import { Track } from '../App'

interface TransferProcessProps {
  selectedTracks: Track[]
  spotifyPlaylistId: string
  progress: number
  setProgress: (progress: number) => void
  onComplete: () => void
  onBack: () => void
}

function TransferProcess({ 
  selectedTracks, 
  spotifyPlaylistId, 
  progress, 
  setProgress, 
  onComplete, 
  onBack 
}: TransferProcessProps) {
  const [isTransferring, setIsTransferring] = useState(false)
  const [transferComplete, setTransferComplete] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const toast = useToast()

  const startTransfer = async () => {
    setIsTransferring(true)
    setError(null)
    setProgress(0)

    try {
      // 模拟进度更新
      const interval = setInterval(() => {
        setProgress(prev => {
          if (prev >= 90) {
            clearInterval(interval)
            return 90
          }
          return prev + 10
        })
      }, 500)

      // 执行实际的迁移
      await api.transferTracks(spotifyPlaylistId, selectedTracks)
      
      clearInterval(interval)
      setProgress(100)
      setTransferComplete(true)
      
      toast({
        title: '迁移完成！',
        description: `成功迁移 ${selectedTracks.length} 首歌曲`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      })
      
      setTimeout(() => {
        onComplete()
      }, 2000)
      
    } catch (error) {
      setError('迁移过程中发生错误，请稍后重试')
      toast({
        title: '迁移失败',
        description: '请检查网络连接或稍后重试',
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
    } finally {
      setIsTransferring(false)
    }
  }

  useEffect(() => {
    // 组件加载时自动开始迁移
    startTransfer()
  }, [])

  return (
    <VStack spacing={6} align="stretch" h="full">
      <Text fontSize="2xl" fontWeight="bold" textAlign="center" color="#319795">
        {transferComplete ? '迁移完成' : '正在迁移歌曲'}
      </Text>

      <Box p={6} borderWidth={1} borderRadius="lg" borderColor="#319795" bg="gray.50">
        <VStack spacing={4}>
          <HStack justify="space-between" w="full">
            <Text fontSize="lg" fontWeight="semibold">
              迁移进度
            </Text>
            <Badge 
              colorScheme={transferComplete ? 'green' : 'blue'} 
              variant="subtle"
              px={3}
              py={1}
            >
              {selectedTracks.length} 首歌曲
            </Badge>
          </HStack>
          
          <Box w="full">
            <Progress 
              value={progress} 
              size="lg" 
              colorScheme={transferComplete ? 'green' : 'teal'}
              hasStripe={!transferComplete}
              isAnimated={!transferComplete}
              borderRadius="md"
            />
            <Text textAlign="center" mt={2} fontSize="sm" color="gray.600">
              {Math.round(progress)}% 完成
            </Text>
          </Box>

          {isTransferring && (
            <Text fontSize="sm" color="gray.600" textAlign="center">
              正在将歌曲添加到您的 Spotify 歌单...
            </Text>
          )}

          {transferComplete && (
            <Alert status="success">
              <AlertIcon />
              <Box>
                <AlertTitle>迁移成功！</AlertTitle>
                <AlertDescription>
                  所有选中的歌曲已成功添加到您的 Spotify 歌单
                </AlertDescription>
              </Box>
            </Alert>
          )}

          {error && (
            <Alert status="error">
              <AlertIcon />
              <Box>
                <AlertTitle>迁移失败</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Box>
            </Alert>
          )}
        </VStack>
      </Box>

      <Box 
        flex={1}
        overflowY="auto"
        borderWidth={1} 
        borderRadius="lg" 
        borderColor="#319795"
        p={4}
        maxH="300px"
      >
        <Text fontSize="lg" fontWeight="semibold" mb={3}>
          迁移歌曲列表
        </Text>
        <List spacing={2}>
          {selectedTracks.map((track, index) => (
            <ListItem key={track.match_key || index}>
              <HStack justify="space-between">
                <VStack align="start" spacing={1} flex={1}>
                  <Text fontWeight="medium" fontSize="sm">{track.title}</Text>
                  <Text fontSize="xs" color="gray.600">
                    {track.artist} {track.album && `• ${track.album}`}
                  </Text>
                </VStack>
                <Badge 
                  colorScheme={transferComplete ? 'green' : 'gray'} 
                  variant="subtle"
                  size="sm"
                >
                  {transferComplete ? '完成' : '等待中'}
                </Badge>
              </HStack>
            </ListItem>
          ))}
        </List>
      </Box>

      <HStack justify="space-between" pt={4}>
        <Button
          variant="outline"
          colorScheme="teal"
          onClick={onBack}
          size="lg"
          isDisabled={isTransferring}
        >
          返回
        </Button>
        
        {error && (
          <Button
            colorScheme="teal"
            onClick={startTransfer}
            size="lg"
            isLoading={isTransferring}
            loadingText="重试中"
          >
            重试
          </Button>
        )}
        
        {transferComplete && (
          <Button
            colorScheme="green"
            onClick={onComplete}
            size="lg"
          >
            完成
          </Button>
        )}
      </HStack>
    </VStack>
  )
}

export default TransferProcess
