package framework

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type InputTransactionMiddleware struct {
	// UseAuthHeader bool
}

func TransactionMiddleware(inp *InputTransactionMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		var appReq AppRequest

		// check url if there's prefix /admin and
		if strings.Split(c.Request.URL.Path, "/")[1] == "admin" {
			if len(c.Request.Header["User_id"]) == 0 {
				c.AbortWithStatus(401)
			}

			userId, err := strconv.Atoi(c.Request.Header["User_id"][0])
			if err != nil {
				c.AbortWithStatus(401)
			}

			user := GetUser(userId)
			if user == nil {
				c.AbortWithStatus(401)
			}

			// check authorization based on rbac permission matrix
			rbacMatrix := GetRBACMatrix()
			fullPath := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
			IsAuth, exist := (*rbacMatrix)[RBACItem{FullPath: fullPath, Role: user.Role}]
			if !exist {
				fmt.Println("not exist")
				c.AbortWithStatus(401)
			}

			if !IsAuth {
				fmt.Println("not authorized")
				c.AbortWithStatus(401)
			}

			appReq.User = user
		}

		fmt.Println(">>> START>>> ----------------------------------------------------------------------")
		fmt.Println("Beginning transaction")
		fmt.Printf("Requested URL : %s\n", c.Request.URL.Path)

		defer func() {
			// recovering from panic in app
			if r := recover(); r != nil {
				c.String(500, fmt.Sprint(r))
				// txHandle.Rollback()

				// get a stack of app
				buf := make([]byte, 1<<10)
				runtime.Stack(buf, false)

				fmt.Println("Error Message : " + fmt.Sprint(r))
				fmt.Printf("%s\n", buf)

				//send slack message
				fmt.Println("Rolling back transaction caused by server error")
				fmt.Println("---------------------------------------------------------------------- <<<  END")
			}
		}()

		appReq.PathParams = GetPathParams(c)
		appReq.QueryParams = GetQueryParams(c)
		appReq.ReqBody = GetRequestBodyPtr(c)

		c.Set("app_request", appReq)

		c.Next()

		if StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			fmt.Println("Committing transactions")
			fmt.Println("---------------------------------------------------------------------- <<<  END  <<<")
			// if err := txHandle.Commit(); err != nil {
			// 	fmt.Println("trx commit error: ", err)
			// }
		} else {
			fmt.Println("Rolling back transaction due to status code   : " + fmt.Sprintf("%d", c.Writer.Status()))
			fmt.Println("---------------------------------------------------------------------- <<<  END  <<<")
			// txHandle.Rollback()
		}
	}
}

//StatusInList -> checks if the given status is in the list
func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

func GetTx(c *gin.Context) *sql.Tx {
	tx, exist := c.Get("db_trx")
	if !exist {
		panic("error on getting tx")
	}

	return tx.(*sql.Tx)
}

func GetAppRequest(c *gin.Context) *AppRequest {
	appRequest, exist := c.Get("app_request")
	if !exist {
		panic("error on getting app request")
	}

	ret := appRequest.(AppRequest)

	return &ret
}

func GetRequestBody(c *gin.Context) string {
	body, err := c.GetRawData()
	if err != nil {
		panic("error on getting Request Body")
	}

	// logging input
	fmt.Printf("Request Body : \n%s\n", string(body))

	return string(body)
}

func GetRequestBodyPtr(c *gin.Context) *string {
	body, err := c.GetRawData()

	fmt.Println("GetRequestBodyPtr : ", string(body))

	if err != nil {
		panic("error on getting Request Body")
	}

	// logging input
	fmt.Printf("Request Body : \n%s\n", string(body))

	ret := string(body)

	ret = strings.ReplaceAll(ret, "'", "''")

	return &ret
}

func GetPathParams(c *gin.Context) *map[string]string {
	ret := make(map[string]string)

	for _, v := range c.Params {
		ret[v.Key] = v.Value
	}

	// logging input
	fmt.Printf("Path Parameters : %+v\n", ret)

	return &ret
}

func GetQueryParams(c *gin.Context) *map[string]string {
	ret := make(map[string]string)

	for k, v := range c.Request.URL.Query() {
		ret[k] = v[0]
	}

	// logging input
	fmt.Printf("Query Parameters : %+v\n", ret)

	return &ret
}
