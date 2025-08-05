package util

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tencentyun/cos-go-sdk-v5"
	"os"
)

func PutBucketEncryption(c *cos.Client, bucketEncryptionSettings BucketEncryptionSettings) error {
	var err error
	opt := &cos.BucketPutEncryptionOptions{
		Rule: &cos.BucketEncryptionConfiguration{
			SSEAlgorithm:   bucketEncryptionSettings.SSEAlgorithm,
			KMSMasterKeyID: bucketEncryptionSettings.KMSMasterKeyID,
		},
	}

	_, err = c.Bucket.PutEncryption(context.Background(), opt)

	if err != nil {
		return err
	}

	return nil
}

func GetBucketEncryption(c *cos.Client) error {
	var err error
	var encryption *cos.BucketGetEncryptionResult
	encryption, _, err = c.Bucket.GetEncryption(context.Background())

	if err != nil {
		return err
	}
	// 渲染表格
	renderEncryptionTable(encryption)
	return nil
}

func DeleteBucketEncryption(c *cos.Client) error {
	_, err := c.Bucket.DeleteEncryption(context.Background())
	return err
}

func renderEncryptionTable(encryption *cos.BucketGetEncryptionResult) {
	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Section", "Key", "Value"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCaption(true, "COS Bucket Encryption Configuration")

	// Check if encryption is configured
	if encryption == nil || encryption.Rule == nil {
		table.Append([]string{"Encryption", "Status", "Not Configured"})
		table.Render()
		return
	}

	rule := encryption.Rule

	// Add encryption algorithm information
	table.Append([]string{"Encryption", "Algorithm", rule.SSEAlgorithm})

	// Add KMS Key ID if exists
	if rule.KMSMasterKeyID != "" {
		table.Append([]string{"Encryption", "KMS Key ID", rule.KMSMasterKeyID})
	} else {
		table.Append([]string{"Encryption", "KMS Key ID", "Not Specified"})
	}

	// Add encryption status
	status := "Enabled"
	if rule.SSEAlgorithm == "" {
		status = "Disabled"
	}
	table.Append([]string{"Encryption", "Status", status})

	// Render table
	table.Render()

	// Add additional information
	printEncryptionDetails(rule)
}

// Print encryption details
func printEncryptionDetails(rule *cos.BucketEncryptionConfiguration) {
	fmt.Println("\nEncryption Details:")

	switch rule.SSEAlgorithm {
	case "AES256":
		fmt.Println(" - Type: Server-Side Encryption with COS-Managed Keys (SSE-COS)")
		fmt.Println(" - Description: Tencent Cloud COS manages encryption keys")
		fmt.Println(" - Security: AES-256 encryption algorithm")
	case "cos/kms":
		fmt.Println(" - Type: Server-Side Encryption with KMS-Managed Keys (SSE-KMS)")
		fmt.Println(" - Description: Tencent Cloud Key Management System (KMS) manages encryption keys")
		if rule.KMSMasterKeyID != "" {
			fmt.Printf(" - KMS Key ID: %s\n", rule.KMSMasterKeyID)
			fmt.Println(" - Key Type: Customer Master Key (CMK)")
		} else {
			fmt.Println(" - KMS Key ID: Default COS Managed Key")
			fmt.Println(" - Key Type: COS Managed Key")
		}
	default:
		fmt.Printf(" - Algorithm: %s\n", rule.SSEAlgorithm)
		fmt.Println(" - Note: Unknown encryption algorithm, please verify configuration")
	}
}
