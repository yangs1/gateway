package provider

//
//// 平滑轮询算法
//func (manager *ProxyManager) check() {
//
//	client := http.Client{}
//	for _, s := range manager.Servers {
//		res, err := client.Head(s.Host)
//		if res != nil {
//			defer res.Body.Close()
//		}
//
//		if err != nil {
//			// 服务器宕机
//			continue
//		}
//
//		if res.StatusCode != 200 {
//			// 请求错误
//			continue
//		}
//	}
//}
//
//func (manager *ProxyManager) CheckerSevers() {
//	t := time.NewTimer(time.Second * 5)
//
//	for {
//		select {
//		case <-t.C:
//			manager.check()
//		}
//	}
//}
