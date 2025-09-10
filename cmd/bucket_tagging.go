package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/spf13/cobra"
)

var bucketTaggingCmd = &cobra.Command{
	Use:   "bucket-tagging",
	Short: "Modify bucket tagging",
	Long: `Modify bucket tagging

Format:
	./coscli bucket-tagging --method [method] cos://<bucket-name> [tag_key]#[tag_value]

Example:
	./coscli bucket-tagging --method put cos://examplebucket tag1#test1 tag2#test2
    ./coscli bucket-tagging --method add cos://examplebucket tag3#test3
	./coscli bucket-tagging --method get cos://examplebucket
	./coscli bucket-tagging --method delete cos://examplebucket
	./coscli bucket-tagging --method delete cos://examplebucket tag1#test1 tag2#test2`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		method, _ := cmd.Flags().GetString("method")

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
			if len(args) < 2 {
				return fmt.Errorf("not enough arguments in call to put bucket tagging")
			}
			err = util.PutBucketTagging(c, args[1:])
		} else if method == "add" {
			err = util.AddBucketTagging(c, args[1:])
		} else if method == "get" {
			err = util.GetBucketTagging(c)
		} else if method == "delete" {
			if len(args) == 1 {
				err = util.DeleteBucketTagging(c)
			} else {
				err = util.DeleteDesBucketTagging(c, args[1:])
			}
		} else {
			err = fmt.Errorf("method '%s' is not supported, valid methods are 'put','add','get', and 'delete'", method)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(bucketTaggingCmd)
	bucketTaggingCmd.Flags().String("method", "", "put/add/get/delete")
}
