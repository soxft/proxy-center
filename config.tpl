Type: 1 # 1 for duomi & 2 for douya 3 for xiequ
Frequency: 5 # 0 秒 则为系统自动判断
Protocol: "http" # http or https or socks5
Api:
  DmApi: "http://api.3ip.cn/dmgetip.asp?apikey=3a019c15&pwd=&getnum=50&httptype=3&geshi=2&fenge=1&fengefu=&Contenttype=1&operate=all"
  DyApi: "https://api.douyadaili.com/proxy/?service=GetUnl&authkey=&num=10&lifetime=1&prot=0&format=json&high=0&detail=1"
  XqApi: "http://api.xiequ.cn/VAD/GetIp.aspx?act=get&uid=&vkey=&num=5&time=30&plat=0&re=0&type=0&so=1&ow=1&spl=1&addr=&db=1"
Server:
  Address: 127.0.0.1:9876