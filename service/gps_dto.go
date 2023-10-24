package service

type GpsDto struct {
	RouteId           int64   `json:"route_id" form:"route_id"`
	UserId            string  `json:"user_id" form:"user_id"`
	AddToUserId       int64   `json:"add_to_user_id" form:"add_to_user_id"`
	Location          string  `json:"location" form:"location"`
	SportsType        int     `json:"sports_type" form:"sports_type"` // 0: outdoor 1: inner
	StartTime         string  `json:"start_time" form:"start_time"`
	EndTime           string  `json:"end_time" form:"end_time"`
	TotalLength       float64 `json:"total_length" form:"total_length"`
	TotalCalories     float64 `json:"total_calories" form:"total_calories"`
	PointsStr         string  `json:"points_str" form:"points_str"`
	ProductId         string  `json:"product_id" form:"product_id"`
	AveragePace       int64   `json:"average_pace" form:"average_pace"`
	AverageSpeed      int64   `json:"average_speed" form:"average_speed"`
	HighestSpeedPerkm int64   `json:"highest_speed_perkm" form:"highest_speed_perkm"`
	PacePerM          string  `json:"pace_per_m" form:"pace_per_m"`
	PacePerMile       string  `json:"pace_per_mile" form:"pace_per_mile"`
	TotalTime         float64 `json:"total_time" form:"total_time"`
	Pause             int32   `json:"pause" form:"pause"`
	Cadences          string  `json:"cadences" form:"cadences"`
	Steps             int64   `json:"steps" form:"steps"`
	HideMap           int     `json:"hide_map" form:"hide_map"`
	HideKmCard        int     `json:"hide_km_card" form:"hide_km_card"`
	LocusUrl          string  `json:"locus_url" form:"locus_url"`   // 轨迹图（方）
	LocusUrl2         string  `json:"locus_url2" form:"locus_url2"` // 轨迹图（长）
	IsHistory         int     `json:"is_history" form:"is_history"` // 是否是第三方的历史数据
	GoalType          int     `json:"goal_type" form:"goal_type"`   // 0 未选择  1 distance  2 time
	GoalValue         float64 `json:"goal_value" form:"goal_value"`
	GoalResult        int     `json:"goal_result" form:"goal_result"` // 1未完成  2完成
	HeartRate         string  `json:"heart_rate" form:"heart_rate"`
	SourceType        int     `json:"source_type" form:"source_type" format:"0: 原生(andriod,ios) 1: 第三方接入 2：runtopia shoes 3: 混合（shoes+App）4:codoon watch s1"`
	Os                string  `json:"os" form:"os"`
	AppVersion        string  `json:"app_version" form:"app_version"`
	Device            string  `json:"device" form:"device"`
	PointNum          int     `json:"point_num" form:"point_num"`
	StartPoint        string  `json:"start_point" form:"start_point"`
	EndPoint          string  `json:"end_point" form:"end_point"`
	UploadTime        int64   `json:"upload_time" form:"upload_time"`
	GPSType           int64   `json:"gps_type" form:"gps_type"`
	//RUNBOX  存储数据文件路径
	FileUrl string `json:"file_url" form:"file_url"`
}

type GetGpsReq struct {
	UserId  string `json:"user_id" form:"user_id"`
	RouteId string `json:"route_id"form:"route_id"`
}

type TestReq struct {
	UserId  string `json:"user_id" form:"user_id"`
	RouteId string `json:"route_id" form:"route_id"`
}

type TestRsp struct {
	UserId  string   `json:"user_id"`
	RouteId string   `json:"route_id"`
	LL      []string `json:"ll"`
	FF      int      `json:"ff"`
}

type PostGpsRspData struct {
	RouteId  string      `json:"route_id"`
	Medals   interface{} `json:"medals"`
	Grade    int         `json:"grade"`
	OldGrade int         `json:"old_grade"`
}
