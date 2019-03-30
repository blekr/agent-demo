package main

import (
	"encoding/json"
	"github.com/blekr/agent-demo/biz"
	"github.com/blekr/agent-demo/errors"
	. "github.com/blekr/agent-demo/util"
	"log"
	"net/http"
)

type response struct {
	ErrorCode string
	ErrorMessage string
	Data interface{}
}

func startHandler(w http.ResponseWriter, r *http.Request) (interface{}, error)  {
	var cmd biz.Command
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		return nil, err
	}
	return biz.StartProcess(&cmd)
}

func stopHandler(w http.ResponseWriter, r *http.Request) (interface{}, error)  {
	var req struct {
		Pid int
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return nil, biz.StopProcess(req.Pid)
}

func showHandler(w http.ResponseWriter, r *http.Request) (interface{}, error)  {
	var req struct {
		Pid int
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return biz.ShowProcess(req.Pid)
}

type CommonHandler func (http.ResponseWriter, *http.Request) (interface{}, error)

// Common processing for all handler such as logging and global error handling
// HTTP header "request-id" must be set for the sake of request tracking
func buildHandler(handler CommonHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ILog.Printf("begin: %v %v", r.Method, r.URL)

		res, err := handler(w, r)
		if err == nil {
			_ = json.NewEncoder(w).Encode(response{
				ErrorCode: "",
				ErrorMessage: "",
				Data: res,
			})
			ILog.Printf("end successfully: %v %v", r.Method, r.URL)
			return
		}

		ELog.Printf("end with error: %v %v: %v", r.Method, r.URL, err)
		switch v := err.(type) {
		case *errors.AppError:
			_ = json.NewEncoder(w).Encode(response{
				ErrorCode:    v.Code,
				ErrorMessage: v.Message,
				Data:         nil,
			})
		default:
			_ = json.NewEncoder(w).Encode(response{
				ErrorCode: "UNKNOWN",
				ErrorMessage: err.Error(),
				Data: nil,
			})
		}
	}

}

func main ()  {
	http.HandleFunc("/start", buildHandler(startHandler))
	http.HandleFunc("/stop", buildHandler(stopHandler))
	http.HandleFunc("/show", buildHandler(showHandler))
	ILog.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}