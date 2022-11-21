package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/fimreal/goutils/parse"
	"github.com/gin-gonic/gin"
)

type Account struct {
	TokenInfo struct {
		LoginToken string `json:"login_token"`
		AppToken   string `json:"app_token"`
		UserID     string `json:"user_id"`
		TTL        int    `json:"ttl"`
		AppTTL     int    `json:"app_ttl"`
	} `json:"token_info"`
	RegistInfo struct {
		IsNewUser   int    `json:"is_new_user"`
		RegistDate  int64  `json:"regist_date"`
		Region      string `json:"region"`
		CountryCode string `json:"country_code"`
	} `json:"regist_info"`
	ThirdpartyInfo struct {
		Nickname string `json:"nickname"`
		Icon     string `json:"icon"`
		ThirdID  string `json:"third_id"`
		Email    string `json:"email"`
	} `json:"thirdparty_info"`
	Result string `json:"result"`
	Domain struct {
		IDDNS string `json:"id-dns"`
	} `json:"domain"`
	Domains []interface{} `json:"domains"`
}

var (
	headers = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
		"User-Agent":   "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2",
	}
	usernameType = "email"
)

// return token, err
func getAccount(username, password string) (*Account, error) {
	account := &Account{}

	if parse.IsPhoneNumber(username) {
		username = "+86" + username
		usernameType = "huami_phone"
	} else if !parse.IsEmail(username) {
		return account, errors.New("invalid username")
	}

	loginUrl := "https://api-user.huami.com/registrations/" + username + "/tokens"
	postData := url.Values{
		"client_id":    {"HuaMi"},
		"country_code": {"CN"},
		"password":     {password},
		"redirect_uri": {"https://s3-us-west-2.amazonaws.com/hm-registration/successsignin.html"},
		"state":        {"REDIRECTION"},
		"token":        {"access"},
	}.Encode()
	headers["apptoken"] = ""
	res, err := HttpPost(loginUrl, postData, headers)
	if err != nil {
		return account, err
	}
	location, _ := res.Location()
	access := location.Query().Get("access")
	ezap.Debug("获取到登录 token: ", access)

	accountUrl := "https://account.huami.com/v2/client/login"
	postData = url.Values{
		"allow_registration": {"false"},
		"app_name":           {"com.xiaomi.hm.health"},
		"app_version":        {"6.3.5"},
		"code":               {access},
		"country_code":       {"CN"},
		"device_id":          {"2C8B4939-0CCD-4E94-8CBA-CB8EA6E613A1"},
		"device_id_type":     {"uuid"},
		"device_model":       {"phone"},
		"dn":                 {"api-user.huami.com%2Capi-mifit.huami.com%2Capp-analytics.huami.com"},
		"grant_type":         {"access_token"},
		"lang":               {"zh_CN"},
		"os_version":         {"1.5.0"},
		"source":             {"com.xiaomi.hm.health"},
		"third_name":         {usernameType},
	}.Encode()
	resp, err := HttpPost(accountUrl, postData, headers)
	if err != nil {
		return account, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return account, err
	}
	err = json.Unmarshal(body, &account)
	if err != nil {
		return account, err
	}

	return account, err
}

func (a *Account) Set(step int) error {
	if step > 99999 || step < 1 {
		return errors.New("输入步数不符合要求，数量应在 1-99999 之间")
	}

	userID := a.TokenInfo.UserID
	appToken := a.TokenInfo.AppToken
	stepUrl := "https://api-mifit-cn.huami.com/v1/data/band_data.json"
	lastDeviceID := "DA932FFFFE8816E7"
	dataJson := `[{"summary":"{\"stp\":{\"runCal\":7,\"cal\":111,\"conAct\":0,\"stage\":[],\"ttl\":` + strconv.Itoa(step) + `,\"dis\":3102,\"rn\":2,\"wk\":43,\"runDist\":146,\"ncal\":0},\"v\":5,\"goal\":2000}","data":[{"stop":1439,"value":"","did":"【last_deviceid】","tz":32,"src":24,"start":0}],"data_hr":"","summary_hr":"{\"ct\":0,\"id\":[]}","date":"` + time.Now().Format("2006-01-02") + `"}]`
	postData := url.Values{
		"data_json":           {dataJson},
		"device_type":         {"0"},
		"enableMutiDevice":    {"1"},
		"last_deviceid":       {lastDeviceID},
		"last_source":         {"24"},
		"last_sync_data_time": {"1668682800"},
		"userid":              {userID},
		// "uuid":                {""},
	}.Encode()
	headers["apptoken"] = appToken
	resp, err := HttpPost(stepUrl, postData, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}

func run(ctx *gin.Context) {
	username, exist := ctx.GetQuery("username")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not found username"})
		return
	}
	password, exist := ctx.GetQuery("password")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not found password"})
		return
	}
	step, exist := ctx.GetQuery("step")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not found step"})
		return
	}
	ezap.Infof("准备上传步数 username: %s, password: %s, step: %s", username, password, step)
	a, err := getAccount(username, password)
	if err != nil {
		ezap.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = a.Set(10000)
	if err != nil {
		ezap.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	} else {
		ezap.Info("成功设置步数")
		ctx.String(http.StatusOK, username+" 成功设置步数 "+step)
	}
}
