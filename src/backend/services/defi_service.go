package services

import (
	"errors"
	"time"

	"defi-backend/models"

	"gorm.io/gorm"
)

type DefiService struct {
	db *gorm.DB
}

func NewDefiService(db *gorm.DB) *DefiService {
	return &DefiService{db: db}
}

// DEX 相关服务
func (s *DefiService) CreateTrade(userID uint, pairID uint, tradeType string, amount, price float64) (*models.Trade, error) {
	trade := &models.Trade{
		UserID:     userID,
		PairID:     pairID,
		Type:       tradeType,
		Amount:     amount,
		Price:      price,
		TotalValue: amount * price,
		Status:     "pending",
	}

	if err := s.db.Create(trade).Error; err != nil {
		return nil, err
	}

	return trade, nil
}

func (s *DefiService) GetTradingPairs() ([]models.TradingPair, error) {
	var pairs []models.TradingPair
	if err := s.db.Find(&pairs).Error; err != nil {
		return nil, err
	}
	return pairs, nil
}

// 借贷相关服务
func (s *DefiService) CreateLendingPosition(userID uint, token string, amount float64, positionType string) (*models.LendingPosition, error) {
	position := &models.LendingPosition{
		UserID:       userID,
		Token:        token,
		Amount:       amount,
		Type:         positionType,
		Status:       "active",
		StartTime:    time.Now(),
		InterestRate: 0.05, // 示例利率
	}

	if err := s.db.Create(position).Error; err != nil {
		return nil, err
	}

	return position, nil
}

func (s *DefiService) GetUserPositions(userID uint) ([]models.LendingPosition, error) {
	var positions []models.LendingPosition
	if err := s.db.Where("user_id = ?", userID).Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

// 挖矿相关服务
func (s *DefiService) CreateFarmingPosition(userID uint, poolID uint, token string, amount float64) (*models.FarmingPosition, error) {
	position := &models.FarmingPosition{
		UserID:        userID,
		PoolID:        poolID,
		Token:         token,
		Amount:        amount,
		APY:           0.15, // 示例年化收益率
		StartTime:     time.Now(),
		LastClaimTime: time.Now(),
		Status:        "active",
	}

	if err := s.db.Create(position).Error; err != nil {
		return nil, err
	}

	return position, nil
}

func (s *DefiService) ClaimRewards(userID uint, positionID uint) (*models.Reward, error) {
	// 获取挖矿仓位
	var position models.FarmingPosition
	if err := s.db.First(&position, positionID).Error; err != nil {
		return nil, errors.New("position not found")
	}

	// 计算奖励
	timeSinceLastClaim := time.Since(position.LastClaimTime)
	rewardAmount := position.Amount * position.APY * float64(timeSinceLastClaim.Hours()) / (24 * 365)

	// 创建奖励记录
	reward := &models.Reward{
		UserID:     userID,
		PositionID: positionID,
		Token:      position.Token,
		Amount:     rewardAmount,
		Type:       "farming",
		ClaimTime:  time.Now(),
	}

	if err := s.db.Create(reward).Error; err != nil {
		return nil, err
	}

	// 更新最后领取时间
	position.LastClaimTime = time.Now()
	if err := s.db.Save(&position).Error; err != nil {
		return nil, err
	}

	return reward, nil
}

func (s *DefiService) GetUserRewards(userID uint) ([]models.Reward, error) {
	var rewards []models.Reward
	if err := s.db.Where("user_id = ?", userID).Find(&rewards).Error; err != nil {
		return nil, err
	}
	return rewards, nil
}
