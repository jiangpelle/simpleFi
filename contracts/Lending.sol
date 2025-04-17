// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Lending is ReentrancyGuard, Ownable {
    struct Market {
        uint256 totalBorrows;
        uint256 totalSupply;
        uint256 borrowRate;
        uint256 supplyRate;
        uint256 exchangeRate;
        uint256 collateralFactor;
        bool isListed;
    }

    mapping(address => Market) public markets;
    mapping(address => mapping(address => uint256)) public borrows;
    mapping(address => mapping(address => uint256)) public supplies;
    mapping(address => uint256) public prices;

    uint256 public constant BASE = 1e18;
    uint256 public constant COLLATERAL_FACTOR = 0.75 * BASE;

    event MarketListed(address indexed token);
    event Supply(address indexed user, address indexed token, uint256 amount);
    event Withdraw(address indexed user, address indexed token, uint256 amount);
    event Borrow(address indexed user, address indexed token, uint256 amount);
    event Repay(address indexed user, address indexed token, uint256 amount);

    function listMarket(address token, uint256 price) external onlyOwner {
        require(!markets[token].isListed, "Market already listed");
        markets[token] = Market({
            totalBorrows: 0,
            totalSupply: 0,
            borrowRate: 0,
            supplyRate: 0,
            exchangeRate: BASE,
            collateralFactor: COLLATERAL_FACTOR,
            isListed: true
        });
        prices[token] = price;
        emit MarketListed(token);
    }

    function supply(address token, uint256 amount) external nonReentrant {
        require(markets[token].isListed, "Market not listed");
        require(amount > 0, "Amount must be greater than 0");

        IERC20(token).transferFrom(msg.sender, address(this), amount);
        
        Market storage market = markets[token];
        uint256 supplyTokens = (amount * BASE) / market.exchangeRate;
        
        market.totalSupply += amount;
        supplies[msg.sender][token] += supplyTokens;
        
        emit Supply(msg.sender, token, amount);
    }

    function withdraw(address token, uint256 amount) external nonReentrant {
        require(markets[token].isListed, "Market not listed");
        require(amount > 0, "Amount must be greater than 0");

        Market storage market = markets[token];
        uint256 supplyTokens = supplies[msg.sender][token];
        require(supplyTokens >= amount, "Insufficient supply");

        uint256 withdrawAmount = (amount * market.exchangeRate) / BASE;
        require(market.totalSupply >= withdrawAmount, "Insufficient liquidity");

        market.totalSupply -= withdrawAmount;
        supplies[msg.sender][token] -= amount;
        
        IERC20(token).transfer(msg.sender, withdrawAmount);
        
        emit Withdraw(msg.sender, token, withdrawAmount);
    }

    function borrow(address token, uint256 amount) external nonReentrant {
        require(markets[token].isListed, "Market not listed");
        require(amount > 0, "Amount must be greater than 0");

        Market storage market = markets[token];
        require(market.totalSupply >= amount, "Insufficient liquidity");

        uint256 borrowValue = amount * prices[token];
        uint256 collateralValue = calculateCollateralValue(msg.sender);
        require(borrowValue <= collateralValue * market.collateralFactor / BASE, "Insufficient collateral");

        market.totalBorrows += amount;
        borrows[msg.sender][token] += amount;
        
        IERC20(token).transfer(msg.sender, amount);
        
        emit Borrow(msg.sender, token, amount);
    }

    function repay(address token, uint256 amount) external nonReentrant {
        require(markets[token].isListed, "Market not listed");
        require(amount > 0, "Amount must be greater than 0");

        Market storage market = markets[token];
        require(borrows[msg.sender][token] >= amount, "Insufficient borrow");

        IERC20(token).transferFrom(msg.sender, address(this), amount);
        
        market.totalBorrows -= amount;
        borrows[msg.sender][token] -= amount;
        
        emit Repay(msg.sender, token, amount);
    }

    function calculateCollateralValue(address user) public view returns (uint256) {
        uint256 totalValue;
        for (uint256 i = 0; i < getListedMarketsCount(); i++) {
            address token = getListedMarket(i);
            uint256 supply = supplies[user][token];
            if (supply > 0) {
                totalValue += (supply * prices[token]) / BASE;
            }
        }
        return totalValue;
    }

    function getListedMarketsCount() public view returns (uint256) {
        uint256 count;
        for (uint256 i = 0; i < 100; i++) {
            if (markets[address(uint160(i))].isListed) {
                count++;
            }
        }
        return count;
    }

    function getListedMarket(uint256 index) public view returns (address) {
        uint256 count;
        for (uint256 i = 0; i < 100; i++) {
            address token = address(uint160(i));
            if (markets[token].isListed) {
                if (count == index) {
                    return token;
                }
                count++;
            }
        }
        revert("Market not found");
    }
} 