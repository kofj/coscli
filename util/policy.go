package util

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	logger "github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"os"
	"strings"
)

func PutBucketPolicy(c *cos.Client, policy string) error {
	configurationContent, err := GetContent(policy)
	if err != nil {
		return err
	}
	var opt cos.BucketPutPolicyOptions
	err = ParseContent(configurationContent, &opt)
	if err != nil {
		return err
	}
	_, err = c.Bucket.PutPolicy(context.Background(), &opt)

	if err != nil {
		return err
	}
	logger.Info("Put Bucket Policy Success")
	return nil
}

func GetBucketPolicy(c *cos.Client) error {
	var err error
	var policy *cos.BucketGetPolicyResult
	policy, _, err = c.Bucket.GetPolicy(context.Background())

	if err != nil {
		return err
	}
	// 渲染表格
	renderPolicyTable(policy)
	return nil
}

func DeleteBucketPolicy(c *cos.Client) error {
	_, err := c.Bucket.DeletePolicy(context.Background())
	if err != nil {
		return err
	}
	logger.Info("Delete Bucket Policy Success")
	return err
}

func renderPolicyTable(policy *cos.BucketGetPolicyResult) {
	if policy == nil {
		fmt.Println("No policy found")
		return
	}

	// 创建表格
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Section", "Key", "Value"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetCaption(true, "Bucket Policy Information")

	// 添加策略版本信息
	table.Append([]string{"Policy", "Version", policy.Version})

	// 添加分隔行
	table.Append([]string{"", "", ""})

	// 遍历所有策略声明
	for i, statement := range policy.Statement {
		section := fmt.Sprintf("Statement #%d", i+1)

		// 添加声明ID
		table.Append([]string{section, "SID", statement.Sid})

		// 添加效果
		table.Append([]string{section, "Effect", statement.Effect})

		// 添加主体信息
		if len(statement.Principal) > 0 {
			principalStr := formatPrincipal(statement.Principal)
			table.Append([]string{section, "Principal", principalStr})
		}

		// 添加操作信息
		if len(statement.Action) > 0 {
			actionStr := strings.Join(statement.Action, "\n")
			table.Append([]string{section, "Action", actionStr})
		}

		// 添加资源信息
		if len(statement.Resource) > 0 {
			resourceStr := strings.Join(statement.Resource, "\n")
			table.Append([]string{section, "Resource", resourceStr})
		}

		// 添加条件信息
		if len(statement.Condition) > 0 {
			conditionStr := formatCondition(statement.Condition)
			table.Append([]string{section, "Condition", conditionStr})
		}

		// 添加分隔行（最后一个不添加）
		if i < len(policy.Statement)-1 {
			table.Append([]string{"", "", ""})
		}
	}

	// 渲染表格
	table.Render()
}

// 格式化主体信息
func formatPrincipal(principal map[string][]string) string {
	var builder strings.Builder
	for key, values := range principal {
		builder.WriteString(fmt.Sprintf("%s:\n", key))
		for _, value := range values {
			builder.WriteString(fmt.Sprintf("  - %s\n", value))
		}
	}
	return builder.String()
}

// 格式化条件信息
func formatCondition(condition map[string]map[string]interface{}) string {
	var builder strings.Builder
	for conditionType, conditionMap := range condition {
		builder.WriteString(fmt.Sprintf("%s:\n", conditionType))
		for key, value := range conditionMap {
			switch v := value.(type) {
			case string:
				builder.WriteString(fmt.Sprintf("  %s: %s\n", key, v))
			case []interface{}:
				builder.WriteString(fmt.Sprintf("  %s:\n", key))
				for _, item := range v {
					builder.WriteString(fmt.Sprintf("    - %v\n", item))
				}
			default:
				builder.WriteString(fmt.Sprintf("  %s: %v\n", key, v))
			}
		}
	}
	return builder.String()
}
