package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"gopkg.in/eapache/queue.v1"
	"log"
	"time"
)

var Que *queue.Queue

// Get proxy Ip Address 获取代理IP
func Get(ctx context.Context) (Proxy, error) {
	// 要求获取到的代理IP剩余时间大于 10 秒
	var endTimeNeed int64 = 30

	var _pData Proxy
	// 一直remove 直到找到一个有效的
	for {
		select {
		case <-ctx.Done():
			return Proxy{}, errors.New("context timeout")
		default:
			if Que.Length() == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			_pData = Que.Remove().(Proxy)

			if _pData.EndTime-time.Now().Unix() >= endTimeNeed {
				log.Printf("[INFO] proxy get %s (%s | %d) length: %d", _pData.Addr, _pData.City, _pData.EndTime-time.Now().Unix(), Que.Length())
				return _pData, nil
			}
		}
	}
}

func ClearProxy() {
	for Que.Length() > 0 {
		_pData := Que.Peek().(Proxy)
		// 如果过期时间小于当前时间 则移除
		if _pData.EndTime >= time.Now().Unix()+10 {
			break
		} else {
			//log.Printf("[NOTICE] proxy remove %s (%s | %s)", _pData.Addr, _pData.City, time.Unix(_pData.EndTime, 0).Format("2006-01-02 15:04:05"))
			Que.Remove()
		}
	}
}

func MainCron() {

	defer func() {
		if err := recover(); err != nil {
			MainCron()
		}
	}()

	for {
		log.Printf("[INFO] start to get proxys")

		// 尝试取出 end_time 最小值
		var endTimeMin int64 = 600

		if proxys, err := getProxy(); err != nil {
			log.Printf("[Error] Cant get Proxy: %s", err.Error())
			log.Printf("[Info] 2 秒后重新获取 Proxys")
			time.Sleep(time.Second * 2)

			continue
		} else {
			for _, _pData := range proxys {
				ttl := _pData.EndTime - time.Now().Unix()

				log.Printf("[INFO] proxy add %s (%s | %s | %d)", _pData.Addr, _pData.City, time.Unix(_pData.EndTime, 0).Format("2006-01-02 15:04:05"), ttl)

				Que.Add(_pData)

				// 找出最小值
				if ttl < endTimeMin {
					endTimeMin = ttl
				}
			}
		}

		// endTimeMin 是 这些代理中最小的剩余过期时间

		// 如果 endTimeMin >= 60s 则 duration = endTimeMin - 30s 这是自动判断 保证在过期前自动重新获取
		duration := endTimeMin - 5
		if endTimeMin >= 60 {
			duration = endTimeMin - 30
		}

		// 如果用户指定
		if viper.GetInt("Frequency") > 0 {
			duration = viper.GetInt64("Frequency")
		}

		log.Printf("[INFO] 将在 %d 秒后 再次获取 proxy 列表", duration)
		time.Sleep(time.Second * time.Duration(duration))
	}
}

func getProxy() ([]Proxy, error) {
	switch viper.GetInt("type") {
	case 1:
		return getProxyDm()
	case 2:
		return getProxyDy()
	default:
		return getProxyXq()
	}
}

func getProxyDy() ([]Proxy, error) {
	client := resty.New().SetTimeout(time.Second * 2)

	var result DyResp
	if _, err := client.
		R().SetResult(&result).
		ForceContentType("application/json").
		Get(viper.GetString("Api.DyApi")); err != nil {

		log.Printf("[Warning] Get Proxy Api Request Failed, Retry in 500 milliseconds | Err: %s", err.Error())

		time.Sleep(time.Millisecond * 500)

		return getProxy()
	}
	if result.Ret != 200 {
		return []Proxy{}, errors.New(fmt.Sprintf("Code: %d, Msg: %s", result.Ret, result.Msg))
	}

	var _proxy []Proxy

	for _, pD := range result.Data {
		r, err := timeParse(pD.Deadline)
		if err != nil {
			log.Printf("[Warning] error time format %s", pD.Deadline)
			continue
		}

		_proxy = append(_proxy, Proxy{
			Addr:    fmt.Sprintf("%s:%s", pD.Ip, pD.Port),
			EndTime: r.Unix() - 2, // 安全起见 减少 10 秒
			City:    pD.City,
		})
	}

	return _proxy, nil
}

// 多米
func getProxyDm() ([]Proxy, error) {
	client := resty.New().SetTimeout(time.Second * 2)

	var result DmResp
	if _, err := client.
		R().SetResult(&result).
		ForceContentType("application/json").
		Get(viper.GetString("Api.DmApi")); err != nil {

		log.Printf("[Warning] Get Proxy Api Request Failed, Retry in 500 milliseconds | Err: %s", err.Error())

		time.Sleep(time.Millisecond * 500)

		return getProxy()
	}
	if result.Code != 0 {
		return []Proxy{}, errors.New(fmt.Sprintf("Code: %d, Msg: %s", result.Code, result.Msg))
	}

	var _proxy []Proxy

	for _, pD := range result.Data {
		r, err := timeParse(pD.Endtime)
		if err != nil {
			log.Printf("[Warning] error time format %s", pD.Endtime)
			continue
		}

		_proxy = append(_proxy, Proxy{
			Addr:    fmt.Sprintf("%s:%d", pD.Ip, pD.Port),
			EndTime: r.Unix() - 2, // 安全起见 减少 10 秒
			City:    pD.City,
		})
	}

	return _proxy, nil
}

func getProxyXq() ([]Proxy, error) {
	client := resty.New().SetTimeout(time.Second * 2)

	var result XqResp
	if _, err := client.
		R().SetResult(&result).
		ForceContentType("application/json").
		Get(viper.GetString("Api.XqApi")); err != nil {

		log.Printf("[Warning] Get Proxy Api Request Failed, Retry in 500 milliseconds | Err: %s", err.Error())

		time.Sleep(time.Millisecond * 500)

		return getProxy()
	}
	if result.Code != 0 {
		return []Proxy{}, errors.New(fmt.Sprintf("Code: %d, Msg: %s", result.Code, result.Msg))
	}

	var _proxy []Proxy

	for _, pD := range result.Data {
		_proxy = append(_proxy, Proxy{
			Addr:    fmt.Sprintf("%s:%d", pD.Ip, pD.Port),
			EndTime: time.Now().Unix() + 28, // 安全起见 减少 10 秒
			City:    pD.IpAddress,
		})
	}

	return _proxy, nil
}

func timeParse(t string) (time.Time, error) {
	var timeFmt = []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"2006-1-02 15:04:05",
		"2006/1/02 15:04:05",
		"2006-01-2 15:04:05",
		"2006/01/2 15:04:05",
		"2006-1-2 15:04:05",
		"2006/1/2 15:04:05",
	}

	var r time.Time
	var err error
	for _, f := range timeFmt {
		r, err = time.ParseInLocation(f, t, time.Local)
		if err == nil {
			break
		}
	}

	return r, err
}
