package util

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
	"os"
	"strings"
)

// PutBucketTagging todo
func PutBucketTagging(c *cos.Client, tags []string) error {

	tg := &cos.BucketPutTaggingOptions{}
	for i := 0; i < len(tags); i += 1 {
		tmp := strings.Split(tags[i], "#")
		if len(tmp) >= 2 {
			tg.TagSet = append(tg.TagSet, cos.BucketTaggingTag{Key: tmp[0], Value: tmp[1]})
		} else {
			return fmt.Errorf("invalid tag")
		}
	}

	_, err := c.Bucket.PutTagging(context.Background(), tg)
	if err != nil {
		return err
	}

	return nil
}

// GetBucketTagging todo
func GetBucketTagging(c *cos.Client) error {
	v, _, err := c.Bucket.GetTagging(context.Background())
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	for _, t := range v.TagSet {
		table.Append([]string{t.Key, t.Value})
	}
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.Render()

	return nil
}

// DeleteBucketTagging todo
func DeleteBucketTagging(c *cos.Client) error {
	_, err := c.Bucket.DeleteTagging(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// DeleteDesBucketTagging todo
func DeleteDesBucketTagging(c *cos.Client, tags []string) error {
	d, _, err := c.Bucket.GetTagging(context.Background())
	if err != nil {
		return err
	}
	table := make(map[string]string)
	for _, t := range d.TagSet {
		table[t.Key] = t.Value
	}
	for i := 0; i < len(tags); i += 1 {
		tmp := strings.Split(tags[i], "#")
		if len(tmp) != 2 || tmp[0] == "" || tmp[1] == "" {
			return fmt.Errorf("invalid tag format: %s", tags[i])
		}
		if _, exist := table[tmp[0]]; exist {
			delete(table, tmp[0])
		} else {
			return fmt.Errorf("the BucketTagging %s is not exist", tmp[0])
		}
	}
	tg := &cos.BucketPutTaggingOptions{}
	for a, b := range table {
		tg.TagSet = append(tg.TagSet, cos.BucketTaggingTag{Key: a, Value: b})
	}

	_, err = c.Bucket.PutTagging(context.Background(), tg)
	if err != nil {
		return err
	}
	return nil
}

// AddBucketTagging todo
func AddBucketTagging(c *cos.Client, tags []string) error {
	d, _, err := c.Bucket.GetTagging(context.Background())
	if err != nil {
		return err
	}
	table := make(map[string]string)
	for _, t := range d.TagSet {
		table[t.Key] = t.Value
	}
	for i := 0; i < len(tags); i += 1 {
		tmp := strings.Split(tags[i], "#")
		if len(tmp) != 2 || tmp[0] == "" || tmp[1] == "" {
			return fmt.Errorf("invalid tag format: %s", tags[i])
		}
		if _, exist := table[tmp[0]]; exist {
			return fmt.Errorf("the BucketTagging %s is already exist", tmp[0])
		}
		table[tmp[0]] = tmp[1]
	}
	tg := &cos.BucketPutTaggingOptions{}
	for a, b := range table {
		tg.TagSet = append(tg.TagSet, cos.BucketTaggingTag{Key: a, Value: b})
	}

	_, err = c.Bucket.PutTagging(context.Background(), tg)
	if err != nil {
		return err
	}
	return nil
}

// PutObjectTagging todo
func PutObjectTagging(c *cos.Client, object string, tags []string, versionId, bucketType string) error {
	var err error
	tg := &cos.ObjectPutTaggingOptions{}
	for i := 0; i < len(tags); i += 1 {
		tmp := strings.Split(tags[i], "#")
		if len(tmp) >= 2 {
			tg.TagSet = append(tg.TagSet, cos.ObjectTaggingTag{Key: tmp[0], Value: tmp[1]})
		} else {
			return fmt.Errorf("invalid tag")
		}
	}
	if bucketType == BucketTypeOfs {
		_, err = c.Object.PutTagging(context.Background(), object, tg)
	} else {
		_, err = c.Object.PutTagging(context.Background(), object, tg, versionId)
	}

	if err != nil {
		return err
	}

	return nil
}

// GetObjectTagging todo
func GetObjectTagging(c *cos.Client, object, versionId, bucketType string) error {
	var err error
	var v *cos.ObjectGetTaggingResult
	if bucketType == BucketTypeOfs {
		v, _, err = c.Object.GetTagging(context.Background(), object)
	} else {
		v, _, err = c.Object.GetTagging(context.Background(), object, versionId)
	}

	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	for _, t := range v.TagSet {
		table.Append([]string{t.Key, t.Value})
	}
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.Render()

	return nil
}

// DeleteObjectTagging todo
func DeleteObjectTagging(c *cos.Client, object, versionId, bucketType string) error {
	var err error
	if bucketType == BucketTypeOfs {
		_, err = c.Object.DeleteTagging(context.Background(), object)
	} else {
		_, err = c.Object.DeleteTagging(context.Background(), object, versionId)
	}

	if err != nil {
		return err
	}
	return nil
}

// DeleteDesObjectTagging todo
func DeleteDesObjectTagging(c *cos.Client, object string, tags []string, versionId, bucketType string) error {
	var err error
	var res *cos.ObjectGetTaggingResult
	if bucketType == BucketTypeOfs {
		res, _, err = c.Object.GetTagging(context.Background(), object)
	} else {
		res, _, err = c.Object.GetTagging(context.Background(), object, versionId)
	}

	if err != nil {
		return err
	}
	table := make(map[string]string)
	for _, t := range res.TagSet {
		table[t.Key] = t.Value
	}
	for i := 0; i < len(tags); i += 1 {
		tmp := strings.Split(tags[i], "#")
		if len(tmp) != 2 || tmp[0] == "" || tmp[1] == "" {
			return fmt.Errorf("invalid tag format: %s", tags[i])
		}
		if _, exist := table[tmp[0]]; exist {
			delete(table, tmp[0])
		} else {
			return fmt.Errorf("the ObjectTagging %s is not exist", tmp[0])
		}
	}
	tg := &cos.ObjectPutTaggingOptions{}
	for a, b := range table {
		tg.TagSet = append(tg.TagSet, cos.ObjectTaggingTag{Key: a, Value: b})
	}

	if bucketType == BucketTypeOfs {
		_, err = c.Object.PutTagging(context.Background(), object, tg)
	} else {
		_, err = c.Object.PutTagging(context.Background(), object, tg, versionId)
	}

	if err != nil {
		return err
	}
	return nil
}

// AddObjectTagging todo
func AddObjectTagging(c *cos.Client, object string, tags []string, versionId, bucketType string) error {
	var err error
	var res *cos.ObjectGetTaggingResult
	if bucketType == BucketTypeOfs {
		res, _, err = c.Object.GetTagging(context.Background(), object)
	} else {
		res, _, err = c.Object.GetTagging(context.Background(), object, versionId)
	}

	if err != nil {
		return err
	}
	table := make(map[string]string)
	for _, t := range res.TagSet {
		table[t.Key] = t.Value
	}
	for i := 0; i < len(tags); i += 1 {
		tmp := strings.Split(tags[i], "#")
		if len(tmp) != 2 || tmp[0] == "" || tmp[1] == "" {
			return fmt.Errorf("invalid tag format: %s", tags[i])
		}
		if _, exist := table[tmp[0]]; exist {
			return fmt.Errorf("the ObjectTagging %s is already exist", tmp[0])
		}
		table[tmp[0]] = tmp[1]
	}
	tg := &cos.ObjectPutTaggingOptions{}
	for a, b := range table {
		tg.TagSet = append(tg.TagSet, cos.ObjectTaggingTag{Key: a, Value: b})
	}
	if bucketType == BucketTypeOfs {
		_, err = c.Object.PutTagging(context.Background(), object, tg)
	} else {
		_, err = c.Object.PutTagging(context.Background(), object, tg, versionId)
	}

	if err != nil {
		return err
	}
	return nil
}

// EncodeTagging 对输入的标签字符串进行处理并编码
func EncodeTagging(tags string) (string, error) {
	if len(tags) == 0 {
		return "", nil
	}
	// 去掉字符串中的所有空格
	cleanTags := strings.ReplaceAll(tags, " ", "")

	// 将标签字符串按 & 拆分为多个键值对
	tagsList := strings.Split(cleanTags, "&")
	var encodedTags []string

	// 遍历每个键值对，分别对 Key 和 Value 进行 URL 编码
	for _, tag := range tagsList {
		// 按等号拆分成 Key 和 Value
		kv := strings.SplitN(tag, "=", 2) // 最多拆分成两部分
		if len(kv) == 2 {
			encodedKey := url.QueryEscape(kv[0])
			encodedValue := url.QueryEscape(kv[1])
			encodedTags = append(encodedTags, fmt.Sprintf("%s=%s", encodedKey, encodedValue))
		} else {
			return "", fmt.Errorf("error tags format:%s", tag)
		}
	}

	// 使用 & 符号连接多个编码后的键值对
	return strings.Join(encodedTags, "&"), nil
}
