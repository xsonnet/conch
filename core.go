package conch

import (
	"net/http"
	"regexp"
	"time"
)

type Router struct {
	Pattern 	string
	Func 		func(ctx Context)
}

type App struct {
	Routers 	[]Router
	LogPath		string
}

// 中间件
func (app App)middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(res, req)
		app.Log("%s (%v)", req.URL.Path, time.Since(start))
	})
}

// 路由处理
func (app App) handleRouter(res http.ResponseWriter, req *http.Request) {
	isFound := false
	for _, router := range app.Routers {
		// 循环匹配，先添加的先匹配
		reg, err := regexp.Compile(router.Pattern)
		if err != nil {
			continue
		}
		if reg.MatchString(req.URL.Path) {
			isFound = true
			router.Func(Context{
				Response: res,
				Request: req,
				App: app,
			})
		}
	}
	if !isFound {
		// 未匹配到路由
		_, _ = res.Write([]byte("Url Not Found!"))
	}
}

// 处理静态文件
func (app App) Static(pattern string, path string) {
	http.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(path))))
}

// Handle
func (app App) Handle(pattern string, handle http.Handler) {
	http.Handle(pattern, handle)
}

// 启动监听
func (app App) Run(addr string) {
	http.Handle("/", app.middleware(http.HandlerFunc(app.handleRouter)))
	app.Log("Server started: " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		app.Log("Http listened failed: %s", err.Error())
	}
}

// 打印日志
func (app App) Log(format string, arg ...interface{}) {
	Log{Path: app.LogPath}.Out(format, arg...)
}