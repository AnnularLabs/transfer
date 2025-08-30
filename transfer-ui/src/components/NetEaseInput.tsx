import { useState } from 'react'
import {
  VStack,
  HStack,
  Input,
  Button,
  Text,
  useToast,
  Spinner,
  Box,
  Heading,
  Icon,
  Center,
  InputGroup,
  InputRightElement,
  FormControl,
  FormLabel,
  FormHelperText,
} from '@chakra-ui/react'
import { FaMusic, FaSearch, FaExternalLinkAlt } from 'react-icons/fa'
import { api } from '../services/api'
import { Playlist } from '../App'

interface NetEaseInputProps {
  playlistId: string
  setPlaylistId: (id: string) => void
  onNext: (playlist: Playlist) => void
}

function NetEaseInput({ playlistId, setPlaylistId, onNext }: NetEaseInputProps) {
  const [isLoading, setIsLoading] = useState(false)
  const toast = useToast()

  const handleSearch = async () => {
    if (!playlistId.trim()) {
      toast({
        title: 'Playlist ID Required',
        description: 'Please enter a NetEase playlist ID',
        status: 'warning',
        duration: 3000,
        isClosable: true,
      })
      return
    }

    setIsLoading(true)
    try {
      const playlist = await api.fetchNeteasePlaylist(playlistId.trim())
      onNext(playlist)
    } catch (error) {
      toast({
        title: 'Failed to Fetch Playlist',
        description: 'Please check if the playlist ID is correct and try again',
        status: 'error',
        duration: 3000,
        isClosable: true,
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch()
    }
  }

  return (
    <Center minH="500px">
      <VStack spacing={8} maxW="lg" w="full">
        {/* Header Section */}
        <VStack spacing={4} textAlign="center">
          <Icon as={FaMusic} boxSize={16} color="netease.500" />
          <Heading size="xl" color="gray.800">
            NetEase Music Playlist
          </Heading>
          <Text fontSize="lg" color="gray.600" maxW="md">
            Enter your NetEase Music playlist ID to start the transfer process
          </Text>
        </VStack>

        {/* Input Form */}
        <Box 
          w="full" 
          bg="white" 
          p={8} 
          borderRadius="20px" 
          border="1px solid"
          borderColor="gray.200"
          boxShadow="0 4px 20px rgba(0, 0, 0, 0.08)"
        >
          <VStack spacing={6}>
            <FormControl>
              <FormLabel fontWeight="600" color="gray.700">
                Playlist ID
              </FormLabel>
              <InputGroup size="lg">
                <Input
                  placeholder="Enter NetEase playlist ID (e.g., 123456789)"
                  value={playlistId}
                  onChange={(e) => setPlaylistId(e.target.value)}
                  onKeyPress={handleKeyPress}
                  isDisabled={isLoading}
                  bg="gray.50"
                  border="2px solid"
                  borderColor="gray.200"
                  _focus={{ 
                    borderColor: 'netease.500', 
                    boxShadow: '0 0 0 1px var(--chakra-colors-netease-500)',
                    bg: 'white'
                  }}
                  _hover={{
                    borderColor: 'gray.300'
                  }}
                />
                <InputRightElement>
                  <Icon as={FaMusic} color="gray.400" />
                </InputRightElement>
              </InputGroup>
              <FormHelperText color="gray.500" fontSize="sm">
                <HStack spacing={1} align="center">
                  <Icon as={FaExternalLinkAlt} boxSize={3} />
                  <Text>
                    Find the ID in the URL: music.163.com/playlist?id=<strong>123456789</strong>
                  </Text>
                </HStack>
              </FormHelperText>
            </FormControl>

            <Button
              colorScheme="netease"
              size="lg"
              w="full"
              h="14"
              onClick={handleSearch}
              isLoading={isLoading}
              loadingText="Fetching Playlist..."
              isDisabled={!playlistId.trim()}
              leftIcon={<FaSearch />}
              fontSize="lg"
              fontWeight="600"
              _hover={{
                transform: 'translateY(-2px)',
                boxShadow: '0 8px 25px rgba(214, 48, 49, 0.25)',
              }}
            >
              {isLoading ? 'Fetching Playlist...' : 'Get Playlist'}
            </Button>
          </VStack>
        </Box>

        {/* Example Section */}
        <Box 
          w="full" 
          p={6} 
          bg="gray.50" 
          borderRadius="16px"
          border="1px solid"
          borderColor="gray.200"
        >
          <VStack spacing={3} textAlign="center">
            <Text fontSize="sm" fontWeight="600" color="gray.700">
              How to find your playlist ID?
            </Text>
            <Text fontSize="sm" color="gray.600">
              1. Go to NetEase Music and open your playlist
            </Text>
            <Text fontSize="sm" color="gray.600">
              2. Copy the numbers after "id=" in the URL
            </Text>
            <Text fontSize="sm" color="gray.600">
              3. Paste the ID in the input field above
            </Text>
          </VStack>
        </Box>
      </VStack>
    </Center>
  )
}

export default NetEaseInput
