import { Flex, Text, Spacer, Button } from "@chakra-ui/react"

export default function Header() {
  return (
    <Flex bg="gray.800" color="white" p="4" align="center" position="sticky" top="0" zIndex="10">
      <Text fontWeight="bold">LiteShield Admin</Text>
      <Spacer />
      <Button colorScheme="red" size="sm">Logout</Button>
    </Flex>
  )
}
