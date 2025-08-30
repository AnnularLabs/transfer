import { useState } from 'react'
import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  Checkbox,
  List,
  ListItem,
  Divider,
  Badge,
} from '@chakra-ui/react'
import { Playlist, Track } from '../App'

interface TrackSelectionProps {
  playlist: Playlist
  selectedTracks: Track[]
  setSelectedTracks: (tracks: Track[]) => void
  onNext: () => void
  onBack: () => void
}

function TrackSelection({ 
  playlist, 
  selectedTracks, 
  setSelectedTracks, 
  onNext, 
  onBack 
}: TrackSelectionProps) {
  const [selectAll, setSelectAll] = useState(false)

  const handleSelectAll = () => {
    if (selectAll) {
      setSelectedTracks([])
    } else {
      setSelectedTracks([...playlist.tracks])
    }
    setSelectAll(!selectAll)
  }

  const handleTrackToggle = (track: Track) => {
    const isSelected = selectedTracks.some(t => t.match_key === track.match_key)
    
    if (isSelected) {
      setSelectedTracks(selectedTracks.filter(t => t.match_key !== track.match_key))
    } else {
      setSelectedTracks([...selectedTracks, track])
    }
  }

  const isTrackSelected = (track: Track) => {
    return selectedTracks.some(t => t.match_key === track.match_key)
  }

  const canProceed = selectedTracks.length > 0

  return (
    <VStack spacing={6} align="stretch" h="full">
      <Text fontSize="2xl" fontWeight="bold" textAlign="center" color="#319795">
        选择要迁移的歌曲
      </Text>

      <Box p={4} borderWidth={1} borderRadius="lg" borderColor="#319795" bg="gray.50">
        <HStack justify="space-between" align="center">
          <HStack>
            <Text fontSize="lg" fontWeight="semibold">
              {playlist.name}
            </Text>
            <Badge colorScheme="teal" variant="subtle">
              {playlist.tracks.length} 首歌曲
            </Badge>
          </HStack>
          <HStack>
            <Badge colorScheme="green" variant="subtle">
              已选择 {selectedTracks.length} 首
            </Badge>
            <Checkbox
              isChecked={selectAll}
              onChange={handleSelectAll}
              colorScheme="teal"
            >
              全选
            </Checkbox>
          </HStack>
        </HStack>
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
        <List spacing={2}>
          {playlist.tracks.map((track, index) => (
            <ListItem key={track.match_key || index}>
              <HStack spacing={3} align="start">
                <Checkbox
                  isChecked={isTrackSelected(track)}
                  onChange={() => handleTrackToggle(track)}
                  colorScheme="teal"
                  mt={1}
                />
                <VStack align="start" spacing={1} flex={1}>
                  <Text fontWeight="medium">{track.title}</Text>
                  <Text fontSize="sm" color="gray.600">
                    {track.artist} {track.album && `• ${track.album}`}
                  </Text>
                </VStack>
              </HStack>
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
          isDisabled={!canProceed}
        >
          开始迁移 ({selectedTracks.length} 首歌曲)
        </Button>
      </HStack>
    </VStack>
  )
}

export default TrackSelection
