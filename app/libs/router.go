/**
  create by yy on 2019-07-03
*/

package libs

import (
	"strings"
)

var (
	HTTPMETHOD = map[string]bool{
		"GET":  true,
		"POST": true,
		// "PUT":       true,
		// "DELETE":    true,
		// "PATCH":     true,
		// "OPTIONS":   true,
		// "HEAD":      true,
		// "TRACE":     true,
		// "CONNECT":   true,
		// "MKCOL":     true,
		// "COPY":      true,
		// "MOVE":      true,
		// "PROPFIND":  true,
		// "PROPPATCH": true,
		// "LOCK":      true,
		// "UNLOCK":    true,
	}
	exceptMethod = []string{"Mapping", "Init"}
)

// func AutoRoute(ginRouter *gin.Engine, controller interface{}) {
//	/**
//	自动路由 (Auto Router)
//	*/
//	prefix := "/"
//	reflectValue := reflect.ValueOf(controller)
//	// Get controller name
//	ct := reflect.Indirect(reflectValue).Type()
//	controllerOriginName := ct.Name()
//	controllerName := strings.TrimSuffix(controllerOriginName, "Controller")
//	controllerName = strings.ToLower(controllerName)
//	ginRouter.AddToAutoRouterGroup(controllerName, "Controller", filterControllerName(reflectValue.Type().String()))
//	ginRouter.AddToAutoRouterController(controllerName, &ct)
//	// Get methods name
//	reflectType := reflectValue.Type()
//	length := reflectType.NumMethod()
//	for i := 0; i < length; i++ {
//		methodsName := reflectType.Method(i).Name
//		if !InSlice(methodsName, exceptMethod) {
//			ginRouter.AddToAutoRouterGroup(controllerName, strings.ToLower(methodsName), methodsName)
//			pattern := Join(prefix, controllerName, strings.ToLower(methodsName))
//			for m := range HTTPMETHOD {
//				ginRouter.Handle(m, pattern, AutoRouteExecute)
//			}
//		}
//	}
// }

// func AutoRouteExecute(ctx *gin.Context) {
//	httpPath := ctx.Request.URL.Path
//	httpPathArr := strings.Split(httpPath, "/")
//	// get controlle name
//	controllerName := httpPathArr[1]
//	// get method name
//	methodName := httpPathArr[2]
//	// 拼接控制器
//	routerGroup := ctx.GetAutoRouterGroup()
//	routerController := ctx.GetAutoRouterController()
//	controller := reflect.New((*routerController)[controllerName])
//	in := make([]reflect.Value, 1)
//	in[0] = reflect.ValueOf(ctx)
//	controller.MethodByName(routerGroup[controllerName][methodName]).Call(in)
// }

func filterControllerName(controllerName string) string {
	name := strings.Split(controllerName, "*")
	return name[1]
}

func Join(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return Clean(strings.Join(elem[i:], "/"))
		}
	}
	return ""
}

type lazybuf struct {
	s   string
	buf []byte
	w   int
}

func (b *lazybuf) index(i int) byte {
	if b.buf != nil {
		return b.buf[i]
	}
	return b.s[i]
}

func (b *lazybuf) append(c byte) {
	if b.buf == nil {
		if b.w < len(b.s) && b.s[b.w] == c {
			b.w++
			return
		}
		b.buf = make([]byte, len(b.s))
		copy(b.buf, b.s[:b.w])
	}
	b.buf[b.w] = c
	b.w++
}

func (b *lazybuf) string() string {
	if b.buf == nil {
		return b.s[:b.w]
	}
	return string(b.buf[:b.w])
}

func Clean(path string) string {
	if path == "" {
		return "."
	}

	rooted := path[0] == '/'
	n := len(path)

	// Invariants:
	//	reading from path; r is index of next byte to process.
	//	writing to buf; w is index of next byte to write.
	//	dotdot is index in buf where .. must stop, either because
	//		it is the leading slash or it is a leading ../../.. prefix.
	out := lazybuf{s: path}
	r, dotdot := 0, 0
	if rooted {
		out.append('/')
		r, dotdot = 1, 1
	}

	for r < n {
		switch {
		case path[r] == '/':
			// empty path element
			r++
		case path[r] == '.' && (r+1 == n || path[r+1] == '/'):
			// . element
			r++
		case path[r] == '.' && path[r+1] == '.' && (r+2 == n || path[r+2] == '/'):
			// .. element: remove to last /
			r += 2
			switch {
			case out.w > dotdot:
				// can backtrack
				out.w--
				for out.w > dotdot && out.index(out.w) != '/' {
					out.w--
				}
			case !rooted:
				// cannot backtrack, but not rooted, so append .. element.
				if out.w > 0 {
					out.append('/')
				}
				out.append('.')
				out.append('.')
				dotdot = out.w
			}
		default:
			// real path element.
			// add slash if needed
			if rooted && out.w != 1 || !rooted && out.w != 0 {
				out.append('/')
			}
			// copy element
			for ; r < n && path[r] != '/'; r++ {
				out.append(path[r])
			}
		}
	}

	// Turn empty string into "."
	if out.w == 0 {
		return "."
	}

	return out.string()
}
