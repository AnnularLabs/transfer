import { extendTheme } from '@chakra-ui/react'

const theme = extendTheme({
  config: {
    initialColorMode: 'light',
    useSystemColorMode: false,
  },
  styles: {
    global: {
      body: {
        bg: 'gray.50',
        fontFamily: "'Inter', 'Segoe UI', 'Roboto', sans-serif",
      },
    },
  },
  fonts: {
    heading: "'Inter', 'Segoe UI', 'Roboto', sans-serif",
    body: "'Inter', 'Segoe UI', 'Roboto', sans-serif",
  },
  components: {
    Button: {
      baseStyle: {
        borderRadius: '16px',
        fontWeight: '600',
        transition: 'all 0.2s ease',
      },
      variants: {
        solid: {
          borderRadius: '16px',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
          _hover: {
            transform: 'translateY(-2px)',
            boxShadow: '0 8px 20px rgba(0, 0, 0, 0.15)',
          },
          _active: {
            transform: 'translateY(0px)',
          },
        },
        outline: {
          borderRadius: '16px',
          borderWidth: '2px',
          _hover: {
            transform: 'translateY(-1px)',
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
          },
        },
        ghost: {
          borderRadius: '16px',
          _hover: {
            transform: 'translateY(-1px)',
          },
        },
      },
    },
    Card: {
      baseStyle: {
        container: {
          borderRadius: '20px',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.08)',
          border: 'none',
          bg: 'white',
        },
      },
    },
    Input: {
      variants: {
        outline: {
          field: {
            borderRadius: '12px',
            borderWidth: '2px',
            _focus: {
              borderColor: 'spotify.500',
              boxShadow: '0 0 0 1px var(--chakra-colors-spotify-500)',
            },
          },
        },
      },
    },
    Select: {
      variants: {
        outline: {
          field: {
            borderRadius: '12px',
            borderWidth: '2px',
            _focus: {
              borderColor: 'spotify.500',
              boxShadow: '0 0 0 1px var(--chakra-colors-spotify-500)',
            },
          },
        },
      },
    },
  },
  colors: {
    spotify: {
      50: '#e8f5e8',
      100: '#c6e7c6',
      200: '#a3d9a3',
      300: '#7fcb7f',
      400: '#5bb65b',
      500: '#1db954',
      600: '#1aa34a',
      700: '#178d3f',
      800: '#147735',
      900: '#11612a',
    },
    netease: {
      50: '#ffe8e8',
      100: '#fcc6c6',
      200: '#f7a3a3',
      300: '#f27f7f',
      400: '#ed5b5b',
      500: '#d63031',
      600: '#c71c1c',
      700: '#b81717',
      800: '#a91313',
      900: '#9a0f0f',
    },
  },
})

export default theme
