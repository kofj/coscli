package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/spf13/cobra"
)

var objectTaggingCmd = &cobra.Command{
	Use:   "object-tagging",
	Short: "Modify object tagging",
	Long: `Modify object tagging

Format:
	./coscli object-tagging --method [method] cos://<bucket-name> [tag_key]#[tag_value]

Example:
	./coscli object-tagging --method put cos://examplebucket/exampleobject tag1#test1 tag2#test2
    ./coscli object-tagging --method add cos://examplebucket/exampleobject tag3#test3
	./coscli object-tagging --method get cos://examplebucket/exampleobject
	./coscli object-tagging --method delete cos://examplebucket/exampleobject
	./coscli object-tagging --method delete cos://examplebucket/exampleobject tag1#test1 tag2#test2`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		method, _ := cmd.Flags().GetString("method")
		versionId, _ := cmd.Flags().GetString("version-id")

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

		if bucketType == util.BucketTypeOfs {
			return fmt.Errorf("ofs bucket not Implemented")
		}

		if method == "put" {
			if len(args) < 2 {
				return fmt.Errorf("not enough arguments in call to put object tagging")
			}
			err = util.PutObjectTagging(c, object, args[1:], versionId, bucketType)
		} else if method == "add" {
			if len(args) < 2 {
				return fmt.Errorf("not enough arguments in call to get object tagging")
			}
			err = util.AddObjectTagging(c, object, args[1:], versionId, bucketType)
		} else if method == "get" {
			if len(args) < 1 {
				return fmt.Errorf("not enough arguments in call to get object tagging")
			}
			err = util.GetObjectTagging(c, object, versionId, bucketType)
		} else if method == "delete" {
			if len(args) < 1 {
				return fmt.Errorf("not enough arguments in call to delete object tagging")
			} else if len(args) == 1 {
				err = util.DeleteObjectTagging(c, object, versionId, bucketType)
			} else {
				err = util.DeleteDesObjectTagging(c, object, args[1:], versionId, bucketType)
			}
		} else {
			err = fmt.Errorf("method '%s' is not supported, valid methods are 'put','add','get', and 'delete'", method)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(objectTaggingCmd)
	objectTaggingCmd.Flags().String("method", "", "put/add/get/delete")
	objectTaggingCmd.Flags().String("version-id", "", "tagging a specified version of a file , only available if bucket versioning is enabled.")
}
