package activity

type ActivityGoods struct {
	ID             int32     `gorm:column:id json:"id"`
	Name           string    `gorm:column:name json:"name"`
	UniqueCode     string    `gorm:column:unique_code json:"unique_code"`
	Title          string    `gorm:column:title json:"title"`
	Type           int8      `gorm:column:type json:"type"`
	ActType        int32     `gorm:column:act_type json:"act_type"`
	StartTime      time.Time `gorm:column:start_time json:"start_time"`
	EndTime        time.Time `gorm:column:end_time json:"end_time"`
	IsDeleted      int8      `gorm:column:is_deleted json:"is_deleted"`
	IndexShow      int8      `gorm:column:index_show json:"index_show"`
	Status         int8      `gorm:column:status json:"status"`
	Banner         string    `gorm:column:banner json:"banner"`
	Layout         int8      `gorm:column:layout json:"layout"`
	Config         string    `gorm:column:config json:"config"`
	CreateAt       time.Time `gorm:column:create_at json:"create_at"`
	UpdateAt       time.Time `gorm:column:update_at json:"update_at"`
	Watermarker    string    `gorm:column:watermarker json:"watermarker"`
	SortType       int8      `gorm:column:sort_type json:"sort_type"`
	GroupId        int32     `gorm:column:group_id json:"group_id"`
	UserId         int32     `gorm:column:user_id json:"user_id"`
	Remark         string    `gorm:column:remark json:"remark"`
	LockStartTime  time.Time `gorm:column:lock_start_time json:"lock_start_time"`
	LockEndTime    time.Time `gorm:column:lock_end_time json:"lock_end_time"`
	Rank           int32     `gorm:column:rank json:"rank"`
	CategoryId     string    `gorm:column:category_id json:"category_id"`
	GoodsType      int8      `gorm:column:goods_type json:"goods_type"`
	MtCategoryId   string    `gorm:column:mt_category_id json:"mt_category_id"`
	GoodsSource    int8      `gorm:column:goods_source json:"goods_source"`
	BrandId        int32     `gorm:column:brand_id json:"brand_id"`
	SystemDiscount int8      `gorm:column:system_discount json:"system_discount"`
	CoinDiscount   int8      `gorm:column:coin_discount json:"coin_discount"`
}

func (ActivityGoods) TableName() string {
	return "activity_goods"
}
