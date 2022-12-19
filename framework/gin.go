package framework

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

type InputTransactionMiddleware struct {
	DbConn          *sql.DB
	UseJwt          bool
	UseNotification bool
	ChannelID       string
}

// TransactionMiddleware : new!! to setup the database transaction middleware
func TransactionMiddleware(inp *InputTransactionMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		txHandle, err := inp.DbConn.BeginTx(ctx, nil)
		if err != nil {
			panic(err)
		}

		fmt.Println(">>> START>>> ----------------------------------------------------------------------")
		fmt.Println("Beginning database transaction")
		fmt.Printf("Requested URL : %s\n", c.Request.URL.Path)

		defer func() {
			// recovering from panic in app
			if r := recover(); r != nil {
				c.String(500, fmt.Sprint(r))
				txHandle.Rollback()

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

		// get jwt claim
		var appReq AppRequest

		appReq.Tx = txHandle

		appReq.PathParams = GetPathParams(c)
		appReq.QueryParams = GetQueryParams(c)
		appReq.ReqBody = GetRequestBodyPtr(c)
		appReq.UseJwt = inp.UseJwt

		c.Set("app_request", appReq)

		c.Next()

		if StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			fmt.Println("Committing transactions")
			fmt.Println("---------------------------------------------------------------------- <<<  END  <<<")
			if err := txHandle.Commit(); err != nil {
				fmt.Println("trx commit error: ", err)
			}
		} else {
			fmt.Println("Rolling back transaction due to status code   : " + fmt.Sprintf("%d", c.Writer.Status()))
			fmt.Println("---------------------------------------------------------------------- <<<  END  <<<")
			txHandle.Rollback()
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

// func GetRouter() *gin.Engine {
// 	gin.DefaultWriter = os.Stdout

// 	router := gin.New()

// 	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 		// your custom format
// 		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \" %s\" %s\"\n",
// 			param.ClientIP,
// 			param.TimeStamp.Format(time.RFC1123),
// 			param.Method,
// 			param.Path,
// 			param.Request.Proto,
// 			param.StatusCode,
// 			param.Latency,
// 			param.Request.UserAgent(),
// 			param.ErrorMessage,
// 		)
// 	}))

// 	return router
// }
