package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/spf13/cobra"
)

var bucketEncryptionCmd = &cobra.Command{
	Use:   "bucket-encryption",
	Short: "Modify bucket encryption",
	Long: `Modify bucket encryption

Format:
	./coscli bucket-encryption --method [method] cos://<bucket-name>

Example:
	./coscli bucket-encryption --method put cos://examplebucket  --grant-read="id=\"100000000003\",id=\"100000000002\""
	./coscli bucket-encryption --method get cos://examplebucket
	./coscli bucket-encryption --method delete cos://examplebucket`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var bucketEncryptionSettings util.BucketEncryptionSettings
		method, _ := cmd.Flags().GetString("method")
		bucketEncryptionSettings.SSEAlgorithm, _ = cmd.Flags().GetString("sse-algorithm")
		bucketEncryptionSettings.KMSMasterKeyID, _ = cmd.Flags().GetString("kms-master-key-id")
		bucketEncryptionSettings.KMSAlgorithm, _ = cmd.Flags().GetString("kms-algorithm")

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
			err = util.PutBucketEncryption(c, bucketEncryptionSettings)
		} else if method == "get" {
			err = util.GetBucketEncryption(c)
		} else if method == "delete" {
			err = util.DeleteBucketEncryption(c)
		} else {
			err = fmt.Errorf("method '%s' is not supported, valid methods are 'put', 'get', and 'delete'", method)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(bucketEncryptionCmd)
	bucketEncryptionCmd.Flags().String("method", "", "put/get/delete")
	bucketEncryptionCmd.Flags().String("sse-algorithm", "", "Support enumeration values: AES256, SM4, KMS. AES256 represents the use of the SSE-COS mode with the AES256 encryption algorithm. SM4 represents the use of the SSE-COS mode with the SM4 encryption algorithm. KMS represents the SSE-KMS mode.")
	bucketEncryptionCmd.Flags().String("kms-master-key-id", "", "When the value of SSEAlgorithm is KMS, it is used to specify the user's Customer Master Key (CMK) for KMS. If not specified, the default CMK created by COS will be used.")
	bucketEncryptionCmd.Flags().String("kms-algorithm", "", "When the value of SSEAlgorithm is KMS, it is used to specify the encryption algorithm for KMS. Supported enumeration values are AES256 and SM4. If not specified, the default value is AES256.  AES256 represents the encryption algorithm AES256. SM4 represents the encryption algorithm SM4.")
}
