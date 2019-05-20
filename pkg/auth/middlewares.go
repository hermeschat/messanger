package auth

//
//// Allow to check authentication of routes
//func Allow(h echo.HandlerFunc, roles ...string) echo.MiddlewareFunc {
//	return func(h echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAaaa")
//			identity, err := Authenticate(&c, "", roles...)
//			if err != nil {
//				// Auth failed
//				if _, ok := err.(UnauthorizedError); ok {
//					return GetUnAuthorizedError("")
//				}
//				if _, ok := err.(ForbiddenError); ok {
//					return GetForbiddenError("")
//				}
//				fmt.Println("error on validate jwt", err)
//				return GetUnAuthorizedError("")
//			}
//			c.Set("identity", identity)
//			return h(c)
//		}
//	}
//}
