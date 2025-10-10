package db

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ppt/log"
	"ppt/model"
	"time"
)

// BatchUpdateExpiredCoupons 使用游标批量更新过期优惠券(超大数据)
func BatchUpdateExpiredCoupons(db *gorm.DB, batchSize int) error {
	var cursorID uint64
	var cursorCouponID uuid.UUID
	currentTime := time.Now()

	log.Info("Begin BatchUpdateExpiredCoupons", zap.String("current_time", currentTime.Format(time.RFC3339)))

	for {
		var userCoupons []model.UserCoupon

		// 构建基础查询：查询未过期的记录
		query := db.Where("expired_at <= ? AND coupon_status = ?",
			currentTime.UnixMilli(), model.CouponStatusAvailable)

		// 如果有游标，添加游标条件（ID > 上一批最后一条记录的ID）
		if cursorID != 0 {
			query = query.Where("(user_id, coupon_id) > (?, ?)",
				cursorID, cursorCouponID)
		}

		// 按主键排序并限制批次大小
		result := query.Order("user_id, coupon_id").
			Limit(batchSize).
			Find(&userCoupons)

		if result.Error != nil {
			return result.Error
		}

		// 如果没有更多数据，退出循环
		if len(userCoupons) == 0 {
			break
		}

		// 批量更新当前批次的优惠券状态
		if err := updateCurrentBatch(db, userCoupons); err != nil {
			return err
		}

		// 更新游标位置
		lastCoupon := userCoupons[len(userCoupons)-1]
		cursorID = lastCoupon.UserID
		cursorCouponID = lastCoupon.CouponID

		log.Info("Current process info", zap.Int("cur_batch_size", len(userCoupons)), zap.Uint64("cur_cursor_user_id", cursorID), zap.String("cur_cursor_coupon_id", cursorCouponID.String()))
	}

	log.Info("End BatchUpdateExpiredCoupons")
	return nil
}

// updateCurrentBatch 更新当前批次的优惠券状态
func updateCurrentBatch(db *gorm.DB, coupons []model.UserCoupon) error {
	// 提取所有要更新的优惠券ID
	couponIDs := make([]struct {
		UserID   uint64
		CouponID uuid.UUID
	}, 0, len(coupons))

	for _, coupon := range coupons {
		couponIDs = append(couponIDs, struct {
			UserID   uint64
			CouponID uuid.UUID
		}{
			UserID:   coupon.UserID,
			CouponID: coupon.CouponID,
		})
	}

	// 执行批量更新
	result := db.Model(&model.UserCoupon{}).
		Where("(user_id, coupon_id) IN ?", couponIDs).
		Updates(map[string]interface{}{
			"coupon_status": model.CouponStatusExpired,
			"updated_at":    time.Now().UnixMilli(),
		})

	if result.Error != nil {
		return result.Error
	}

	log.Info("updateCurrentBatch success", zap.Int64("success_updated", result.RowsAffected))
	return nil
}

// DirectedUpdateExpiredCoupons 批量更新方法(数据量不大)
func DirectedUpdateExpiredCoupons(db *gorm.DB) error {
	currentTime := time.Now()

	log.Info("Begin DirectedUpdateExpiredCoupons", zap.String("current_time", currentTime.Format(time.RFC3339)))

	// 直接使用UPDATE语句批量更新，避免数据迁移
	result := db.Model(&model.UserCoupon{}).
		Where("expired_at <= ? AND coupon_status = ?",
			currentTime.UnixMilli(), model.CouponStatusAvailable).
		Updates(map[string]interface{}{
			"coupon_status": model.CouponStatusExpired,
			"updated_at":    currentTime.UnixMilli(),
		})

	if result.Error != nil {
		return result.Error
	}

	log.Info("End DirectedUpdateExpiredCoupons", zap.Int64("success_updated", result.RowsAffected))
	return nil
}
