import React from "react"
import ReactDOM from "react-dom/client"
import { ChakraProvider, Box, Heading, Button } from "@chakra-ui/react"

function App() {
  return (
    <ChakraProvider>
      <Box
        bgGradient="linear(to-br, blue.500, purple.600)"
        h="100vh"
        color="white"
        textAlign="center"
        pt="40"
      >
        <Heading mb={6}>LiteShield Dashboard ⚡</Heading>
        <Button colorScheme="whiteAlpha" variant="outline" size="lg">
          Tester
        </Button>
      </Box>
    </ChakraProvider>
  )
}

ReactDOM.createRoot(document.getElementById("root")!).render(<App />)
