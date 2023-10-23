package service

import (
	"database/sql"
	"fmt"
	"git.in.codoon.com/Overseas/runbox/first-test/common"
	"git.in.codoon.com/Overseas/runbox/first-test/conf"
	"git.in.codoon.com/Overseas/runbox/first-test/model"
	"html/template"
	"log"
	"math"
	"net/http"
	"time"
)

//一个模版
var (
	indexTmpl = template.Must(template.New("index").Parse(indexHTML))
)

// Votes handles HTTP requests to alternatively show the voting app or to save a
// vote.
//Votes处理HTTP请求，以显示投票应用程序或保存
//应该是程序入口
func Votes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderIndex(w, r, conf.Fb_mysql)
	case http.MethodPost:
		saveVote(w, r, conf.Fb_mysql)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func Gps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		gpsget(w, r, conf.Fb_mysql)
	case http.MethodPost:
		gpsPost(w, r, conf.Fb_mysql)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		TestGet(w, r)
	case http.MethodPost:
		TestPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// vote contains a single row from the votes table in the database. Each vote
// includes a candidate ("TABS" or "SPACES") and a timestamp.
//投票的数据结构
type vote struct {
	//候选人
	Candidate string
	//投票时间
	VoteTime time.Time
}

// renderIndex renders the HTML application with the voting form, current
// totals, and recent votes.

//以html的 方式呈现最近的投票结果
func renderIndex(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	t, err := currentTotals(db)
	if err != nil {
		log.Printf("renderIndex: failed to read current totals: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = indexTmpl.Execute(w, t)
	if err != nil {
		log.Printf("renderIndex: failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// saveVote saves a vote passed as http.Request form data.
// 保存投票
func saveVote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err := r.ParseForm(); err != nil {
		log.Printf("saveVote: failed to parse form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	team := r.FormValue("team")
	if team == "" {
		log.Printf("saveVote: \"team\" property missing from form submission")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if team != "TABS" && team != "SPACES" {
		log.Printf("saveVote: \"team\" property should be \"TABS\" or \"SPACES\", was %q", team)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// [START cloud_sql_mysql_databasesql_connection]
	insertVote := "INSERT INTO votes(candidate, created_at) VALUES(?, NOW())"
	_, err := db.Exec(insertVote, team)
	// [END cloud_sql_mysql_databasesql_connection]

	if err != nil {
		log.Printf("saveVote: unable to save vote: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Vote successfully cast for %s!", team)
}

// recentVotes returns the last five votes cast.
// 查询语句  查询投票
func recentVotes(db *sql.DB) ([]vote, error) {
	rows, err := db.Query("SELECT candidate, created_at FROM votes ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		return nil, fmt.Errorf("DB.Query: %w", err)
	}
	defer rows.Close()

	var votes []vote
	for rows.Next() {
		var (
			candidate string
			voteTime  time.Time
		)
		err := rows.Scan(&candidate, &voteTime)
		if err != nil {
			return nil, fmt.Errorf("Rows.Scan: %w", err)
		}
		votes = append(votes, vote{Candidate: candidate, VoteTime: voteTime})
	}
	return votes, nil
}

// formatMargin calculates the difference between votes and returns a human
// friendly margin (e.g., 2 votes)
//计算选票之间的差额并返回一个人
func formatMargin(a, b int) string {
	diff := int(math.Abs(float64(a - b)))
	margin := fmt.Sprintf("%d votes", diff)
	// remove pluralization when diff is just one
	if diff == 1 {
		margin = "1 vote"
	}
	return margin
}

// votingData is used to pass data to the HTML template.
//将数据传递到html模版
type votingData struct {
	TabsCount   int
	SpacesCount int
	VoteMargin  string
	RecentVotes []vote
}

// currentTotals retrieves all voting data from the database.
//从数据库中检索所有投票数据。
func currentTotals(db *sql.DB) (votingData, error) {
	var (
		tabs   int
		spaces int
	)
	err := db.QueryRow("SELECT count(id) FROM votes WHERE candidate='TABS'").Scan(&tabs)
	if err != nil {
		return votingData{}, fmt.Errorf("DB.QueryRow: %w", err)
	}
	err = db.QueryRow("SELECT count(id) FROM votes WHERE candidate='SPACES'").Scan(&spaces)
	if err != nil {
		return votingData{}, fmt.Errorf("DB.QueryRow: %w", err)
	}

	recent, err := recentVotes(db)
	if err != nil {
		return votingData{}, fmt.Errorf("recentVotes: %w", err)
	}

	return votingData{
		TabsCount:   tabs,
		SpacesCount: spaces,
		VoteMargin:  formatMargin(tabs, spaces),
		RecentVotes: recent,
	}, nil
}

var indexHTML = `
<html lang="en">
<head>
    <title>Tabs VS Spaces</title>
    <link rel="icon" type="image/png" href="data:image/png;base64,iVBORw0KGgo=">
    <link rel="stylesheet"
          href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
</head>
<body>
<nav class="red lighten-1">
    <div class="nav-wrapper">
        <a href="#" class="brand-logo center">Tabs VS Spaces</a>
    </div>
</nav>
<div class="section">
    <div class="center">
        <h4>
            {{ if eq .TabsCount .SpacesCount }}
                TABS and SPACES are evenly matched!
            {{ else if gt .TabsCount .SpacesCount }}
                TABS are winning by {{ .VoteMargin }}
            {{ else if gt .SpacesCount .TabsCount }}
                SPACES are winning by {{ .VoteMargin }}
            {{ end }}
        </h4>
    </div>
    <div class="row center">
        <div class="col s6 m5 offset-m1">
            {{ if gt .TabsCount .SpacesCount }}
			<div class="card-panel green lighten-3">
			{{ else }}
			<div class="card-panel">
			{{ end }}
                <i class="material-icons large">keyboard_tab</i>
                <h3>{{ .TabsCount }} votes</h3>
                <button id="voteTabs" class="btn green">Vote for TABS</button>
            </div>
        </div>
        <div class="col s6 m5">
            {{ if lt .TabsCount .SpacesCount }}
			<div class="card-panel blue lighten-3">
			{{ else }}
			<div class="card-panel">
			{{ end }}
                <i class="material-icons large">space_bar</i>
                <h3>{{ .SpacesCount }} votes</h3>
                <button id="voteSpaces" class="btn blue">Vote for SPACES</button>
            </div>
        </div>
    </div>
    <h4 class="header center">Recent Votes</h4>
    <ul class="container collection center">
        {{ range .RecentVotes }}
            <li class="collection-item avatar">
                {{ if eq .Candidate "TABS" }}
                    <i class="material-icons circle green">keyboard_tab</i>
                {{ else if eq .Candidate "SPACES" }}
                    <i class="material-icons circle blue">space_bar</i>
                {{ end }}
                <span class="title">
                    A vote for <b>{{.Candidate}}</b> was cast at {{.VoteTime.Format "2006-01-02T15:04:05Z07:00" }}
                </span>
            </li>
        {{ end }}
    </ul>
</div>
<script>
    function vote(team) {
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function () {
            if (this.readyState == 4) {
                window.location.reload();
            }
        };
        xhr.open("POST", "/Votes", true);
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xhr.send("team=" + team);
    }

    document.getElementById("voteTabs").addEventListener("click", function () {
        vote("TABS");
    });
    document.getElementById("voteSpaces").addEventListener("click", function () {
        vote("SPACES");
    });
</script>
</body>
</html>
`

func TestGet(w http.ResponseWriter, r *http.Request) {
	req := &GetGpsReq{}
	booll := common.Bind(r, req)

	if !booll {
		log.Printf("bing request failed parse form: %v ", r)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	vfvf := make([]string, 0)
	vfvf = append(vfvf, "dsfds")
	vfvf = append(vfvf, "dsfds")
	data := TestRsp{
		UserId:  req.UserId,
		RouteId: req.RouteId,
		LL:      vfvf,
		FF:      1,
	}

	common.Render(w, 200, data)
}
func TestPost(w http.ResponseWriter, r *http.Request) {
	req := &GetGpsReq{}
	booll := common.Bind(r, req)

	if !booll {
		log.Printf("bing request failed parse form: %v  xxx %v ", r.Body, r.Form)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("bing request failed parse form: %v ", req)

	vfvf := make([]string, 0)
	vfvf = append(vfvf, "dsfds")
	vfvf = append(vfvf, "dsfds")
	data := TestRsp{
		UserId:  req.UserId,
		RouteId: req.RouteId,
		LL:      vfvf,
		FF:      1,
	}

	common.Render(w, 200, data)
}

//获取一条记录
func gpsget(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	///
	req := &GetGpsReq{}
	booll := common.Bind(r, req)

	if !booll {
		log.Printf("bing request failed parse form: %v ", r)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//根据route_id 和 获取路线

	data := model.Gps{}
	err := data.Select(req.RouteId)
	if err != nil {
		log.Printf("gpsGet: failed req : %v", req)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "hello Go Web get")
	common.Render(w, 200, data)
}

//获取一条记录
func GetGpsList(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	///
	req := &GetGpsReq{}
	booll := common.Bind(r, req)

	if !booll {
		log.Printf("bing request failed parse form: %v ", r)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//根据route_id 和 获取路线

	data := model.Gps{}
	err := data.Select(req.RouteId)
	if err != nil {
		log.Printf("gpsGet: failed req : %v", req)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "hello Go Web get")
	common.Render(w, 200, data)
}

func gpsPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if err := r.ParseForm(); err != nil {
		log.Printf("gpsPost: failed to parse form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//
	userId := common.GetUserId(r)

	if userId == "" {
		log.Printf("gpsPost: userId is empty ")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	gpsDto := &GpsDto{
		UserId: userId,
	}
	booll := common.Bind(r, gpsDto)

	if !booll {
		log.Printf("bing request failed parse form: %v ", r)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Upload_time := time.Now()
	//
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gpsDto.StartTime, time.Now().Location())
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", gpsDto.StartTime, time.Now().Location())

	//转换数据

	routeId := common.GetRandomRouteId()
	gps := model.Gps{
		UserId:        userId,
		RouteId:       routeId,
		AddToUserId:   gpsDto.AddToUserId,
		Location:      gpsDto.Location,
		SportsType:    gpsDto.SportsType,
		StartTime:     startTime,
		EndTime:       endTime,
		TotalLength:   gpsDto.TotalLength,
		TotalCalories: gpsDto.TotalCalories,
		ProductId:     gpsDto.ProductId,

		AveragePace:       gpsDto.AveragePace,
		AverageSpeed:      gpsDto.AverageSpeed,
		HighestSpeedPerkm: gpsDto.HighestSpeedPerkm,

		PacePerM: gpsDto.PacePerM,

		PacePerMile: gpsDto.PacePerMile,
		TotalTime:   gpsDto.TotalTime,
		Pause:       gpsDto.Pause,
		Cadences:    gpsDto.Cadences,
		Steps:       gpsDto.Steps,
		HideMap:     gpsDto.HideMap,
		HideKmCard:  gpsDto.HideKmCard,
		LocusUrl:    gpsDto.LocusUrl,
		LocusUrl2:   gpsDto.LocusUrl2,
		IsHistory:   gpsDto.IsHistory,
		GoalType:    gpsDto.GoalType,
		GoalResult:  gpsDto.GoalResult,
		GoalValue:   gpsDto.GoalValue,
		HeartRate:   gpsDto.HeartRate,
		SourceType:  gpsDto.SourceType,
		Os:          gpsDto.Os,
		AppVersion:  gpsDto.AppVersion,
		PointNum:    gpsDto.PointNum,
		StartPoint:  gpsDto.StartPoint,
		EndPoint:    gpsDto.EndPoint,
		UploadTime:  Upload_time,
	}

	//创建数据
	err := gps.Create()

	if err != nil {
		log.Printf("gpsPost: failed to parse form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//fmt.Fprintf(w, "hello Go Web Post")
	common.Render(w, 200, gps.RouteId)
}
