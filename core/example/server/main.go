package main

import "github.com/lighthouse-and-whale/quark/core"

func main() {
	HTTP := core.NewHTTP(":2233")
	HTTP.Add("", "", func(ct core.Context) error {
		return ct.ToJson(200, "success", map[string]interface{}{
			"response": "data",
		})
	})
	core.ShowAPIs(true)
	if e := HTTP.Server.ListenAndServe(); e != nil {
		println(e.Error())
	}
}

// 连接数据库
//var DB = store.NewStorageDatabase(false,
//	fmt.Sprintf("%s", core.CONFIG["db_nodes"]),
//	fmt.Sprintf("%s", core.CONFIG["db_name"]),
//	fmt.Sprintf("%s", core.CONFIG["db_user"]),
//	fmt.Sprintf("%s", core.CONFIG["db_pwd"]),
//)
//
//func Verify(ct core.Context) error {
//	if ct.User == "" {
//		return ct.To401()
//	}
//	return core.Next
//}
//
//func main() {
//	defer core.PProfCPU().Stop()
//
//	HTTP := core.NewHTTP(":2233")
//	DB.Name()
//	v1 := HTTP.Add("", "/v1.app", core.UserRequestInterceptor)
//	v1.Add("online", "/online", core.UserOnlineStateMachine)
//	// api.Group(v1)
//	core.ShowAPIs(true)
//	if e := HTTP.Server.ListenAndServe(); e != nil {
//		println(e.Error())
//	}
//}

// go tool pprof -http=127.0.0.1:3000 cpu.pprof
