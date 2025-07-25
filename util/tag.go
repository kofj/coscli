package util

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tencentyun/cos-go-sdk-v5"
	"os"
	"strings"
)

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

func DeleteBucketTagging(c *cos.Client) error {
	_, err := c.Bucket.DeleteTagging(context.Background())
	if err != nil {
		return err
	}
	return nil
}

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

func PutObjectTagging(c *cos.Client, object string, tags []string) error {
	tg := &cos.ObjectPutTaggingOptions{}
	for i := 0; i < len(tags); i += 1 {
		tmp := strings.Split(tags[i], "#")
		if len(tmp) >= 2 {
			tg.TagSet = append(tg.TagSet, cos.ObjectTaggingTag{Key: tmp[0], Value: tmp[1]})
		} else {
			return fmt.Errorf("invalid tag")
		}
	}

	_, err := c.Object.PutTagging(context.Background(), object, tg)
	if err != nil {
		return err
	}

	return nil
}

func GetObjectTagging(c *cos.Client, object string) error {
	v, _, err := c.Object.GetTagging(context.Background(), object)
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

func DeleteObjectTagging(c *cos.Client, object string) error {
	_, err := c.Object.DeleteTagging(context.Background(), object)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDesObjectTagging(c *cos.Client, object string, tags []string) error {
	d, _, err := c.Object.GetTagging(context.Background(), object)
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
			return fmt.Errorf("the ObjectTagging %s is not exist", tmp[0])
		}
	}
	tg := &cos.ObjectPutTaggingOptions{}
	for a, b := range table {
		tg.TagSet = append(tg.TagSet, cos.ObjectTaggingTag{Key: a, Value: b})
	}

	_, err = c.Object.PutTagging(context.Background(), object, tg)
	if err != nil {
		return err
	}
	return nil
}

func AddObjectTagging(c *cos.Client, object string, tags []string) error {
	d, _, err := c.Object.GetTagging(context.Background(), object)
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
			return fmt.Errorf("the ObjectTagging %s is already exist", tmp[0])
		}
		table[tmp[0]] = tmp[1]
	}
	tg := &cos.ObjectPutTaggingOptions{}
	for a, b := range table {
		tg.TagSet = append(tg.TagSet, cos.ObjectTaggingTag{Key: a, Value: b})
	}

	_, err = c.Object.PutTagging(context.Background(), object, tg)
	if err != nil {
		return err
	}
	return nil
}
