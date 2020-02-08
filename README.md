「 Lighthouse and Whale 」

# Quark
Quark is a core web development library for the Go programming language

<br>

## Quick start

```go
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
```
