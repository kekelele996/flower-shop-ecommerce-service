package service

import (
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/repository"
)

// StatsService 销售统计服务。
type StatsService struct {
	orderRepo *repository.OrderRepository
	prodRepo  *repository.ProductRepository
	userRepo  *repository.UserRepository
}

func NewStatsService(orderRepo *repository.OrderRepository, prodRepo *repository.ProductRepository, userRepo *repository.UserRepository) *StatsService {
	return &StatsService{orderRepo: orderRepo, prodRepo: prodRepo, userRepo: userRepo}
}

func (s *StatsService) Stats() (*dto.StatsVO, error) {
	totalSales, err := s.orderRepo.SumSales()
	if err != nil {
		return nil, err
	}
	orderCount, err := s.orderRepo.CountAll()
	if err != nil {
		return nil, err
	}
	productCount, err := s.prodRepo.CountAll()
	if err != nil {
		return nil, err
	}
	userCount, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}
	pendingShipment, err := s.orderRepo.CountByStatus("PENDING_SHIPMENT")
	if err != nil {
		return nil, err
	}
	statusMap, err := s.orderRepo.CountByStatusMap()
	if err != nil {
		return nil, err
	}
	topProducts, err := s.orderRepo.TopProductsByRevenue(5)
	if err != nil {
		return nil, err
	}
	dailySales, err := s.orderRepo.DailySales(7)
	if err != nil {
		return nil, err
	}
	return &dto.StatsVO{
		TotalSales:       totalSales,
		OrderCount:       orderCount,
		ProductCount:     productCount,
		UserCount:        userCount,
		PendingShipment:  pendingShipment,
		OrderStatusCount: statusMap,
		TopProducts:      topProducts,
		DailySales:       dailySales,
	}, nil
}
