package cmd

import (
	"context"
	"coscli/util"
	"fmt"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var mbCmd = &cobra.Command{
	Use:   "mb",
	Short: "Create bucket",
	Long: `Create bucket

Format:
  ./coscli mb cos://<bucket-name>-<appid> -e <endpoint>

Example:
  ./coscli mb cos://examplebucket-1234567890 -e cos.ap-beijing.myqcloud.com`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		bucketIDName, cosPath := util.ParsePath(args[0])
		if bucketIDName == "" || cosPath != "" {
			return fmt.Errorf("Invalid arguments! ")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := createBucket(cmd, args)
		return err
	},
}

func init() {
	rootCmd.AddCommand(mbCmd)

	mbCmd.Flags().StringP("region", "r", "", "Region")
	mbCmd.Flags().BoolP("ofs", "o", false, "Ofs")
	mbCmd.Flags().BoolP("maz", "m", false, "Maz")
	mbCmd.Flags().String("acl", "", "Defines the Access Control List (ACL) property of an bucket. The default value is default.")
	mbCmd.Flags().String("grant-read", "", "Grants the grantee permission to read the bucket. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	mbCmd.Flags().String("grant-write", "", "Grants the grantee permission to write the object. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id='100000000001',id=\"100000000002\".")
	mbCmd.Flags().String("grant-read-acp", "", "Grants the grantee permission to read the bucket's Access Control List (ACL). The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	mbCmd.Flags().String("grant-write-acp", "", "Grants the grantee permission to write the bucket's Access Control List (ACL). The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	mbCmd.Flags().String("grant-full-control", "", "Grants the grantee full permissions to operate on the bucket. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	mbCmd.Flags().String("tags", "", "The set of tags for the bucket, with a maximum of 10 tags (e.g., Key1=Value1 & Key2=Value2). The Key and Value in the tag set must be URL-encoded beforehand.")
}

func createBucket(cmd *cobra.Command, args []string) error {
	flagRegion, _ := cmd.Flags().GetString("region")
	flagOfs, _ := cmd.Flags().GetBool("ofs")
	flagMaz, _ := cmd.Flags().GetBool("maz")
	acl, _ := cmd.Flags().GetString("acl")
	grantRead, _ := cmd.Flags().GetString("grant-read")
	grantWrite, _ := cmd.Flags().GetString("grant-write")
	grantReadACP, _ := cmd.Flags().GetString("grant-read-acp")
	grantWriteACP, _ := cmd.Flags().GetString("grant-write-acp")
	grantFullControl, _ := cmd.Flags().GetString("grant-full-control")
	tags, _ := cmd.Flags().GetString("tags")
	if param.Endpoint == "" && flagRegion != "" {
		param.Endpoint = fmt.Sprintf("cos.%s.myqcloud.com", flagRegion)
	}
	bucketIDName, _ := util.ParsePath(args[0])

	c, err := util.CreateClient(&config, &param, bucketIDName)
	if err != nil {
		return err
	}

	// 解析tags
	tags, err = util.EncodeTagging(tags)
	if err != nil {
		return err
	}

	opt := &cos.BucketPutOptions{
		XCosACL:                   acl,
		XCosGrantRead:             grantRead,
		XCosGrantWrite:            grantWrite,
		XCosGrantFullControl:      grantFullControl,
		XCosGrantReadACP:          grantReadACP,
		XCosGrantWriteACP:         grantWriteACP,
		CreateBucketConfiguration: &cos.CreateBucketConfiguration{},
		XCosTagging:               tags,
	}

	if flagOfs {
		opt.CreateBucketConfiguration.BucketArchConfig = "OFS"
	}

	if flagMaz {
		opt.CreateBucketConfiguration.BucketAZConfig = "MAZ"
	}

	_, err = c.Bucket.Put(context.Background(), opt)
	if err != nil {
		return err
	}
	logger.Infof("Create a new bucket! name: %s\n", bucketIDName)
	return nil
}
