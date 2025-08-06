package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/spf13/cobra"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Modify inventory",
	Long: `Modify inventory

Format:
	./coscli inventory --method [method] cos://<bucket-name> 

Example:
	./coscli inventory --method put cos://examplebucket
	./coscli inventory --method get cos://examplebucket
	./coscli inventory --method list cos://examplebucket
	./coscli inventory --method delete cos://examplebucket
	./coscli inventory --method post cos://examplebucket`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		method, _ := cmd.Flags().GetString("method")
		taskId, _ := cmd.Flags().GetString("task-id")
		configuration, _ := cmd.Flags().GetString("configuration")

		cosUrl, err := util.FormatUrl(args[0])
		if err != nil {
			return fmt.Errorf("cos url format error:%v", err)
		}

		if !cosUrl.IsCosUrl() {
			return fmt.Errorf("cospath needs to contain cos://")
		}

		bucketName := cosUrl.(*util.CosUrl).Bucket

		c, err := util.NewClient(&config, &param, bucketName)
		if err != nil {
			return err
		}

		if method == "put" {
			err = util.PutBucketInventory(c, taskId, configuration)
		} else if method == "get" {
			err = util.GetBucketInventory(c, taskId)
		} else if method == "list" {
			err = util.ListBucketInventory(c)
		} else if method == "delete" {
			err = util.DeleteBucketInventory(c, taskId)
		} else if method == "post" {
			err = util.PostBucketInventory(c, taskId, configuration)
		} else {
			err = fmt.Errorf("method '%s' is not supported, valid methods are 'put', 'get', 'list', 'delete', and 'post'", method)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(inventoryCmd)
	inventoryCmd.Flags().String("method", "", "put/get")
	inventoryCmd.Flags().String("task-id", "", "The name of the checklist task.Default value: None.Valid characters: a-z, A-Z, 0-9, -, _, .")
	inventoryCmd.Flags().String("configuration", "", "Store checklist configuration information, supporting both XML and JSON syntax; if the input contains the file:// prefix, it indicates reading the configuration from a file.")
}
