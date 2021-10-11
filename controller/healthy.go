package controller

//import (
//	log "github.com/cihub/seelog"
//	"github.com/julienschmidt/httprouter"
//	"net/http"
//)
//
//var healthJson = []byte{0x7B, 0x22, 0x6D, 0x73, 0x67, 0x22, 0x3A, 0x22, 0x50, 0x72, 0x6F, 0x67, 0x72, 0x61, 0x6D, 0x20, 0x69, 0x73, 0x20, 0x72, 0x75, 0x6E, 0x6E, 0x69, 0x6E, 0x67, 0x22, 0x7D}
//
//func HealthCheck(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	method, contentType := r.Method, r.Header.Get("Content-Type")
//	log.Tracef("HealthCheck-INFO: Method:%v Content-Type:%v IP:%v UserAgent:%v",
//		method,
//		contentType,
//		util.GetRemoteIp(r),
//		util.GetUserAgent(r))
//	w.WriteHeader(http.StatusOK)
//	w.Write(healthJson)
//	return
//}
