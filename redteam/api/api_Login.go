package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

func Login(c *gin.Context) {
	var num int //http 상태정보를 반환받을 변수

	user := model.User{}
	err := c.BindJSON(&user)
	if err != nil {
		model.SugarLogger.Error(err.Error())
		// 입력값이 제대로 바인딩 되지 않은경우 400 에러를 반환한다.
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk": false,
		})
		return
	}

	// 로그인 자격증명을 검사한다.
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	// 로그인 시도 횟수 제한 5회
	err, num, loginCount := user.IsAuthenticated(&conn) // 비밀번호 확인
	if err != nil {
		c.JSON(num, gin.H{
			"isOk":       false,
			"loginCount": loginCount,
		})
		return
	}

	accessToken, refreshToken, err := user.GetAuthToken()
	if err == nil { //여기서 토큰을 쿠키에 붙인다. // 각 1시간, 1주일
		c.SetCookie("access-token", accessToken, 604800, "", "", false, true)
		c.SetCookie("refresh-token", refreshToken, 604800, "", "", false, true)
		// https 사용시 refresh-token 의 secure -> true 로 변경한다.
		// (maxAge) 1800 -> 30분

		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
		return
	} else {
		// access 토큰이 발급되지 않은 경우 500에러를 반환한다.
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk": false,
		})
		log.Print("Login error occurred, account : ", Account)
		return
	}

}

func DelUser(c *gin.Context) {
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)
	user := model.User{}
	user.UserNo = num

	err := user.DelUser(&conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isOk": true,
	})
	return
}
