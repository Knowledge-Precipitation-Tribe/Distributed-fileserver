package oss

import(
	cfg "Distributed-fileserver/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

//创建oss Client对象
func Client() *oss.Client{
	if ossCli != nil{
		return ossCli
	}
	ossCli, err := oss.New(cfg.OSSEndpoint,
		cfg.OSSAccesskeyID, cfg.OSSAccessKeySecret)
	if err != nil{
		fmt.Println(err.Error())
		return nil
	}
	return ossCli
}


//获取bucket存储空间
func Bucket() *oss.Bucket{
	cli := Client()
	if cli != nil{
		bucket, err := cli.Bucket(cfg.OSSBucket)
		if err != nil{
			fmt.Println(err.Error())
			return nil
		}
		return bucket
	}
	return nil
}

//DownloadURL: 临时授权下载url
func DownloadURL(objName string) string{
	signURL, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil{
		fmt.Println(err.Error())
		return ""
	}
	return signURL
}

//设置生命周期规则
func BuildLifecycleRule(bucketName string){
	//制定bucket中以test开头的对象，30天内没有修改则自动删除
	ruleTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rules := []oss.LifecycleRule{ruleTest1}
	Client().SetBucketLifecycle(bucketName, rules)
}