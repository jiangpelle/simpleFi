import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Web3Provider } from './contexts/Web3Context';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Swap from './pages/Swap';
import Pool from './pages/Pool';
import Stake from './pages/Stake';

const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90caf9',
    },
    secondary: {
      main: '#f48fb1',
    },
  },
});

function App() {
  return (
    <Web3Provider>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Router>
          <Layout>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/swap" element={<Swap />} />
              <Route path="/pool" element={<Pool />} />
              <Route path="/stake" element={<Stake />} />
            </Routes>
          </Layout>
        </Router>
      </ThemeProvider>
    </Web3Provider>
  );
}

export default App; 