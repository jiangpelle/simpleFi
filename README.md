# Web3 DeFi Platform

A decentralized finance (DeFi) platform built with Solidity smart contracts and Go backend.

## Features

- **Decentralized Exchange (DEX)**: Trade tokens directly on the blockchain
- **Lending Protocol**: Borrow and lend assets with smart contract automation
- **Yield Farming**: Stake tokens to earn rewards
- **Authentication**: Secure user authentication with JWT tokens

## Tech Stack

- **Smart Contracts**: Solidity
- **Backend**: Go
- **Frontend**: React (planned)
- **Blockchain**: Ethereum

## Project Structure

```
.
├── contracts/           # Solidity smart contracts
│   ├── Dex.sol         # Decentralized exchange
│   ├── Lending.sol     # Lending protocol
│   └── Farming.sol     # Yield farming
├── backend/            # Go backend
│   └── auth.go         # Authentication service
└── frontend/          # React frontend (planned)
```

## Setup

### Prerequisites

- Node.js (v14+)
- Go (v1.16+)
- Hardhat
- Ganache or local Ethereum node

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/web3.git
cd web3
```

2. Install dependencies:
```bash
npm install
```

3. Compile smart contracts:
```bash
npx hardhat compile
```

4. Run tests:
```bash
npx hardhat test
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 