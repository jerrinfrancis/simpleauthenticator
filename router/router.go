package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Router struct {
	mux      map[string]map[string]http.Handler
	NotFound http.Handler
}

func corsMiddlewareHandler(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Reached CORS")
		/*	if (r.Method == http.MethodPost || r.Method == http.MethodPatch) && os.Getenv("BLOG_IN_PROD") == "X" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Post method not supported"))
			return
		}*/
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(w, r)

	})
}

func NewRouter() *Router {
	var mux map[string]map[string]http.Handler
	mux = make(map[string]map[string]http.Handler)
	mux["get"] = make(map[string]http.Handler)
	mux["post"] = make(map[string]http.Handler)
	mux["delete"] = make(map[string]http.Handler)
	mux["options"] = make(map[string]http.Handler)
	mux["patch"] = make(map[string]http.Handler)
	return &Router{
		mux:      mux,
		NotFound: corsMiddlewareHandler(http.NotFoundHandler()),
	}
}

func (r *Router) SetHandlerFunc(method, path string, fn http.HandlerFunc) {
	r.mux[strings.ToLower(method)][path] = corsMiddlewareHandler(http.Handler(fn))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	method := strings.ToLower(req.Method)
	urlstring := req.URL.String()
	fmt.Println(urlstring)
	var basePath string
	ctx := req.Context()
	if strings.Contains(urlstring, "?") {
		basePath = strings.Split(urlstring, "?")[0]
		filterstring := strings.Split(urlstring, "?")[1]
		m, err := url.ParseQuery(filterstring)
		fmt.Println("m:", m)
		if err != nil {
			log.Println(err.Error())
			return
		}
		js, err := json.Marshal(m)
		if err != nil {
			log.Println(err.Error())
			return
		}
		filters := strings.ReplaceAll(string(js), ",", " ,")
		ctx = context.WithValue(ctx, "filter", filters)

	} else {
		basePath = urlstring
	}
	fmt.Println("basepath" + basePath + method)
	if hm, ok := r.mux[method]; ok {
		fmt.Println("reached method")
		//	/uploadfile/abcd.png
		if strings.Count(basePath, "/") == 2 {
			params := strings.Split(basePath, "/")
			param1 := params[2]
			basePath = "/" + params[1]
			ctx = context.WithValue(ctx, "param1", param1)
		}
		if h, ok := hm[basePath]; ok {
			h.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	r.NotFound.ServeHTTP(w, req)

}
