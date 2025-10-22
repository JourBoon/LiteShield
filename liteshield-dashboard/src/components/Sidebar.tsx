import { Box, VStack, Text, Button } from "@chakra-ui/react"
import { useNavigate } from "react-router-dom"

export default function Sidebar() {
  const navigate = useNavigate()
  return (
    <Box bg="gray.900" color="white" w="240px" h="100vh" p="6" position="fixed">
      <Text fontSize="2xl" fontWeight="bold" mb="10">
        LiteShield ⚡
      </Text>
      <VStack align="stretch" spacing="4">
        <Button variant="ghost" colorScheme="whiteAlpha" onClick={() => navigate("/")}>Dashboard</Button>
        <Button variant="ghost" colorScheme="whiteAlpha" onClick={() => navigate("/clients")}>Clients</Button>
        <Button variant="ghost" colorScheme="whiteAlpha" onClick={() => navigate("/logs")}>Logs</Button>
      </VStack>
    </Box>
  )
}
