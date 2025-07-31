package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/spf13/cobra"
)

var objectAclCmd = &cobra.Command{
	Use:   "object-acl",
	Short: "Modify object acl",
	Long: `Modify object acl

Format:
	./coscli object-acl --method [method] cos://<bucket-name>

Example:
	./coscli object-acl --method put cos://examplebucket/exampleobject  --grant-read="id=\"100000000003\",id=\"100000000002\""
	./coscli object-acl --method get cos://examplebucket/exampleobject`,

	RunE: func(cmd *cobra.Command, args []string) error {
		var aclSettings util.ACLSettings
		method, _ := cmd.Flags().GetString("method")
		versionId, _ := cmd.Flags().GetString("version-id")
		aclSettings.ACL, _ = cmd.Flags().GetString("acl")
		aclSettings.GrantRead, _ = cmd.Flags().GetString("grant-read")
		//aclSettings.GrantWrite, _ = cmd.Flags().GetString("grant-write")
		aclSettings.GrantReadACP, _ = cmd.Flags().GetString("grant-read-acp")
		aclSettings.GrantWriteACP, _ = cmd.Flags().GetString("grant-write-acp")
		aclSettings.GrantFullControl, _ = cmd.Flags().GetString("grant-full-control")

		var err error
		cosPath := args[0]
		if !util.IsCosPath(cosPath) {
			return fmt.Errorf("cospath needs to contain cos://")
		}

		bucketName, object := util.ParsePath(cosPath)
		c, err := util.NewClient(&config, &param, bucketName)
		if err != nil {
			return err
		}

		bucketType, err := util.GetBucketType(c, &param, &config, bucketName)
		if err != nil {
			return err
		}

		if method == "put" {
			err = util.PutObjectAcl(c, object, versionId, bucketType, aclSettings)
		}

		if method == "get" {
			err = util.GetObjectAcl(c, object, versionId, bucketType)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(objectAclCmd)
	objectAclCmd.Flags().String("method", "", "put/get")
	objectAclCmd.Flags().String("version-id", "", "put/get acl of a file , only available if bucket versioning is enabled.")
	objectAclCmd.Flags().String("acl", "", "Defines the Access Control List (ACL) property of an object. The default value is default.")
	objectAclCmd.Flags().String("grant-read", "", "Grants the grantee permission to read the object. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	//objectAclCmd.Flags().String("grant-write", "", "Grants the grantee permission to write the object. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id='100000000001',id=\"100000000002\".")
	objectAclCmd.Flags().String("grant-read-acp", "", "Grants the grantee permission to read the object's Access Control List (ACL). The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	objectAclCmd.Flags().String("grant-write-acp", "", "Grants the grantee permission to write the object's Access Control List (ACL). The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
	objectAclCmd.Flags().String("grant-full-control", "", "Grants the grantee full permissions to operate on the object. The format is id=\"[OwnerUin]\", for example, id=\"100000000001\". Multiple grantees can be specified using commas (,), for example, id=\"100000000001\",id=\"100000000002\".")
}
