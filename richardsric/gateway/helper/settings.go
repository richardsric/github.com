package helper

import (
	"fmt"
	"time"
)

type ApiSettings struct {
	BaseUrl string
	ApiVersion string
	ReqTimeOut time.Duration
}

// Settings this is use to to load database settings to the struct
func Settings()(baseUrl, apiVersion string, reqTimeOut time.Duration)  {
	fmt.Println("entered settings function")
	con, err := OpenConnection()
    if err != nil {
      fmt.Println(err)
      return
    }
	defer con.Close()
	
	sRow := con.Db.QueryRow("SELECT base_url,api_version,req_time_out FROM exapi_settings WHERE exchange_id = $1", 1)
	var API_BASE,API_VERSION string
	var DEFAULT_HTTPCLIENT_TIMEOUT time.Duration
	//var DEFAULT_HTTPCLIENT_TIMEOUT time.Duration

		err = sRow.Scan(&API_BASE,&API_VERSION,&DEFAULT_HTTPCLIENT_TIMEOUT)
		if err != nil {
			fmt.Println("ERROR!!...due to",err)
			return
		}
		return API_BASE, API_VERSION, DEFAULT_HTTPCLIENT_TIMEOUT 
}