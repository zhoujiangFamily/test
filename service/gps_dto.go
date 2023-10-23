package service

type GpsDto struct {
	RouteId           int64   `json:"route_id"`
	UserId            string  `json:"user_id"`
	AddToUserId       int64   `json:"add_to_user_id"`
	Location          string  `json:"location"`
	SportsType        int     `json:"sports_type"` // 0: outdoor 1: inner
	StartTime         string  `json:"start_time"`
	EndTime           string  `json:"end_time"`
	TotalLength       float64 `json:"total_length"`
	TotalCalories     float64 `json:"total_calories"`
	PointsStr         string  `json:"points_str"`
	ProductId         string  `json:"product_id"`
	AveragePace       int64   `json:"average_pace"`
	AverageSpeed      int64   `json:"average_speed"`
	HighestSpeedPerkm int64   `json:"highest_speed_perkm"`
	PacePerM          string  `json:"pace_per_m"`
	PacePerMile       string  `json:"pace_per_mile"`
	TotalTime         float64 `json:"total_time"`
	Pause             int32   `json:"pause"`
	Cadences          string  `json:"cadences"`
	Steps             int64   `json:"steps"`
	HideMap           int     `json:"hide_map"`
	HideKmCard        int     `json:"hide_km_card"`
	LocusUrl          string  `json:"locus_url"`  // 轨迹图（方）
	LocusUrl2         string  `json:"locus_url2"` // 轨迹图（长）
	IsHistory         int     `json:"is_history"` // 是否是第三方的历史数据
	GoalType          int     `json:"goal_type"`  // 0 未选择  1 distance  2 time
	GoalValue         float64 `json:"goal_value"`
	GoalResult        int     `json:"goal_result"` // 1未完成  2完成
	HeartRate         string  `json:"heart_rate"`
	SourceType        int     `json:"source_type" format:"0: 原生(andriod,ios) 1: 第三方接入 2：runtopia shoes 3: 混合（shoes+App）4:codoon watch s1"`
	Os                string  `json:"os"`
	AppVersion        string  `json:"app_version"`
	Device            string  `json:"device"`
	PointNum          int     `json:"point_num"`
	StartPoint        string  `json:"start_point"`
	EndPoint          string  `json:"end_point"`
	UploadTime        int64   `json:"upload_time"`
	GPSType           int64   `json:"gps_type"`
	//RUNBOX  存储数据文件路径
	FileUrl string `json:"file_url"`
}

type GetGpsReq struct {
	UserId  string `json:"user_id"`
	RouteId string `json:"route_id"`
}

type TestReq struct {
	UserId  string `json:"user_id"`
	RouteId string `json:"route_id"`
}

type TestRsp struct {
	UserId  string   `json:"user_id"`
	RouteId string   `json:"route_id"`
	LL      []string `json:"ll"`
	FF      int      `json:"ff"`
}
