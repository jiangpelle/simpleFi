// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Dex is ReentrancyGuard, Ownable {
    struct Pool {
        uint256 token0Reserve;
        uint256 token1Reserve;
        uint256 totalSupply;
    }

    mapping(address => mapping(address => Pool)) public pools;
    mapping(address => mapping(address => uint256)) public liquidity;
    mapping(address => uint256) public fees;

    event Swap(
        address indexed sender,
        address indexed tokenIn,
        address indexed tokenOut,
        uint256 amountIn,
        uint256 amountOut
    );

    event AddLiquidity(
        address indexed provider,
        address indexed token0,
        address indexed token1,
        uint256 amount0,
        uint256 amount1,
        uint256 liquidity
    );

    event RemoveLiquidity(
        address indexed provider,
        address indexed token0,
        address indexed token1,
        uint256 amount0,
        uint256 amount1,
        uint256 liquidity
    );

    function addLiquidity(
        address token0,
        address token1,
        uint256 amount0,
        uint256 amount1
    ) external nonReentrant {
        require(amount0 > 0 && amount1 > 0, "Amounts must be greater than 0");
        
        IERC20(token0).transferFrom(msg.sender, address(this), amount0);
        IERC20(token1).transferFrom(msg.sender, address(this), amount1);

        Pool storage pool = pools[token0][token1];
        uint256 liquidityMinted;
        
        if (pool.totalSupply == 0) {
            liquidityMinted = sqrt(amount0 * amount1);
        } else {
            liquidityMinted = min(
                (amount0 * pool.totalSupply) / pool.token0Reserve,
                (amount1 * pool.totalSupply) / pool.token1Reserve
            );
        }

        require(liquidityMinted > 0, "Insufficient liquidity minted");
        
        pool.token0Reserve += amount0;
        pool.token1Reserve += amount1;
        pool.totalSupply += liquidityMinted;
        liquidity[msg.sender][token0] += liquidityMinted;

        emit AddLiquidity(msg.sender, token0, token1, amount0, amount1, liquidityMinted);
    }

    function removeLiquidity(
        address token0,
        address token1,
        uint256 liquidityAmount
    ) external nonReentrant {
        require(liquidityAmount > 0, "Amount must be greater than 0");
        
        Pool storage pool = pools[token0][token1];
        require(liquidity[msg.sender][token0] >= liquidityAmount, "Insufficient liquidity");

        uint256 amount0 = (liquidityAmount * pool.token0Reserve) / pool.totalSupply;
        uint256 amount1 = (liquidityAmount * pool.token1Reserve) / pool.totalSupply;

        pool.token0Reserve -= amount0;
        pool.token1Reserve -= amount1;
        pool.totalSupply -= liquidityAmount;
        liquidity[msg.sender][token0] -= liquidityAmount;

        IERC20(token0).transfer(msg.sender, amount0);
        IERC20(token1).transfer(msg.sender, amount1);

        emit RemoveLiquidity(msg.sender, token0, token1, amount0, amount1, liquidityAmount);
    }

    function swap(
        address tokenIn,
        address tokenOut,
        uint256 amountIn
    ) external nonReentrant {
        require(amountIn > 0, "Amount must be greater than 0");
        
        Pool storage pool = pools[tokenIn][tokenOut];
        require(pool.totalSupply > 0, "Pool does not exist");

        uint256 amountOut = getAmountOut(amountIn, pool.token0Reserve, pool.token1Reserve);
        require(amountOut > 0, "Insufficient output amount");

        IERC20(tokenIn).transferFrom(msg.sender, address(this), amountIn);
        IERC20(tokenOut).transfer(msg.sender, amountOut);

        pool.token0Reserve += amountIn;
        pool.token1Reserve -= amountOut;

        emit Swap(msg.sender, tokenIn, tokenOut, amountIn, amountOut);
    }

    function getAmountOut(
        uint256 amountIn,
        uint256 reserveIn,
        uint256 reserveOut
    ) public pure returns (uint256) {
        require(amountIn > 0, "INSUFFICIENT_INPUT_AMOUNT");
        require(reserveIn > 0 && reserveOut > 0, "INSUFFICIENT_LIQUIDITY");
        uint256 amountInWithFee = amountIn * 997;
        uint256 numerator = amountInWithFee * reserveOut;
        uint256 denominator = reserveIn * 1000 + amountInWithFee;
        return numerator / denominator;
    }

    function sqrt(uint256 y) internal pure returns (uint256 z) {
        if (y > 3) {
            z = y;
            uint256 x = y / 2 + 1;
            while (x < z) {
                z = x;
                x = (y / x + x) / 2;
            }
        } else if (y != 0) {
            z = 1;
        }
    }

    function min(uint256 a, uint256 b) internal pure returns (uint256) {
        return a < b ? a : b;
    }
} 