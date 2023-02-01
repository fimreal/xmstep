package main

import (
	"net/http"
	"os"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/gin-gonic/gin"
)

var port = ":" + os.Getenv("PORT")

func main() {
	// ezap.SetLevel("debug")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/run", run)
	r.GET("/run", run)
	ezap.Info("listening to ", port)
	r.Run(port)
}

func run(ctx *gin.Context) {
	type RunReq struct {
		Username string `json:"username" form:"username" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
		Step     int    `json:"step" form:"step" validate:"required"`
		Date     string `json:"date" form:"date"`
	}

	var req RunReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		ezap.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var (
		username = req.Username
		password = req.Password
		step     = req.Step
		date     = req.Date
	)
	ezap.Infof("准备设置步数 username: %s, password: %s, step: %d, date: %s", username, password, step, date)

	// 检查是否自定义时间
	t := time.Now()
	if date != "" {
		t, err = time.Parse("2006-01-02 15:04:05  -0700 MST", req.Date+" +0800 CST")
		if err != nil {
			ezap.Error("输入时间格式不正确: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "输入时间格式不正确, 正确格式例如2022-01-01 19:00:00"})
			return
		}
	}

	// 获取账号 token 等
	account, err := getAccount(req.Username, req.Password)
	if err != nil {
		ezap.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置时间
	err = account.set(step, t)
	if err != nil {
		ezap.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ezap.Infof("%s 成功设置步数 %d", username, step)
	ctx.JSON(http.StatusOK, gin.H{"username": username, "step": step, "time": t.String(), "result": "success"})

}
