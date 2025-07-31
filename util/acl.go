package util

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tencentyun/cos-go-sdk-v5"
	"os"
	"strings"
)

func PutBucketAcl(c *cos.Client, aclSettings ACLSettings) error {
	var err error
	opt := &cos.BucketPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL:              aclSettings.ACL,
			XCosGrantRead:        aclSettings.GrantRead,
			XCosGrantWrite:       aclSettings.GrantWrite,
			XCosGrantReadACP:     aclSettings.GrantReadACP,
			XCosGrantWriteACP:    aclSettings.GrantWriteACP,
			XCosGrantFullControl: aclSettings.GrantFullControl,
		},
	}
	_, err = c.Bucket.PutACL(context.Background(), opt)

	if err != nil {
		return err
	}

	return nil
}

func GetBucketAcl(c *cos.Client) error {
	var err error
	var acl *cos.BucketGetACLResult
	acl, _, err = c.Bucket.GetACL(context.Background())

	if err != nil {
		return err
	}
	// 渲染表格
	RenderACLTable(acl)

	return nil
}

func PutObjectAcl(c *cos.Client, object, versionId, bucketType string, aclSettings ACLSettings) error {
	var err error
	opt := &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL:       aclSettings.ACL,
			XCosGrantRead: aclSettings.GrantRead,
			//XCosGrantWrite:       aclSettings.GrantWrite,
			XCosGrantReadACP:     aclSettings.GrantReadACP,
			XCosGrantWriteACP:    aclSettings.GrantWriteACP,
			XCosGrantFullControl: aclSettings.GrantFullControl,
		},
	}
	if bucketType == BucketTypeOfs {
		_, err = c.Object.PutACL(context.Background(), object, opt)
	} else {
		_, err = c.Object.PutACL(context.Background(), object, opt, versionId)
	}

	if err != nil {
		return err
	}

	return nil
}

func GetObjectAcl(c *cos.Client, object, versionId, bucketType string) error {
	var err error
	var acl *cos.ObjectGetACLResult
	if bucketType == BucketTypeOfs {
		acl, _, err = c.Object.GetACL(context.Background(), object)
	} else {
		acl, _, err = c.Object.GetACL(context.Background(), object, versionId)
	}

	if err != nil {
		return err
	}
	// 渲染表格
	RenderACLTable(acl)

	return nil
}

func RenderACLTable(acl *cos.ACLXml) {
	// 创建表格
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Section", "Key", "Value"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetCaption(true, "Access Control List (ACL) Information")

	// 设置颜色
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)

	// 添加所有者信息
	if acl.Owner != nil {
		table.Append([]string{"Owner", "UIN", acl.Owner.UIN})
		table.Append([]string{"Owner", "ID", acl.Owner.ID})
		table.Append([]string{"Owner", "Display Name", acl.Owner.DisplayName})
	} else {
		table.Append([]string{"Owner", "", "No owner information"})
	}

	// 添加分隔行
	table.Append([]string{"", "", ""})

	// 添加 ACL 授权信息
	if len(acl.AccessControlList) == 0 {
		table.Append([]string{"ACL Grants", "", "No access grants"})
	} else {
		for i, grant := range acl.AccessControlList {
			section := fmt.Sprintf("Grant #%d", i+1)

			// 添加权限信息（带颜色）
			permissionColor := tablewriter.Colors{tablewriter.FgHiYellowColor}
			table.Rich([]string{section, "Permission", grant.Permission}, []tablewriter.Colors{
				{},
				{tablewriter.FgHiGreenColor},
				permissionColor,
			})

			// 添加授权者信息
			if grantee := grant.Grantee; grantee != nil {
				// 确定授权者类型
				granteeType := grantee.Type
				if granteeType == "" {
					granteeType = grantee.TypeAttr.Value
				}
				table.Append([]string{section, "Grantee Type", granteeType})

				// 根据类型显示不同的标识符
				switch strings.ToLower(granteeType) {
				case "group":
					table.Append([]string{section, "URI", grantee.URI})
				case "canonicaluser":
					table.Append([]string{section, "ID", grantee.ID})
					table.Append([]string{section, "Display Name", grantee.DisplayName})
				default:
					table.Append([]string{section, "UIN", grantee.UIN})
					table.Append([]string{section, "URI", grantee.URI})
					table.Append([]string{section, "ID", grantee.ID})
					table.Append([]string{section, "Display Name", grantee.DisplayName})
					table.Append([]string{section, "Sub Account", grantee.SubAccount})
				}
			} else {
				table.Append([]string{section, "Grantee", "No grantee information"})
			}

			// 添加分隔行（最后一个不添加）
			if i < len(acl.AccessControlList)-1 {
				table.Append([]string{"", "", ""})
			}
		}
	}

	// 渲染表格
	table.Render()

	// 添加摘要信息
	printACLSummary(acl)
}

// 打印 ACL 摘要
func printACLSummary(acl *cos.ACLXml) {
	fmt.Println("\nSummary:")

	if acl.Owner != nil {
		fmt.Printf(" - Owner: %s (UIN: %s)\n",
			acl.Owner.DisplayName,
			acl.Owner.UIN,
		)
	} else {
		fmt.Println(" - Owner: Not specified")
	}

	fmt.Printf(" - Total Grants: %d\n", len(acl.AccessControlList))

	// 统计权限类型
	permissionCount := make(map[string]int)
	for _, grant := range acl.AccessControlList {
		permissionCount[grant.Permission]++
	}

	fmt.Println(" - Permissions:")
	for perm, count := range permissionCount {
		fmt.Printf("   - %s: %d grants\n", perm, count)
	}
}
