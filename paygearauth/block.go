package paygearauth

//
//// helpers function for checking card is blocked or not
//func CheckIsBlock(h echo.HandlerFunc) echo.MiddlewareFunc {
//	return func(h echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			sourceUserid := c.Param("source")
//			redisCon, err := ConnectRedis()
//			if err != nil {
//				return utils.GetValidationError("تعداد رمز اشتباه بیش از حد مجاز می باشد")
//			}
//			key := "Block:"+sourceUserid
//			_,err = redisCon.Get(key).Result()
//			if err == redis.Nil {
//				return h(c)
//			}else if err != nil {
//				return err
//			}else {
//				return utils.GetValidationError("تعداد رمز اشتباه بیش از حد مجاز می باشد")
//			}
//			return h(c)
//		}
//	}
//}
