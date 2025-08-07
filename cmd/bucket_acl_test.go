package cmd

import (
	"context"
	"coscli/util"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tencentyun/cos-go-sdk-v5"
	"reflect"
	"testing"
)

func TestBucketAclCmd(t *testing.T) {
	fmt.Println("TestBucketAclCmd")
	testBucket = randStr(8)
	testAlias = testBucket + "-alias"
	setUp(testBucket, testAlias, testEndpoint, false, false)
	defer tearDown(testBucket, testAlias, testEndpoint, false)
	clearCmd()
	cmd := rootCmd
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	genDir(testDir, 3)
	defer delDir(testDir)

	Convey("test coscli bucket_acl", t, func() {
		Convey("success", func() {
			Convey("put", func() {
				clearCmd()
				// 创建模拟的 BucketService
				mockBucketService := &cos.BucketService{}
				// 创建打桩
				patches := gomonkey.NewPatches()
				defer patches.Reset()

				// 打桩 BucketService.PutACL 方法
				patches.ApplyMethod(
					reflect.TypeOf(mockBucketService), // 目标类型
					"PutACL",                          // 方法名
					func(_ *cos.BucketService, ctx context.Context, opt *cos.BucketPutACLOptions) (*cos.Response, error) {
						return nil, nil
					},
				)

				cmd := rootCmd
				args := []string{"bucket-acl", "--method", "put",
					fmt.Sprintf("cos://%s", testAlias), "--grant-read", "id=\"100000000003\",id=\"100000000002\""}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("get", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"bucket-acl", "--method", "get",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
		})
		Convey("fail", func() {
			Convey("clinet err", func() {
				clearCmd()
				cmd := rootCmd
				patches := ApplyFunc(util.NewClient, func(config *util.Config, param *util.Param, bucketName string) (client *cos.Client, err error) {
					return nil, fmt.Errorf("test put client error")
				})
				defer patches.Reset()
				args := []string{"bucket-acl", "--method", "put",
					fmt.Sprintf("cos://%s", testAlias), "--grant-read", "id=\"100000000003\",id=\"100000000002\""}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
		})
	})
}
