package main

import (
	"os"
	"fmt"
	"strconv"
	"time"
	"regexp"
	"net/http"
	"io/ioutil"
)

const (
	SohuVideoIdRegStr = "vid\\s*=\\s*\"(\\d+)\""
	SohuVideoUrlRegStr = "http://tv.sohu.com/(\\d+)/n(\\d+).shtml"
)

func file_get_contents(url string) string {
	r, e := http.Get(url)
	if e != nil {
		return ""
	}
	defer r.Body.Close()
	c, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return ""
	}
	return string(c)
}
// http://hot.vrs.sohu.com/ipad1349813.m3m8
// http://hot.vrs.sohu.com/ipad1349813.m3m8
// http://hot.vrs.sohu.com/ipad1717852_4508031888580_4928173.m3u8?plat=h5&uid=1401271517200744&ver=1&prod=h5&pt=2&pg=1&ch=2&isad=1
// 1397712726852
// 4508031888580
// http://js.tv.itc.cn/site/pad/video14041101.js

func get_vid(url string) string {
	html := file_get_contents(url)
	r:= regexp.MustCompile(SohuVideoIdRegStr)
	rs:=r.FindStringSubmatch(html)
	return rs[1]
}

/*
(function() {
    ""._shift_en || (String.prototype._shift_en = function(e) {
        var t = e.length,
        n = 0;
        return this.replace(/[0-9a-zA-Z]/g, 
        function(r) {
            var i = r.charCodeAt(0),
            s = 65,
            o = 26;
            i >= 97 ? s = 97: i < 65 && (s = 48, o = 10);
            var u = i - s;
            return String.fromCharCode((u + e[n++%t]) % o + s)
        })
    })
})()
 */
// 视频编码校验， 上面为js版本， 如果失败， 则需要同步修改规则
func sohu_shift_en(num_str string, code [4]int) string {
	t := len(code)
    n := 0
    str := []byte(num_str)
    strlen := len(str)
    var newStr []byte
    s, o := 65, 26
    u := 0
    for i := 0; i < strlen; i++ {
    	s = 65
    	o = 26
    	if str[i] >= 97 {
    		s = 97
    	} else if str[i] < 65 {
    		s = 48
    		o = 10
    	}
    	u,_ = strconv.Atoi(string(str[i]))
		m := ((u + code[n % t]) % o) + s
		n ++
    	newStr = append(newStr, byte(m))
	}
	newA := ""
	for i := 0; i < len(newStr); i++ {
		newA = newA + string(newStr[i])
	}
	return newA
}

func sohu_timestamp() string {
	tmp := time.Now().UnixNano()
	ss := strconv.FormatInt(tmp, 10)
	return ss[0:13]
}

func main() {

	argsWithProg := os.Args
	url := argsWithProg[1]
	if url[0:19] != "http://tv.sohu.com/" {
		fmt.Println("参数无效1")
	}

	// SohuVideoUrlRegStr
	r:= regexp.MustCompile(SohuVideoUrlRegStr)
	rs:=r.FindStringSubmatch(url)
	if rs == nil {
		fmt.Println("暂时只支持http://tv.sohu.com/下的视频")
		return
	}
	
	t := sohu_timestamp()
	code := [4]int{23,12,131,1321}
	m3u8 := "http://hot.vrs.sohu.com/ipad"
	vid := get_vid(url)
	m3u8 = m3u8 + vid + "_" + sohu_shift_en(t, code) + "_" + sohu_shift_en(vid, code) + ".m3u8"

	fmt.Println(m3u8)
	
}