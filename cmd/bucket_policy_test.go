package cmd

import (
	"context"
	"coscli/util"
	"fmt"
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tencentyun/cos-go-sdk-v5"
	"reflect"
	"testing"
)

func TestBucketPolicyCmd(t *testing.T) {
	fmt.Println("TestBucketPolicyCmd")
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

	Convey("test coscli bucket_policy", t, func() {
		Convey("success", func() {
			Convey("put", func() {
				clearCmd()
				var c *cos.BucketService
				patches := ApplyMethodFunc(reflect.TypeOf(c), "PutPolicy", func(ctx context.Context, opt *cos.BucketPutPolicyOptions) (*cos.Response, error) {
					return nil, nil
				})
				defer patches.Reset()
				cmd := rootCmd
				args := []string{"bucket-policy", "--method", "put",
					fmt.Sprintf("cos://%s", testAlias), "--policy", "{\"Statement\":[{\"Principal\":{\"qcs\":[\"qcs::cam::uin/2832742109:uin/100032195968\",\"qcs::cam::uin/2832742109:uin/100032195968\"]},\"Effect\":\"allow\",\"Action\":[\"name/cos:GetBucket\"],\"Resource\":[\"qcs::cos:ap-nanjing:uid/1253960454:willppantest6-1253960454/*\"],\"condition\":{\"ip_equal\":{\"qcs:ip\":[\"10.9.189.79\"]}}}],\"version\":\"2.0\"}"}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("get", func() {
				clearCmd()
				var c *cos.BucketService
				patches := ApplyMethodFunc(reflect.TypeOf(c), "GetPolicy", func(ctx context.Context) (*cos.BucketGetPolicyResult, *cos.Response, error) {
					return &cos.BucketGetPolicyResult{}, nil, nil
				})
				defer patches.Reset()
				cmd := rootCmd
				args := []string{"bucket-policy", "--method", "get",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("delete", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"bucket-policy", "--method", "delete",
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
				args := []string{"bucket-policy", "--method", "get",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
		})
	})
}
