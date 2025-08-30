import { useState } from 'react'
import {
  Box,
  Flex,
  Text,
  VStack,
  Container,
  Button,
  HStack,
  useToast,
  Heading,
} from '@chakra-ui/react'
import { FaSpotify, FaMusic, FaArrowRight, FaCheck } from 'react-icons/fa'
import NetEaseInput from './components/NetEaseInput'
import PlaylistDisplay from './components/PlaylistDisplay'
import SpotifyAuth from './components/SpotifyAuth'
import SpotifyUserConfirm from './components/SpotifyUserConfirm'
import TrackSelection from './components/TrackSelection'
import TransferProcess from './components/TransferProcess'
import StepIndicator from './components/StepIndicator'
import { api } from './services/api'

export interface Track {
  title: string
  artist: string
  album?: string
  match_key: string
}

export interface Playlist {
  name: string
  id: string
  tracks: Track[]
}

export interface SpotifyPlaylist {
  name: string
  id: string
}

function App() {
  // State Management - Spotify Auth comes first
  const [currentStep, setCurrentStep] = useState(1)
  const [isSpotifyAuthed, setIsSpotifyAuthed] = useState(false)
  const [spotifyUserInfo, setSpotifyUserInfo] = useState<any>(null)
  const [spotifyPlaylists, setSpotifyPlaylists] = useState<SpotifyPlaylist[]>([])
  const [selectedSpotifyPlaylist, setSelectedSpotifyPlaylist] = useState('')
  
  // NetEase related states
  const [neteasePlaylistId, setNeteasePlaylistId] = useState('')
  const [playlist, setPlaylist] = useState<Playlist | null>(null)
  const [selectedTracks, setSelectedTracks] = useState<Track[]>([])
  const [transferProgress, setTransferProgress] = useState(0)

  const toast = useToast()

  const handleLogout = async () => {
    try {
      await api.logout()
      // Reset all states
      setCurrentStep(1)
      setIsSpotifyAuthed(false)
      setSpotifyUserInfo(null)
      setSpotifyPlaylists([])
      setSelectedSpotifyPlaylist('')
      setNeteasePlaylistId('')
      setPlaylist(null)
      setSelectedTracks([])
      setTransferProgress(0)
      
      toast({
        title: 'Logged Out',
        description: 'Successfully logged out from Spotify',
        status: 'info',
        duration: 3000,
        isClosable: true,
      })
    } catch (error) {
      toast({
        title: 'Logout Failed',
        description: 'Error occurred during logout',
        status: 'error',
        duration: 3000,
        isClosable: true,
      })
    }
  }

  // Updated step names in English, with Spotify auth first
  const steps = [
    'Spotify Login',
    'User Confirm', 
    'NetEase Input',
    'Playlist Display',
    'Track Selection',
    'Transfer Music'
  ]

  return (
    <Container maxW="7xl" py={6}>
      <VStack spacing={8} align="stretch">
        {/* Modern Header */}
        <Flex justify="space-between" align="center" mb={4}>
          <VStack align="start" spacing={1}>
            <Heading 
              size="2xl" 
              bgGradient="linear(to-r, spotify.500, spotify.700)"
              bgClip="text"
              fontWeight="800"
            >
              Transfer
            </Heading>
            <Text color="gray.600" fontSize="lg">
              Transfer your playlists from NetEase to Spotify
            </Text>
          </VStack>
          
          {isSpotifyAuthed && spotifyUserInfo && (
            <HStack spacing={4}>
              <VStack spacing={0} align="end">
                <Text fontSize="sm" fontWeight="600" color="gray.700">
                  Welcome, {spotifyUserInfo.display_name || spotifyUserInfo.id}
                </Text>
                <Text fontSize="xs" color="gray.500">
                  Spotify Connected
                </Text>
              </VStack>
              <Button 
                size="sm" 
                variant="outline" 
                colorScheme="red"
                borderRadius="12px"
                onClick={handleLogout}
              >
                Logout
              </Button>
            </HStack>
          )}
        </Flex>

        {/* Step Indicator */}
        <StepIndicator currentStep={currentStep} steps={steps} />

        {/* Main Content Card */}
        <Flex justify="center">
          <Box 
            width="100%" 
            maxW="4xl"
            minH="70vh"
            bg="white"
            borderRadius="24px" 
            boxShadow="0 8px 32px rgba(0, 0, 0, 0.1)"
            border="1px solid"
            borderColor="gray.200"
            overflow="hidden"
          >
            <Box p={8}>
              {/* Step 1: Spotify Authentication */}
              {currentStep === 1 && (
                <SpotifyAuth
                  isAuthed={isSpotifyAuthed}
                  onAuthSuccess={(userInfo: any) => {
                    setIsSpotifyAuthed(true)
                    setSpotifyUserInfo(userInfo)
                    setCurrentStep(2)
                  }}
                  onBack={() => {}} // No back button on first step
                />
              )}

              {/* Step 2: Spotify User Confirmation */}
              {currentStep === 2 && (
                <SpotifyUserConfirm
                  userInfo={spotifyUserInfo}
                  playlists={spotifyPlaylists}
                  setPlaylists={setSpotifyPlaylists}
                  selectedPlaylist={selectedSpotifyPlaylist}
                  setSelectedPlaylist={setSelectedSpotifyPlaylist}
                  onNext={() => setCurrentStep(3)}
                  onBack={() => setCurrentStep(1)}
                />
              )}

              {/* Step 3: NetEase Input */}
              {currentStep === 3 && (
                <NetEaseInput
                  playlistId={neteasePlaylistId}
                  setPlaylistId={setNeteasePlaylistId}
                  onNext={(playlist: any) => {
                    setPlaylist(playlist)
                    setCurrentStep(4)
                  }}
                />
              )}

              {/* Step 4: Playlist Display */}
              {currentStep === 4 && playlist && (
                <PlaylistDisplay
                  playlist={playlist}
                  onNext={() => setCurrentStep(5)}
                  onBack={() => setCurrentStep(3)}
                />
              )}

              {/* Step 5: Track Selection */}
              {currentStep === 5 && playlist && (
                <TrackSelection
                  playlist={playlist}
                  selectedTracks={selectedTracks}
                  setSelectedTracks={setSelectedTracks}
                  onNext={() => setCurrentStep(6)}
                  onBack={() => setCurrentStep(4)}
                />
              )}

              {/* Step 6: Transfer Process */}
              {currentStep === 6 && (
                <TransferProcess
                  selectedTracks={selectedTracks}
                  spotifyPlaylistId={selectedSpotifyPlaylist}
                  progress={transferProgress}
                  setProgress={setTransferProgress}
                  onComplete={() => {
                    toast({
                      title: 'Transfer Complete!',
                      description: 'Your music has been successfully transferred to Spotify',
                      status: 'success',
                      duration: 5000,
                      isClosable: true,
                    })
                  }}
                  onBack={() => setCurrentStep(5)}
                />
              )}
            </Box>
          </Box>
        </Flex>
      </VStack>
    </Container>
  )
}

export default App
