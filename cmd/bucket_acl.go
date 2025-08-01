package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/spf13/cobra"
)

var bucketAclCmd = &cobra.Command{
	Use:   "bucket-acl",
	Short: "Modify bucket acl",
	Long: `Modify bucket acl

Format:
	./coscli bucket-acl --method [method] cos://<bucket-name>

Example:
	./coscli bucket-acl --method put cos://examplebucket  --grant-read="id=\"100000000003\",id=\"100000000002\""
	./coscli bucket-acl --method get cos://examplebucket`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var aclSettings util.ACLSettings
		method, _ := cmd.Flags().GetString("method")
		aclSettings.ACL, _ = cmd.Flags().GetString("acl")
		aclSettings.GrantRead, _ = cmd.Flags().GetString("grant-read")
		aclSettings.GrantWrite, _ = cmd.Flags().GetString("grant-write")
		aclSettings.GrantReadACP, _ = cmd.Flags().GetString("grant-read-acp")
		aclSettings.GrantWriteACP, _ = cmd.Flags().GetString("grant-write-acp")
		aclSettings.GrantFullControl, _ = cmd.Flags().GetString("grant-full-control")

		var err error
		cosPath := args[0]
		if !util.IsCosPath(cosPath) {
			return fmt.Errorf("cospath needs to contain cos://")
		}

		bucketName, _ := util.ParsePath(cosPath)
		c, err := util.NewClient(&config, &param, bucketName)
		if err != nil {
			return err
		}

		if method == "put" {
			err = util.PutBucketAcl(c, aclSettings)
		} else if method == "get" {
			err = util.GetBucketAcl(c)
		} else {
			err = fmt.Errorf("method '%s' is not supported, valid methods are 'put' and 'get'", method)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(bucketAclCmd)
	bucketAclCmd.Flags().String("method", "", "put/get")
	bucketAclCmd.Flags().String("acl", "", "Defines the Access Control List (ACL) property of an bucket. The default value is default.")
	bucketAclCmd.Flags().String("grant-read", "", "Grants the grantee permission to read the bucket. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	bucketAclCmd.Flags().String("grant-write", "", "Grants the grantee permission to write the bucket. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id='100000000001',id=\"100000000002\".")
	bucketAclCmd.Flags().String("grant-read-acp", "", "Grants the grantee permission to read the bucket's Access Control List (ACL). The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	bucketAclCmd.Flags().String("grant-write-acp", "", "Grants the grantee permission to write the bucket's Access Control List (ACL). The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	bucketAclCmd.Flags().String("grant-full-control", "", "Grants the grantee full permissions to operate on the bucket. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
}
