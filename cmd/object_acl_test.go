package cmd

import (
	"coscli/util"
	"fmt"
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tencentyun/cos-go-sdk-v5"
	"testing"
	"time"
)

func TestObjectAclCmd(t *testing.T) {
	fmt.Println("TestObjectAclCmd")
	testBucket = randStr(8)
	testAlias = testBucket + "-alias"
	setUp(testBucket, testAlias, testEndpoint, false, false)
	defer tearDown(testBucket, testAlias, testEndpoint, false)
	setUp(testOfsBucket, testOfsBucketAlias, testEndpoint, true, false)
	defer tearDown(testOfsBucket, testOfsBucketAlias, testEndpoint, false)
	clearCmd()
	cmd := rootCmd
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	genDir(testDir, 3)
	defer delDir(testDir)
	localFileName := fmt.Sprintf("%s/small-file", testDir)
	// 上传cos文件
	cosFileName := fmt.Sprintf("cos://%s/%s", testAlias, "multi-small")
	args := []string{"cp", localFileName, cosFileName, "-r"}
	cmd.SetArgs(args)
	cmd.Execute()

	// 上传ofs文件
	ofsFileName := fmt.Sprintf("cos://%s/%s", testOfsBucketAlias, "multi-small")
	args = []string{"cp", localFileName, ofsFileName, "-r"}
	cmd.SetArgs(args)
	cmd.Execute()
	Convey("test coscli object_acl", t, func() {
		Convey("success", func() {
			Convey("cos", func() {
				Convey("put", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-acl", "--method", "put",
						cosFileName, "--grant-read", "id=\"100000000003\",id=\"100000000002\""}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
				Convey("get", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-acl", "--method", "get",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
			})
			Convey("ofs", func() {
				Convey("put", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-acl", "--method", "put",
						ofsFileName, "--grant-read", "id=\"100000000003\",id=\"100000000002\""}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
				Convey("get", func() {
					time.Sleep(time.Second)
					clearCmd()
					cmd := rootCmd
					args := []string{"object-acl", "--method", "get",
						ofsFileName}
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
					args := []string{"object-acl", "--method", "put",
						cosFileName, "--grant-read", "id=\"100000000003\",id=\"100000000002\""}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
			})
		})
	})
}
