import { SimpleGrid } from "@chakra-ui/react"
import StatCard from "../components/StatCard"
import { useEffect, useState } from "react"
import axios from "axios"

export default function Dashboard() {
  const [stats, setStats] = useState({ requests: 0, cacheHits: 0, uptime: "0s" })

  useEffect(() => {
    axios.get("http://localhost:8090/admin/metrics").then(res => {
      setStats({
        requests: res.data.requests || 0,
        cacheHits: res.data.cache_hits || 0,
        uptime: res.data.uptime || "0s"
      })
    }).catch(() => {
      console.warn("Proxy metrics unavailable")
    })
  }, [])

  return (
    <SimpleGrid columns={[1, 3]} spacing="6">
      <StatCard title="Total Requests" value={stats.requests} />
      <StatCard title="Cache Hits" value={stats.cacheHits} />
      <StatCard title="Uptime" value={stats.uptime} />
    </SimpleGrid>
  )
}
