package model

import (
	"fmt"
	"git.in.codoon.com/Overseas/runbox/first-test/conf"
	"log"
	"time"
)

type Gps struct {
	Id                int       `json:"id"`
	RouteId           string    `json:"route_id"`
	UserId            string    `json:"user_id"`
	AddToUserId       int64     `json:"add_to_user_id"`
	Location          string    `json:"location"`
	SportsType        int       `json:"sports_type"` // 0: outdoor 1: inner
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	TotalLength       float64   `json:"total_length"`
	TotalCalories     float64   `json:"total_calories"`
	PointsStr         string    `json:"points_str"`
	ProductId         string    `json:"product_id"`
	AveragePace       int64     `json:"average_pace"`
	AverageSpeed      int64     `json:"average_speed"`
	HighestSpeedPerkm int64     `json:"highest_speed_perkm"`
	PacePerM          string    `json:"pace_per_m"`
	PacePerMile       string    `json:"pace_per_mile"`
	TotalTime         float64   `json:"total_time"`
	Pause             int32     `json:"pause"`
	Cadences          string    `json:"cadences"`
	Steps             int64     `json:"steps"`
	HideMap           int       `json:"hide_map"`
	HideKmCard        int       `json:"hide_km_card"`
	LocusUrl          string    `json:"locus_url"`  // 轨迹图（方）
	LocusUrl2         string    `json:"locus_url2"` // 轨迹图（长）
	IsHistory         int       `json:"is_history"` // 是否是第三方的历史数据
	GoalType          int       `json:"goal_type"`  // 0 未选择  1 distance  2 time
	GoalValue         float64   `json:"goal_value"`
	GoalResult        int       `json:"goal_result"` // 1未完成  2完成
	HeartRate         string    `json:"heart_rate"`
	SourceType        int       `json:"source_type" format:"0: 原生(andriod,ios) 1: 第三方接入 2：runtopia shoes 3: 混合（shoes+App）4:codoon watch s1"`
	Os                string    `json:"os"`
	AppVersion        string    `json:"app_version"`
	Device            string    `json:"device"`
	PointNum          int       `json:"point_num"`
	StartPoint        string    `json:"start_point"`
	EndPoint          string    `json:"end_point"`

	UploadTime time.Time `json:"upload_time"`

	GPSType int64  `json:"gps_type"`
	FileUrl string `json:"file_url"`
	Md5     string `json:"md_5"`
}

func (g *Gps) Create() error {

	insertStmt, err := conf.Fb_mysql.Prepare("INSERT INTO gps_route_data (route_id,user_id,total_length,total_time,total_calories,location,sports_type,start_time, EndTime,upload_time,locus_url,locus_url2,steps,file_url) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Printf("create gps Prepare error :%v", err)
		return err
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(g.RouteId, g.UserId, g.TotalLength, g.TotalTime, g.TotalCalories, g.Location, g.SportsType, g.StartTime, g.EndTime, g.UploadTime, g.LocusUrl, g.LocusUrl2, g.Steps, g.FileUrl)
	if err != nil {
		log.Printf("create gps Exec error :%v", err)

		return err
	}

	return nil

}

func (g *Gps) Update() error {
	fmt.Print("更新数据 ")
	stmt, err := conf.Fb_mysql.Prepare("update gps_route_data set user_id=? where route_id=?")
	if err != nil {
		log.Printf("Update Gps Prepare error :%v", err)
		return err

	}
	res, err := stmt.Exec(g.UserId, g.RouteId)
	if err != nil {
		log.Printf("Update Gps Exec error :%v", err)
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Update Gps RowsAffected error :%v", err)
		return err

	}
	return nil

}

func (g *Gps) Select(routeId string) error {

	insertStmt, err := conf.Fb_mysql.Prepare("select id,route_id,user_id,total_length,total_time,total_calories,location,sports_type,start_time, EndTime,upload_time,locus_url,locus_url2,steps,file_url from gps_route_data where route_id = ?")
	if err != nil {
		log.Printf("Select Gps select error :%v", err)
		return err
	}
	defer insertStmt.Close()
	rows, err := insertStmt.Query(routeId)
	if err != nil {
		log.Printf("Select Gps select error :%v", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&g.Id,
			&g.RouteId,
			&g.UserId,
			&g.TotalLength,
			&g.TotalTime,
			&g.TotalCalories,
			&g.Location,
			&g.SportsType,
			&g.StartTime,
			&g.EndTime,
			&g.UploadTime,
			&g.LocusUrl,
			&g.LocusUrl2,
			&g.Steps,
			&g.FileUrl,
		)

		if err != nil {
			log.Printf("Select Gps select error :%v", err)
			return err
		}
	}
	if err = rows.Err(); err != nil {
		log.Printf("Select Gps select error :%v", err)
		return err
	}
	return nil

}
