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
	"time"
)

func TestObjectTaggingCmd(t *testing.T) {
	fmt.Println("TestObjectTaggingCmd")
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
	localFileName := fmt.Sprintf("%s/small-file/0", testDir)
	// 上传cos文件
	cosFileName := fmt.Sprintf("cos://%s/%s", testAlias, "multi-small")
	args := []string{"cp", localFileName, cosFileName, "-r"}
	cmd.SetArgs(args)
	cmd.Execute()

	Convey("test coscli object_tagging", t, func() {
		Convey("success", func() {
			Convey("cos", func() {
				Convey("put", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "put",
						cosFileName, "testkey#testval"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
				Convey("add", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "add",
						cosFileName, "testkey2#testval2"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
				Convey("get", func() {
					time.Sleep(time.Second)
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "get",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
				Convey("deleteDes", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "delete",
						cosFileName, "testkey2#testval2"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
				Convey("delete", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "delete",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeNil)
				})
			})
			//Convey("ofs", func() {
			//	Convey("put", func() {
			//		clearCmd()
			//		cmd := rootCmd
			//		args := []string{"object-tagging", "--method", "put",
			//			ofsFileName, "testkey#testval"}
			//		cmd.SetArgs(args)
			//		e := cmd.Execute()
			//		So(e, ShouldBeNil)
			//	})
			//	Convey("add", func() {
			//		clearCmd()
			//		cmd := rootCmd
			//		args := []string{"object-tagging", "--method", "add",
			//			ofsFileName, "testkey2#testval2"}
			//		cmd.SetArgs(args)
			//		e := cmd.Execute()
			//		So(e, ShouldBeNil)
			//	})
			//	Convey("get", func() {
			//		time.Sleep(time.Second)
			//		clearCmd()
			//		cmd := rootCmd
			//		args := []string{"object-tagging", "--method", "get",
			//			ofsFileName}
			//		cmd.SetArgs(args)
			//		e := cmd.Execute()
			//		So(e, ShouldBeNil)
			//	})
			//	Convey("deleteDes", func() {
			//		clearCmd()
			//		cmd := rootCmd
			//		args := []string{"object-tagging", "--method", "delete",
			//			ofsFileName, "testkey2#testval2"}
			//		cmd.SetArgs(args)
			//		e := cmd.Execute()
			//		So(e, ShouldBeNil)
			//	})
			//	Convey("delete", func() {
			//		clearCmd()
			//		cmd := rootCmd
			//		args := []string{"object-tagging", "--method", "delete",
			//			ofsFileName}
			//		cmd.SetArgs(args)
			//		e := cmd.Execute()
			//		So(e, ShouldBeNil)
			//	})
			//})
		})
		Convey("fail", func() {
			Convey("put", func() {
				Convey("not enough arguments", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "put",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("clinet err", func() {
					clearCmd()
					cmd := rootCmd
					patches := ApplyFunc(util.NewClient, func(config *util.Config, param *util.Param, bucketName string) (client *cos.Client, err error) {
						return nil, fmt.Errorf("test put client error")
					})
					defer patches.Reset()
					args := []string{"object-tagging", "--method", "put",
						cosFileName, "testval"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("invalid tag", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "put",
						cosFileName, "testval"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("PutTagging failed", func() {
					clearCmd()
					var c *cos.ObjectService
					patches := ApplyMethodFunc(reflect.TypeOf(c), "PutTagging", func(ctx context.Context, name string, opt *cos.ObjectPutTaggingOptions, id ...string) (*cos.Response, error) {
						return nil, fmt.Errorf("PutTagging failed")
					})
					defer patches.Reset()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "put",
						cosFileName, "qcs:1#testval"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
			})
			Convey("add", func() {
				Convey("not enough arguments", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "add",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("clinet err", func() {
					clearCmd()
					cmd := rootCmd
					patches := ApplyFunc(util.NewClient, func(config *util.Config, param *util.Param, bucketName string) (client *cos.Client, err error) {
						return nil, fmt.Errorf("test add client error")
					})
					defer patches.Reset()
					args := []string{"object-tagging", "--method", "add",
						cosFileName, "testval"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("invalid tag", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "add",
						cosFileName, "testval"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("AddTagging failed", func() {
					clearCmd()
					var c *cos.ObjectService
					patches := ApplyMethodFunc(reflect.TypeOf(c), "PutTagging", func(ctx context.Context, name string, opt *cos.ObjectPutTaggingOptions, id ...string) (*cos.Response, error) {
						return nil, fmt.Errorf("PutTagging failed")
					})
					defer patches.Reset()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "add",
						cosFileName, "qcs:1#testval"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
			})
			Convey("get", func() {
				Convey("not enough arguments", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "get"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("clinet err", func() {
					clearCmd()
					cmd := rootCmd
					patches := ApplyFunc(util.NewClient, func(config *util.Config, param *util.Param, bucketName string) (client *cos.Client, err error) {
						return nil, fmt.Errorf("test get client error")
					})
					defer patches.Reset()
					args := []string{"object-tagging", "--method", "get",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("get tag error", func() {
					clearCmd()
					var c *cos.ObjectService
					patches := ApplyMethodFunc(reflect.TypeOf(c), "GetTagging", func(ctx context.Context, name string, opt ...interface{}) (*cos.ObjectGetTaggingResult, *cos.Response, error) {
						return nil, nil, fmt.Errorf("GetTagging failed")
					})
					defer patches.Reset()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "get",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
			})
			Convey("delete", func() {
				Convey("not enough arguments", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "delete"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("delete bucket not exist", func() {
					clearCmd()
					var c *cos.ObjectService
					patches := ApplyMethodFunc(reflect.TypeOf(c), "DeleteTagging", func(ctx context.Context, name string, opt ...interface{}) (*cos.Response, error) {
						return nil, fmt.Errorf("test delete tagging error")
					})
					defer patches.Reset()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "delete",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("clinet err", func() {
					clearCmd()
					cmd := rootCmd
					patches := ApplyFunc(util.NewClient, func(config *util.Config, param *util.Param, bucketName string) (client *cos.Client, err error) {
						return nil, fmt.Errorf("test delete client error")
					})
					defer patches.Reset()
					args := []string{"object-tagging", "--method", "delete",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
				Convey("DeleteTagging", func() {
					clearCmd()
					cmd := rootCmd
					var c *cos.ObjectService
					patches := ApplyMethodFunc(reflect.TypeOf(c), "DeleteTagging", func(ctx context.Context, name string, opt ...interface{}) (*cos.Response, error) {
						return nil, fmt.Errorf("test delete tagging error")
					})
					defer patches.Reset()
					args := []string{"object-tagging", "--method", "delete",
						cosFileName}
					cmd.SetArgs(args)
					e := cmd.Execute()
					fmt.Printf(" : %v", e)
					So(e, ShouldBeError)
				})
			})
			Convey("deleteDes", func() {
				Convey("invalid tag format", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "delete",
						cosFileName, "testkey2"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeError)
				})
				Convey("ObjectTagging not exist", func() {
					clearCmd()
					cmd := rootCmd
					args := []string{"object-tagging", "--method", "delete",
						cosFileName, "testkey101#11"}
					cmd.SetArgs(args)
					e := cmd.Execute()
					So(e, ShouldBeError)
				})
			})
		})
	})
}
