package cmd

import (
	"context"
	"coscli/util"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var signurlCmd = &cobra.Command{
	Use:   "signurl",
	Short: "Gets the signed download URL",
	Long: `Gets the signed download URL

Format:
  ./coscli signurl cos://<bucket-name>/<key> [flags]

Example:
  ./coscli signurl cos://examplebucket/test.jpg -t 100`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		time, _ := cmd.Flags().GetInt("time")
		mode, _ := cmd.Flags().GetString("mode")
		var err error
		if util.IsCosPath(args[0]) {
			err = GetSignedURL(args[0], time, mode)
		} else {
			return fmt.Errorf("cospath needs to contain cos://")
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(signurlCmd)

	signurlCmd.Flags().IntP("time", "t", 10000, "Set the validity time of the signature(Default 10000)")
	signurlCmd.Flags().StringP("mode", "m", "", "Set output mode")
}

// GetSignedURL 生成签名url
func GetSignedURL(path string, t int, mode string) error {
	bucketName, cosPath := util.ParsePath(path)
	c, err := util.NewClient(&config, &param, bucketName)
	if err != nil {
		return err
	}

	opt := &cos.PresignedURLOptions{
		Query:  &url.Values{},
		Header: &http.Header{},
	}

	presignedURL, err := c.Object.GetPresignedURL2(context.Background(), http.MethodGet, cosPath, time.Second*time.Duration(t), opt)
	if err != nil {
		return err
	}

	if mode == "simple" {
		fmt.Println(presignedURL)
	} else {
		logger.Infoln("Signed URL:")
		logger.Infoln(presignedURL)
	}

	return nil
}
