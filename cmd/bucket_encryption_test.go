package cmd

import (
	"coscli/util"
	"fmt"
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tencentyun/cos-go-sdk-v5"
	"testing"
)

func TestBucketEncryptionCmd(t *testing.T) {
	fmt.Println("TestBucketEncryptionCmd")
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

	Convey("test coscli bucket_encryption", t, func() {
		Convey("success", func() {
			Convey("put", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"bucket-encryption", "--method", "put",
					fmt.Sprintf("cos://%s", testAlias), "--sse-algorithm", "AES256"}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("get", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"bucket-encryption", "--method", "get",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("delete", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"bucket-encryption", "--method", "delete",
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
				args := []string{"bucket-encryption", "--method", "get",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("cos path error", func() {
				clearCmd()
				cmd := rootCmd

				args := []string{"bucket-encryption", "--method", "get",
					fmt.Sprintf("cos:/%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("invalid method", func() {
				clearCmd()
				cmd := rootCmd

				args := []string{"bucket-encryption", "--method", "add",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
		})
	})
}
