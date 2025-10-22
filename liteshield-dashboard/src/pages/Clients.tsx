import { Table, Thead, Tbody, Tr, Th, Td, Spinner } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import axios from "axios"

export default function Clients() {
  const [clients, setClients] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    axios.get("http://localhost:8090/admin/clients", {
      headers: { "X-Admin-Key": "supersecret" }
    }).then(res => setClients(res.data)).finally(() => setLoading(false))
  }, [])

  if (loading) return <Spinner />

  return (
    <Table variant="simple" color="white">
      <Thead>
        <Tr><Th>Client</Th><Th>Requests</Th><Th>Cache Size</Th></Tr>
      </Thead>
      <Tbody>
        {clients.map(c => (
          <Tr key={c.id}>
            <Td>{c.name}</Td>
            <Td>{c.requests}</Td>
            <Td>{c.cache_size}</Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  )
}
