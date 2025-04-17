import React, { useState, useEffect } from 'react';
import { Box, Typography, Paper } from '@mui/material';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

const PriceDisplay = ({ token, interval = '1h' }) => {
  const [price, setPrice] = useState(null);
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPrice = async () => {
      try {
        const response = await fetch(`/api/prices/${token}`);
        const data = await response.json();
        setPrice(data.price);
      } catch (error) {
        console.error('Error fetching price:', error);
      }
    };

    const fetchHistory = async () => {
      try {
        const response = await fetch(`/api/prices/${token}/history?interval=${interval}`);
        const data = await response.json();
        setHistory(data);
      } catch (error) {
        console.error('Error fetching price history:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchPrice();
    fetchHistory();

    const priceInterval = setInterval(fetchPrice, 30000); // Update every 30 seconds
    const historyInterval = setInterval(fetchHistory, 300000); // Update every 5 minutes

    return () => {
      clearInterval(priceInterval);
      clearInterval(historyInterval);
    };
  }, [token, interval]);

  const chartData = {
    labels: history.map(item => new Date(item.timestamp).toLocaleTimeString()),
    datasets: [
      {
        label: `${token} Price`,
        data: history.map(item => item.price),
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: `${token} Price History`
      }
    },
    scales: {
      y: {
        beginAtZero: false
      }
    }
  };

  if (loading) {
    return (
      <Paper sx={{ p: 2 }}>
        <Typography>Loading price data...</Typography>
      </Paper>
    );
  }

  return (
    <Box>
      <Paper sx={{ p: 2, mb: 2 }}>
        <Typography variant="h6" gutterBottom>
          Current {token} Price
        </Typography>
        <Typography variant="h4">
          ${price?.toFixed(2) || 'N/A'}
        </Typography>
      </Paper>

      <Paper sx={{ p: 2 }}>
        <Line data={chartData} options={chartOptions} />
      </Paper>
    </Box>
  );
};

export default PriceDisplay; 