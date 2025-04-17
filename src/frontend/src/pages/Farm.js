import React, { useState, useEffect } from 'react';
import {
  Box,
  Grid,
  Paper,
  Typography,
  Button,
  TextField,
  CircularProgress,
  Card,
  CardContent,
  CardActions
} from '@mui/material';
import { useWeb3 } from '../contexts/Web3Context';

const Farm = () => {
  const { account, provider } = useWeb3();
  const [pools, setPools] = useState([]);
  const [loading, setLoading] = useState(true);
  const [depositAmount, setDepositAmount] = useState('');
  const [selectedPool, setSelectedPool] = useState(null);

  useEffect(() => {
    const fetchPools = async () => {
      try {
        const response = await fetch('/api/farm/pools');
        const data = await response.json();
        setPools(data.pools);
      } catch (error) {
        console.error('Error fetching pools:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchPools();
  }, []);

  const handleDeposit = async (poolId) => {
    if (!account || !provider) {
      alert('Please connect your wallet');
      return;
    }

    try {
      // Implement deposit logic using smart contract
      const amount = ethers.utils.parseEther(depositAmount);
      // Call smart contract deposit function
      // Update UI after successful deposit
    } catch (error) {
      console.error('Error depositing:', error);
      alert('Failed to deposit');
    }
  };

  const handleWithdraw = async (poolId) => {
    if (!account || !provider) {
      alert('Please connect your wallet');
      return;
    }

    try {
      // Implement withdraw logic using smart contract
      // Call smart contract withdraw function
      // Update UI after successful withdraw
    } catch (error) {
      console.error('Error withdrawing:', error);
      alert('Failed to withdraw');
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Liquidity Mining
      </Typography>

      <Grid container spacing={3}>
        {pools.map((pool) => (
          <Grid item xs={12} md={6} key={pool.id}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  {pool.name}
                </Typography>
                <Typography color="textSecondary">
                  APR: {pool.apr}%
                </Typography>
                <Typography color="textSecondary">
                  Total Staked: {pool.totalStaked}
                </Typography>
                <Typography color="textSecondary">
                  Your Stake: {pool.userStake}
                </Typography>
                <Typography color="textSecondary">
                  Pending Rewards: {pool.pendingRewards}
                </Typography>
              </CardContent>
              <CardActions>
                <TextField
                  size="small"
                  label="Amount"
                  value={depositAmount}
                  onChange={(e) => setDepositAmount(e.target.value)}
                  sx={{ mr: 1 }}
                />
                <Button
                  variant="contained"
                  color="primary"
                  onClick={() => handleDeposit(pool.id)}
                >
                  Deposit
                </Button>
                <Button
                  variant="outlined"
                  color="secondary"
                  onClick={() => handleWithdraw(pool.id)}
                >
                  Withdraw
                </Button>
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Farm; 