import {
  Box,
  HStack,
  Text,
  Circle,
  Divider,
  VStack,
  Icon,
} from '@chakra-ui/react'
import { FaSpotify, FaUser, FaMusic, FaList, FaCheck, FaDownload } from 'react-icons/fa'

interface StepIndicatorProps {
  currentStep: number
  steps: string[]
}

function StepIndicator({ currentStep, steps }: StepIndicatorProps) {
  const stepIcons = [
    FaSpotify,    // Spotify Login
    FaUser,       // User Confirm
    FaMusic,      // NetEase Input
    FaList,       // Playlist Display
    FaCheck,      // Track Selection
    FaDownload,   // Transfer Music
  ]

  return (
    <HStack spacing={0} justify="center" align="center" w="full" overflowX="auto" py={4}>
      {steps.map((step, index) => {
        const stepNumber = index + 1
        const isActive = stepNumber === currentStep
        const isCompleted = stepNumber < currentStep
        const IconComponent = stepIcons[index]
        
        return (
          <HStack key={stepNumber} spacing={0} flex="1" justify="center">
            <VStack spacing={2} align="center" minW="120px">
              <Circle
                size={["40px", "48px"]}
                bg={isActive ? 'spotify.500' : isCompleted ? 'spotify.400' : 'gray.200'}
                color={isActive || isCompleted ? 'white' : 'gray.500'}
                fontWeight="bold"
                transition="all 0.3s ease"
                boxShadow={isActive ? "0 4px 12px rgba(29, 185, 84, 0.3)" : "none"}
                border={isActive ? "3px solid" : "none"}
                borderColor={isActive ? "spotify.200" : "transparent"}
              >
                {isCompleted ? (
                  <Icon as={FaCheck} boxSize={5} />
                ) : (
                  <Icon as={IconComponent} boxSize={4} />
                )}
              </Circle>
              <Text
                fontSize={["xs", "sm"]}
                color={isActive ? 'spotify.600' : isCompleted ? 'spotify.500' : 'gray.500'}
                fontWeight={isActive ? '600' : '500'}
                textAlign="center"
                transition="all 0.3s ease"
                noOfLines={2}
              >
                {step}
              </Text>
            </VStack>
            {index < steps.length - 1 && (
              <Box flex="1" px={2}>
                <Divider 
                  orientation="horizontal" 
                  borderColor={isCompleted ? 'spotify.300' : 'gray.300'}
                  borderWidth="2px"
                  transition="all 0.3s ease"
                />
              </Box>
            )}
          </HStack>
        )
      })}
    </HStack>
  )
}

export default StepIndicator
