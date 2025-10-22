import { Box, Flex } from "@chakra-ui/react"
import Sidebar from "./components/Sidebar"
import Header from "./components/Header"
import Dashboard from "./pages/Dashboard"
import Clients from "./pages/Clients"
import Logs from "./pages/Logs"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"

export default function App() {
  return (
    <Router>
      <Flex>
        <Sidebar />
        <Box flex="1" ml="240px" bg="gray.900" minH="100vh" p="6">
          <Header />
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/clients" element={<Clients />} />
            <Route path="/logs" element={<Logs />} />
          </Routes>
        </Box>
      </Flex>
    </Router>
  )
}
