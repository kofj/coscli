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

func TestBucketInventoryCmd(t *testing.T) {
	fmt.Println("TestBucketInventoryCmd")
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

	Convey("test coscli bucket_inventory", t, func() {
		Convey("success", func() {
			Convey("put", func() {
				clearCmd()
				var c *cos.BucketService
				patches := ApplyMethodFunc(reflect.TypeOf(c), "PutInventory", func(ctx context.Context, id string, opt *cos.BucketPutInventoryOptions) (*cos.Response, error) {
					return nil, nil
				})
				defer patches.Reset()
				cmd := rootCmd
				args := []string{"inventory", "--method", "put",
					fmt.Sprintf("cos://%s", testAlias), "--task-id", "list4", "--configuration", "<?xml version=\"1.0\" encoding=\"UTF-8\"?><InventoryConfiguration xmlns=\"http://....\"><Id>list4</Id><IsEnabled>false</IsEnabled><Destination><COSBucketDestination><Format>CSV</Format><AccountId>100000000002</AccountId><Bucket>qcs::cos:ap-nanjing::test-1000000001</Bucket><Prefix>list4</Prefix><Encryption><SSE-COS></SSE-COS></Encryption></COSBucketDestination></Destination><Schedule><Frequency>Weekly</Frequency></Schedule><Filter><And><Prefix>myPrefix</Prefix><Tag><Key>age</Key><Value>18</Value></Tag></And><Period><StartTime>1768688761</StartTime><EndTime>1568688762</EndTime></Period></Filter><IncludedObjectVersions>All</IncludedObjectVersions><OptionalFields><Field>Size</Field><Field>Tag</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>StorageClass</Field><Field>IsMultipartUploaded</Field></OptionalFields></InventoryConfiguration>"}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("get", func() {
				clearCmd()
				var c *cos.BucketService
				patches := ApplyMethodFunc(reflect.TypeOf(c), "GetInventory", func(ctx context.Context, id string) (*cos.BucketGetInventoryResult, *cos.Response, error) {
					return &cos.BucketGetInventoryResult{}, nil, nil
				})
				defer patches.Reset()
				cmd := rootCmd
				args := []string{"inventory", "--method", "get",
					fmt.Sprintf("cos://%s", testAlias), "--task-id", "list4"}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("list", func() {
				clearCmd()
				var c *cos.BucketService
				patches := ApplyMethodFunc(reflect.TypeOf(c), "ListInventoryConfigurations", func(ctx context.Context, token string) (*cos.ListBucketInventoryConfigResult, *cos.Response, error) {
					return &cos.ListBucketInventoryConfigResult{}, nil, nil
				})
				defer patches.Reset()
				cmd := rootCmd
				args := []string{"inventory", "--method", "list",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("delete", func() {
				clearCmd()
				var c *cos.BucketService
				patches := ApplyMethodFunc(reflect.TypeOf(c), "DeleteInventory", func(ctx context.Context, id string) (*cos.Response, error) {
					return nil, nil
				})
				defer patches.Reset()
				cmd := rootCmd
				args := []string{"inventory", "--method", "delete",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
			Convey("post", func() {
				clearCmd()
				var c *cos.BucketService
				patches := ApplyMethodFunc(reflect.TypeOf(c), "PostInventory", func(ctx context.Context, id string, opt *cos.BucketPostInventoryOptions) (*cos.Response, error) {
					return nil, nil
				})
				defer patches.Reset()
				cmd := rootCmd
				args := []string{"inventory", "--method", "post",
					fmt.Sprintf("cos://%s", testAlias), "--task-id", "list4", "--configuration", "<?xml version=\"1.0\" encoding=\"UTF-8\"?><InventoryConfiguration xmlns=\"http://....\"><Id>list4</Id><Destination><COSBucketDestination><Format>CSV</Format><AccountId>100000000002</AccountId><Bucket>qcs::cos:ap-nanjing::test-100000001</Bucket><Prefix>list4</Prefix><Encryption><SSE-COS>111</SSE-COS></Encryption></COSBucketDestination></Destination><Filter><And><Prefix>myPrefix</Prefix><Tag><Key>age</Key><Value>18</Value></Tag></And><Period><StartTime>1768688761</StartTime><EndTime>1568688762</EndTime></Period></Filter><IncludedObjectVersions>All</IncludedObjectVersions><OptionalFields><Field>Size</Field><Field>Tag</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>StorageClass</Field><Field>IsMultipartUploaded</Field></OptionalFields></InventoryConfiguration>"}
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
				args := []string{"inventory", "--method", "list",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("FormatUrl err", func() {
				clearCmd()
				cmd := rootCmd
				patches := ApplyFunc(util.FormatUrl, func(urlStr string) (util.StorageUrl, error) {
					return nil, fmt.Errorf("test format url error")
				})
				defer patches.Reset()
				args := []string{"inventory", "--method", "list",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("cos path error", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"inventory", "--method", "list",
					fmt.Sprintf("cos:/%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("invalid method", func() {
				clearCmd()
				cmd := rootCmd

				args := []string{"inventory", "--method", "add",
					fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
		})
	})
}
