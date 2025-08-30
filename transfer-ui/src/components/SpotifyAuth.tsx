import { useState, useEffect } from 'react'
import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  Spinner,
  Alert,
  AlertIcon,
  AlertTitle,
  AlertDescription,
  Heading,
  Icon,
  Center,
  Card,
  CardBody,
  useColorModeValue,
} from '@chakra-ui/react'
import { FaSpotify, FaMusic, FaArrowRight } from 'react-icons/fa'
import { api } from '../services/api'

interface SpotifyAuthProps {
  isAuthed: boolean
  onAuthSuccess: (userInfo: any) => void
  onBack: () => void
}

function SpotifyAuth({ isAuthed, onAuthSuccess, onBack }: SpotifyAuthProps) {
  const [isChecking, setIsChecking] = useState(false)
  const [authUrl, setAuthUrl] = useState('')

  useEffect(() => {
    // 获取授权URL
    const url = api.getSpotifyAuthUrl()
    setAuthUrl(url)

    // 检查是否已经授权
    checkAuthStatus()

    // 检查URL参数，看是否从授权回调返回
    const urlParams = new URLSearchParams(window.location.search)
    const authResult = urlParams.get('auth')
    const userParam = urlParams.get('user')
    
    if (authResult === 'success' && userParam) {
      // 清除URL参数
      window.history.replaceState({}, document.title, window.location.pathname)
      // 获取用户信息并触发成功回调
      fetchUserInfoAndProceed()
    } else if (authResult === 'error') {
      const error = urlParams.get('error')
      console.error('授权失败:', error)
      // 清除URL参数
      window.history.replaceState({}, document.title, window.location.pathname)
    }
  }, [])

  const checkAuthStatus = async () => {
    setIsChecking(true)
    try {
      const result = await api.checkSpotifyAuthStatus()
      if (result.authenticated) {
        const userInfo = await api.getCurrentSpotifyUser()
        onAuthSuccess(userInfo)
      }
    } catch (error) {
      console.error('Failed to check authorization status:', error)
    } finally {
      setIsChecking(false)
    }
  }

  const fetchUserInfoAndProceed = async () => {
    try {
      const userInfo = await api.getCurrentSpotifyUser()
      onAuthSuccess(userInfo)
    } catch (error) {
      console.error('Failed to get user info:', error)
    }
  }

  const handleAuth = () => {
    // Redirect to Spotify authorization page
    window.location.href = authUrl
  }

  if (isChecking) {
    return (
      <Center minH="500px">
        <VStack spacing={6}>
          <Spinner size="xl" color="spotify.500" thickness="4px" />
          <Text fontSize="lg" color="gray.600">
            Checking authorization status...
          </Text>
        </VStack>
      </Center>
    )
  }

  if (isAuthed) {
    return (
      <Center minH="500px">
        <VStack spacing={8} maxW="md" w="full">
          <Box textAlign="center">
            <Icon as={FaSpotify} boxSize={16} color="spotify.500" mb={4} />
            <Text fontSize="2xl" fontWeight="bold" color="gray.800" mb={2}>
              Already Connected!
            </Text>
            <Text color="gray.600">
              Your Spotify account is successfully connected
            </Text>
          </Box>
          
          <Button 
            colorScheme="spotify" 
            size="lg" 
            onClick={() => fetchUserInfoAndProceed()}
            rightIcon={<FaArrowRight />}
          >
            Continue
          </Button>
        </VStack>
      </Center>
    )
  }

  return (
    <Center minH="500px">
      <VStack spacing={8} maxW="lg" w="full">
        {/* Header Section */}
        <VStack spacing={4} textAlign="center">
          <Icon as={FaSpotify} boxSize={20} color="spotify.500" />
          <Heading size="xl" color="gray.800">
            Connect to Spotify
          </Heading>
          <Text fontSize="lg" color="gray.600" maxW="md">
            Connect your Spotify account to transfer your NetEase playlists seamlessly
          </Text>
        </VStack>

        {/* Permission Card */}
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
            <Text fontSize="md" fontWeight="600" color="gray.700">
              We'll need permission to:
            </Text>
            
            <VStack align="start" spacing={3} w="full">
              <HStack spacing={3}>
                <Box w={2} h={2} bg="spotify.500" borderRadius="full" />
                <Text fontSize="md" color="gray.600">
                  View your Spotify profile information
                </Text>
              </HStack>
              <HStack spacing={3}>
                <Box w={2} h={2} bg="spotify.500" borderRadius="full" />
                <Text fontSize="md" color="gray.600">
                  Access your playlists
                </Text>
              </HStack>
              <HStack spacing={3}>
                <Box w={2} h={2} bg="spotify.500" borderRadius="full" />
                <Text fontSize="md" color="gray.600">
                  Add songs to your playlists
                </Text>
              </HStack>
            </VStack>

            <Button
              colorScheme="spotify"
              size="lg"
              w="full"
              h="14"
              onClick={handleAuth}
              leftIcon={<FaSpotify />}
              fontSize="lg"
              fontWeight="600"
              _hover={{
                transform: 'translateY(-2px)',
                boxShadow: '0 8px 25px rgba(29, 185, 84, 0.25)',
              }}
            >
              Connect Spotify Account
            </Button>
          </VStack>
        </Box>

        {/* Footer Actions */}
        <HStack spacing={4} w="full" justify="center">
          <Button
            variant="ghost"
            colorScheme="gray"
            onClick={checkAuthStatus}
            isLoading={isChecking}
            leftIcon={<FaMusic />}
          >
            Check Connection Status
          </Button>
        </HStack>
      </VStack>
    </Center>
  )
}

export default SpotifyAuth
