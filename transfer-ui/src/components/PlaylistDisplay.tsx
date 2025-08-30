import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  List,
  ListItem,
  Divider,
  Badge,
} from '@chakra-ui/react'
import { Playlist } from '../App'

interface PlaylistDisplayProps {
  playlist: Playlist
  onNext: () => void
  onBack: () => void
}

function PlaylistDisplay({ playlist, onNext, onBack }: PlaylistDisplayProps) {
  return (
    <VStack spacing={6} align="stretch" h="full">
      <Text fontSize="2xl" fontWeight="bold" textAlign="center" color="#319795">
        歌单信息
      </Text>

      <Box 
        p={4} 
        borderWidth={1} 
        borderRadius="lg" 
        borderColor="#319795"
        bg="gray.50"
      >
        <VStack align="start" spacing={2}>
          <Text fontSize="xl" fontWeight="bold">
            {playlist.name}
          </Text>
          <HStack>
            <Badge colorScheme="teal" variant="subtle">
              歌单ID: {playlist.id}
            </Badge>
            <Badge colorScheme="green" variant="subtle">
              {playlist.tracks.length} 首歌曲
            </Badge>
          </HStack>
        </VStack>
      </Box>

      <Box 
        flex={1}
        overflowY="auto"
        borderWidth={1} 
        borderRadius="lg" 
        borderColor="#319795"
        p={4}
        maxH="400px"
      >
        <Text fontSize="lg" fontWeight="semibold" mb={3}>
          歌曲列表
        </Text>
        <List spacing={2}>
          {playlist.tracks.map((track, index) => (
            <ListItem key={index}>
              <VStack align="start" spacing={1}>
                <Text fontWeight="medium">{track.title}</Text>
                <Text fontSize="sm" color="gray.600">
                  {track.artist} {track.album && `• ${track.album}`}
                </Text>
              </VStack>
              {index < playlist.tracks.length - 1 && <Divider my={2} />}
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
        >
          返回
        </Button>
        <Button
          colorScheme="teal"
          onClick={onNext}
          size="lg"
        >
          下一步：Spotify 授权
        </Button>
      </HStack>
    </VStack>
  )
}

export default PlaylistDisplay
