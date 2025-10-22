import { Box, Text } from "@chakra-ui/react"

export default function StatCard({ title, value }: { title: string; value: string | number }) {
  return (
    <Box bg="gray.700" color="white" p="5" rounded="lg" shadow="md" textAlign="center" flex="1">
      <Text fontSize="lg" fontWeight="semibold">{title}</Text>
      <Text fontSize="3xl" fontWeight="bold" mt="2">{value}</Text>
    </Box>
  )
}
