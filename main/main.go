package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Yaml struct {
	Auth_token string `yaml:"auth_token"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetCookie() string {
	conf := new(Yaml)
	yamlFile, err := ioutil.ReadFile("config.yaml")

	//这一步将struct类型和config.yaml配置文件管理关联
	err = yaml.Unmarshal(yamlFile, conf)

	checkError(err)
	auth_token := conf.Auth_token
	return auth_token
}

func setCookiefile(src string) {
	stu := &Yaml{
		Auth_token: "",
	}
	data, err := yaml.Marshal(stu)
	checkError(err)
	err = ioutil.WriteFile(src, data, 0777)
	checkError(err)
}

//结果去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func main() {

	var name string
	var page int
	var child bool
	var rate int
	var deep int
	var id []string
	var company []string
	var domainlist []string
	var applist []string
	var weixinlist []string
	auth_token := GetCookie()

	flag.StringVar(&name, "n", "", "查询目标名称")
	flag.IntVar(&page, "page", 0, "查询页面数")
	flag.BoolVar(&child, "child", false, "是否查询子公司,默认为false")
	flag.IntVar(&rate, "rate", 90, "控股比例,默认为90%控股")
	flag.IntVar(&deep, "deep", 1, "子公司查询递归深度,默认为1")
	flag.Parse()

	f, err := os.Open("config.yaml")
	if err != nil && os.IsNotExist(err) {
		setCookiefile("config.yaml")
	}
	f.Close()

	if page != 0 {
		id, company = GetPage(auth_token, name, page)
	} else {
		id, company = GetFirstCompany(auth_token, name)
	}
	for i := range id {
		domainlist = GetDomain(id[i], company[i], auth_token)
		applist = Getapp(id[i], company[i], auth_token)
		weixinlist = Getweixin(id[i], company[i], auth_token)
		if child == true {
			domainlisttmp, applisttmp, weixintmp := GetChildCompany_jt(id[i], auth_token, rate)
			domainlist = append(domainlist, domainlisttmp...)
			applist = append(applist, applisttmp...)
			weixinlist = append(weixinlist, weixintmp...)
			domainlisttmp, applisttmp, weixintmp = GetChildCompany_dc(id[i], auth_token, rate)
			domainlist = append(domainlist, domainlisttmp...)
			applist = append(applist, applisttmp...)
			weixinlist = append(weixinlist, weixintmp...)
			domainlisttmp, applisttmp, weixintmp = GetChildCompany_yh(id[i], auth_token, rate)
			domainlist = append(domainlist, domainlisttmp...)
			applist = append(applist, applisttmp...)
			weixinlist = append(weixinlist, weixintmp...)
			domainlisttmp, applisttmp, weixintmp = GetChildCompany_gq(id[i], auth_token, rate, deep)
			domainlist = append(domainlist, domainlisttmp...)
			applist = append(applist, applisttmp...)
			weixinlist = append(weixinlist, weixintmp...)
		}
	}

	//去重
	domainlist = RemoveRepeatedElement(domainlist)
	applist = RemoveRepeatedElement(applist)
	weixinlist = RemoveRepeatedElement(weixinlist)

	//输出汇总结果
	fmt.Println()
	fmt.Println("[+]域名资产汇总")
	fmt.Println("==================================")

	for i := range domainlist {
		fmt.Println(domainlist[i])
	}
	fmt.Println()
	fmt.Println("[+]app资产汇总")
	fmt.Println("==================================")
	for i := range applist {
		fmt.Println(applist[i])
	}
	fmt.Println()
	fmt.Println("[+]微信公众号资产汇总")
	fmt.Println("==================================")
	for i := range weixinlist {
		fmt.Println(weixinlist[i])
	}

}
