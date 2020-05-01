package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

//获取s3连接
func GetCephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}
	//初始化ceph的信息
	auth := aws.Auth{
		AccessKey: "",
		SecretKey: "",
	}
	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "127.0.0.1:9080",
		S3Endpoint:           "127.0.0.1:9080",
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}

	//创建s3连接
	return s3.New(auth, curRegion)
}

//获取制定的bucket对象
func GetCephBucket(bucket string) * s3.Bucket{
	conn := GetCephConnection()
	return conn.Bucket(bucket)
}