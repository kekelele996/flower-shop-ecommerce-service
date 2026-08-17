package database

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// Seed 初始化种子数据：管理员、测试用户、分类、商品、优惠券。
func Seed(db *gorm.DB, minioClient *minio.Client, bucket string, logger *slog.Logger) error {
	var userCount int64
	if err := db.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		logger.Info(constants.LogDBSeedDone, "users", 0, "products", 0, "categories", 0)
		return nil
	}

	adminPwd, _ := util.HashPassword("admin123")
	userPwd, _ := util.HashPassword("user123")
	admin := &model.User{Username: "admin", Password: adminPwd, Nickname: "管理员", Email: "admin@flowershop.com", Phone: "13800000000", Role: constants.RoleAdmin, Status: constants.UserStatusActive}
	user := &model.User{Username: "user", Password: userPwd, Nickname: "花友小张", Email: "user@flowershop.com", Phone: "13900000000", Role: constants.RoleUser, Status: constants.UserStatusActive}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	categories := []model.Category{
		{Name: "鲜切花", Level: 1, Sort: 1},
		{Name: "盆栽绿植", Level: 1, Sort: 2},
		{Name: "多肉植物", Level: 1, Sort: 3},
		{Name: "花盆花器", Level: 1, Sort: 4},
		{Name: "园艺工具", Level: 1, Sort: 5},
		{Name: "营养土肥", Level: 1, Sort: 6},
	}
	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	// 生成占位图并上传 MinIO
	imgURL := func(name string, r, g, b uint8) string {
		data, err := util.PlaceholderPNG(600, 400, r, g, b)
		if err != nil {
			return ""
		}
		objectName := "seed/" + name + ".png"
		_, err = minioClient.PutObject(context.Background(), bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "image/png"})
		if err != nil {
			logger.Warn("seed image upload failed", "object", objectName, "err", err)
			return ""
		}
		return "/minio/" + bucket + "/" + objectName
	}

	products := []model.Product{
		{CategoryID: categories[0].ID, Name: "香槟玫瑰礼盒 20枝", SubTitle: "浪漫告白优选", Price: 129, OriginalPrice: 169, Stock: 200, Sales: 320, Rating: 4.9, RatingCount: 156, ShippingFrom: "云南昆明", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("rose", 230, 160, 150), Images: util.MarshalJSON([]string{imgURL("rose", 230, 160, 150), imgURL("rose2", 220, 150, 140)}), Description: "精选云南高原玫瑰，花头饱满，花期持久。", Detail: "<p>香槟玫瑰 20 枝，花语为『我只钟情你一个』。</p><p>顺丰冷链包邮，附保鲜剂与养护卡。</p>"},
		{CategoryID: categories[0].ID, Name: "向日葵花束 10枝", SubTitle: "阳光治愈系", Price: 59, OriginalPrice: 79, Stock: 300, Sales: 210, Rating: 4.8, RatingCount: 98, ShippingFrom: "云南昆明", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("sunflower", 240, 200, 80), Images: util.MarshalJSON([]string{imgURL("sunflower", 240, 200, 80)}), Description: "明亮向日葵，送给需要正能量的人。", Detail: "<p>向日葵 10 枝，花头直径约 10-12cm。</p>"},
		{CategoryID: categories[0].ID, Name: "白玫瑰 30枝 婚礼用花", SubTitle: "纯洁优雅", Price: 99, Stock: 150, Sales: 88, Rating: 4.7, RatingCount: 45, ShippingFrom: "云南昆明", FreeShipping: false, Status: constants.ProductStatusOnSale, CoverImage: imgURL("whiterose", 240, 240, 240), Images: util.MarshalJSON([]string{imgURL("whiterose", 240, 240, 240)}), Description: "白玫瑰 30 枝，适合婚礼与纪念日。", Detail: "<p>厄瓜多尔白玫瑰同款种植，A 级品质。</p>"},
		{CategoryID: categories[1].ID, Name: "琴叶榕 1.2米 大盆栽", SubTitle: "北欧风客厅绿植", Price: 199, OriginalPrice: 259, Stock: 80, Sales: 132, Rating: 4.8, RatingCount: 76, ShippingFrom: "广东广州", FreeShipping: false, Status: constants.ProductStatusOnSale, CoverImage: imgURL("fiddle", 90, 140, 90), Images: util.MarshalJSON([]string{imgURL("fiddle", 90, 140, 90)}), Description: "琴叶榕盆栽，含陶瓷盆，净化空气。", Detail: "<p>株高约 1.2m，含盆发货，适合客厅/办公室。</p>"},
		{CategoryID: categories[1].ID, Name: "龟背竹 室内绿植", SubTitle: "ins风拍照绿植", Price: 69, Stock: 200, Sales: 410, Rating: 4.9, RatingCount: 233, ShippingFrom: "广东广州", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("monstera", 60, 120, 70), Images: util.MarshalJSON([]string{imgURL("monstera", 60, 120, 70)}), Description: "龟背竹盆栽，叶形独特，好养耐阴。", Detail: "<p>含营养土与盆，办公室/卧室皆宜。</p>"},
		{CategoryID: categories[1].ID, Name: "绿萝 水培盆栽", SubTitle: "新手必入", Price: 29, Stock: 500, Sales: 680, Rating: 4.7, RatingCount: 342, ShippingFrom: "福建漳州", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("pothos", 80, 160, 80), Images: util.MarshalJSON([]string{imgURL("pothos", 80, 160, 80)}), Description: "绿萝水培，好养易活，净化空气。", Detail: "<p>水培绿萝，玻璃瓶装。</p>"},
		{CategoryID: categories[2].ID, Name: "多肉植物组合 10颗", SubTitle: "办公桌治愈系", Price: 39, OriginalPrice: 49, Stock: 260, Sales: 520, Rating: 4.8, RatingCount: 289, ShippingFrom: "山东青州", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("succulent", 150, 100, 120), Images: util.MarshalJSON([]string{imgURL("succulent", 150, 100, 120)}), Description: "随机 10 颗景天科多肉，含盆土。", Detail: "<p>含拇指盆与颗粒土，新手友好。</p>"},
		{CategoryID: categories[2].ID, Name: "桃蛋 老桩多肉", SubTitle: "粉粉圆滚滚", Price: 25, Stock: 180, Sales: 190, Rating: 4.6, RatingCount: 87, ShippingFrom: "山东青州", FreeShipping: false, Status: constants.ProductStatusOnSale, CoverImage: imgURL("peach", 220, 120, 160), Images: util.MarshalJSON([]string{imgURL("peach", 220, 120, 160)}), Description: "桃蛋老桩，状态粉嫩，带盆发货。", Detail: "<p>单头直径 3-4cm，精品桩。</p>"},
		{CategoryID: categories[3].ID, Name: "北欧陶瓷花盆 白色 大口径", SubTitle: "百搭简约", Price: 45, Stock: 400, Sales: 300, Rating: 4.7, RatingCount: 154, ShippingFrom: "江西景德镇", FreeShipping: false, Status: constants.ProductStatusOnSale, CoverImage: imgURL("pot", 230, 230, 235), Images: util.MarshalJSON([]string{imgURL("pot", 230, 230, 235)}), Description: "简约白陶瓷盆，口径 18cm，带托盆。", Detail: "<p>釉面光滑，底部带孔。</p>"},
		{CategoryID: categories[3].ID, Name: "水泥花盆 三件套", SubTitle: "工业风组合", Price: 59, Stock: 220, Sales: 145, Rating: 4.5, RatingCount: 62, ShippingFrom: "浙江义乌", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("cementpot", 150, 150, 150), Images: util.MarshalJSON([]string{imgURL("cementpot", 150, 150, 150)}), Description: "水泥灰花盆三件套，适合多肉/小型绿植。", Detail: "<p>口径 8/10/12cm 各一。</p>"},
		{CategoryID: categories[4].ID, Name: "园艺三件套 铲子+耙子+喷壶", SubTitle: "家庭园艺必备", Price: 35, Stock: 350, Sales: 260, Rating: 4.6, RatingCount: 118, ShippingFrom: "浙江永康", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("tools", 120, 90, 60), Images: util.MarshalJSON([]string{imgURL("tools", 120, 90, 60)}), Description: "不锈钢园艺工具三件套，含 500ml 喷壶。", Detail: "<p>防锈材质，握感舒适。</p>"},
		{CategoryID: categories[4].ID, Name: "电动修枝剪 锂电", SubTitle: "省力园艺神器", Price: 159, OriginalPrice: 199, Stock: 60, Sales: 70, Rating: 4.4, RatingCount: 31, ShippingFrom: "浙江永康", FreeShipping: false, Status: constants.ProductStatusOnSale, CoverImage: imgURL("pruner", 90, 90, 60), Images: util.MarshalJSON([]string{imgURL("pruner", 90, 90, 60)}), Description: "锂电修枝剪，一剪两用，续航 8 小时。", Detail: "<p>Type-C 充电，含备用刀片。</p>"},
		{CategoryID: categories[5].ID, Name: "有机营养土 10L", SubTitle: "通用型配方", Price: 29, Stock: 600, Sales: 730, Rating: 4.8, RatingCount: 402, ShippingFrom: "江苏宿迁", FreeShipping: true, Status: constants.ProductStatusOnSale, CoverImage: imgURL("soil", 100, 70, 50), Images: util.MarshalJSON([]string{imgURL("soil", 100, 70, 50)}), Description: "通用营养土 10L，透气保肥。", Detail: "<p>含泥炭、珍珠岩、有机肥。</p>"},
	}
	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			return err
		}
	}

	// 给测试用户发放两张优惠券
	now := time.Now()
	coupons := []model.Coupon{
		{UserID: user.ID, TemplateID: 1, Name: "新人满减券", Type: constants.CouponTypeFullReduction, Threshold: 100, Amount: 20, Status: constants.CouponStatusUnused, ValidFrom: now.AddDate(0, 0, -1), ValidTo: now.AddDate(0, 0, 29)},
		{UserID: user.ID, TemplateID: 2, Name: "9折花卉券", Type: constants.CouponTypeDiscount, Threshold: 50, DiscountRate: 0.9, Status: constants.CouponStatusUnused, ValidFrom: now.AddDate(0, 0, -1), ValidTo: now.AddDate(0, 0, 29)},
	}
	if err := db.Create(&coupons).Error; err != nil {
		return err
	}

	logger.Info(constants.LogDBSeedDone, "users", 2, "products", len(products), "categories", len(categories))
	fmt.Println("seed done")
	return nil
}
