package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/spf13/cobra"
)

var bucketPolicyCmd = &cobra.Command{
	Use:   "bucket-policy",
	Short: "Modify bucket policy",
	Long: `Modify bucket policy

Format:
	./coscli bucket-policy --method [method] cos://<bucket-name>

Example:
	./coscli bucket-policy --method put cos://examplebucket  --policy="{\"Statement\":[{\"Principal\":{\"qcs\":[\"qcs::cam::uin/100000000001:uin/100000000011\"]},\"Effect\":\"allow\",\"Action\":[\"name/cos:GetBucket\"],\"Resource\":[\"qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/*\"]}],\"version\":\"2.0\"}"
	./coscli bucket-policy --method get cos://examplebucket
    ./coscli bucket-policy --method delete cos://examplebucket`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		method, _ := cmd.Flags().GetString("method")
		policy, _ := cmd.Flags().GetString("policy")

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
			if policy == "" {
				return fmt.Errorf("no policy provided")
			}
			err = util.PutBucketPolicy(c, policy)
		} else if method == "get" {
			err = util.GetBucketPolicy(c)
		} else if method == "delete" {
			err = util.DeleteBucketPolicy(c)
		} else {
			err = fmt.Errorf("method '%s' is not supported, valid methods are 'put', 'get', and 'delete'", method)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(bucketPolicyCmd)
	bucketPolicyCmd.Flags().String("method", "", "put/get/delete")
	bucketPolicyCmd.Flags().String("policy", "", "Bucket policy (JSON format).")
}
