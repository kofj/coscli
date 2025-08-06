package cmd

import (
	"bytes"
	"context"
	"coscli/util"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// mockClient 创建模拟的 COS 客户端
func mockClient(t *testing.T) *cos.Client {
	return &cos.Client{}
}

// captureOutput 捕获标准输出
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// generateTestConfig 生成测试用清单配置
func generateTestConfig(id string) *cos.BucketGetInventoryResult {
	return &cos.BucketGetInventoryResult{
		ID:        id,
		IsEnabled: "true",
		Schedule: &cos.BucketInventorySchedule{
			Frequency: "Daily",
		},
		Destination: &cos.BucketInventoryDestination{
			Bucket:    "test-bucket",
			Format:    "CSV",
			AccountId: "123456789012",
			Prefix:    "inventory/",
		},
		IncludedObjectVersions: "All",
		Filter: &cos.BucketInventoryFilter{
			Prefix: "prefix/",
			Tags: []cos.ObjectTaggingTag{
				{Key: "env", Value: "production"},
			},
			StorageClass: "STANDARD",
			Period: &cos.BucketInventoryFilterPeriod{
				StartTime: time.Now().Unix(),
				EndTime:   time.Now().AddDate(0, 0, 7).Unix(),
			},
		},
		OptionalFields: &cos.BucketInventoryOptionalFields{
			BucketInventoryFields: []string{"Size", "LastModifiedDate"},
		},
	}
}

// generateTestListConfig 生成测试用清单配置
func generateTestListConfig(id string) *cos.BucketListInventoryConfiguartion {
	return &cos.BucketListInventoryConfiguartion{
		ID:        id,
		IsEnabled: "true",
		Schedule: &cos.BucketInventorySchedule{
			Frequency: "Daily",
		},
		Destination: &cos.BucketInventoryDestination{
			Bucket:    "test-bucket",
			Format:    "CSV",
			AccountId: "123456789012",
			Prefix:    "inventory/",
		},
		IncludedObjectVersions: "All",
		Filter: &cos.BucketInventoryFilter{
			Prefix: "prefix/",
			Tags: []cos.ObjectTaggingTag{
				{Key: "env", Value: "production"},
			},
			StorageClass: "STANDARD",
			Period: &cos.BucketInventoryFilterPeriod{
				StartTime: time.Now().Unix(),
				EndTime:   time.Now().AddDate(0, 0, 7).Unix(),
			},
		},
		OptionalFields: &cos.BucketInventoryOptionalFields{
			BucketInventoryFields: []string{"Size", "LastModifiedDate"},
		},
	}
}

func TestPutBucketInventory_Success(t *testing.T) {
	// 准备测试数据
	testID := "test-config"
	configContent := `{
		"Id": "test-config",
		"IsEnabled": "true",
		"Schedule": {"Frequency": "Daily"},
		"Destination": {
			"COSBucketDestination": {
				"Bucket": "test-bucket",
				"Format": "CSV"
			}
		}
	}`

	// 打桩
	patches := gomonkey.ApplyFunc(util.GetContent, func(input string) ([]byte, error) {
		return []byte(configContent), nil
	})
	defer patches.Reset()

	patches.ApplyFunc(util.ParseContent[cos.BucketPutInventoryOptions], func(content []byte, v interface{}) error {
		return json.Unmarshal(content, v)
	})

	patches.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "PutInventory", func(_ *cos.BucketService, ctx context.Context, id string, opt *cos.BucketPutInventoryOptions) (*cos.Response, error) {
		assert.Equal(t, testID, id)
		assert.Equal(t, "test-config", opt.ID)
		assert.Equal(t, "true", opt.IsEnabled)
		assert.Equal(t, "Daily", opt.Schedule.Frequency)
		return &cos.Response{}, nil
	})

	// 执行测试
	err := util.PutBucketInventory(mockClient(t), testID, configContent)

	// 验证结果
	assert.NoError(t, err)
}

func TestPutBucketInventory_InvalidConfig(t *testing.T) {
	// 准备测试数据
	testID := "test-config"
	invalidConfig := "invalid json"

	// 打桩
	patches := gomonkey.ApplyFunc(util.GetContent, func(input string) ([]byte, error) {
		return []byte(invalidConfig), nil
	})
	defer patches.Reset()

	patches.ApplyFunc(util.ParseContent[cos.BucketPutInventoryOptions], func(content []byte, v interface{}) error {
		return fmt.Errorf("unrecognized configuration format, must be JSON or XML")
	})

	// 执行测试
	err := util.PutBucketInventory(mockClient(t), testID, invalidConfig)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unrecognized configuration format, must be JSON or XML")
}

func TestPutBucketInventory_FileReadError(t *testing.T) {
	// 准备测试数据
	testID := "test-config"
	filePath := "file:///invalid/path.json"

	// 打桩
	patches := gomonkey.ApplyFunc(util.GetContent, func(input string) ([]byte, error) {
		return nil, fmt.Errorf("file not found")
	})
	defer patches.Reset()

	// 执行测试
	err := util.PutBucketInventory(mockClient(t), testID, filePath)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestPutBucketInventory_COSError(t *testing.T) {
	// 准备测试数据
	testID := "test-config"
	configContent := `{"Id": "test-config"}`

	// 打桩
	patches := gomonkey.ApplyFunc(util.GetContent, func(input string) ([]byte, error) {
		return []byte(configContent), nil
	})
	defer patches.Reset()

	patches.ApplyFunc(util.ParseContent[cos.BucketPutInventoryOptions], func(content []byte, v interface{}) error {
		return json.Unmarshal(content, v)
	})

	patches.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "PutInventory", func(_ *cos.BucketService, ctx context.Context, id string, opt *cos.BucketPutInventoryOptions) (*cos.Response, error) {
		return nil, fmt.Errorf("access denied")
	})

	// 执行测试
	err := util.PutBucketInventory(mockClient(t), testID, configContent)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestGetBucketInventory_Success(t *testing.T) {
	// 准备测试数据
	testID := "test-config"
	expectedConfig := generateTestConfig(testID)

	patches := gomonkey.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "GetInventory", func(_ *cos.BucketService, ctx context.Context, id string) (*cos.BucketGetInventoryResult, *cos.Response, error) {
		assert.Equal(t, testID, id)
		return expectedConfig, &cos.Response{}, nil
	})
	defer patches.Reset()

	// 捕获输出
	output := captureOutput(func() {
		err := util.GetBucketInventory(mockClient(t), testID)
		assert.NoError(t, err)
	})

	// 验证输出
	assert.Contains(t, output, "Inventory Configuration Details")
	assert.Contains(t, output, testID)
	assert.Contains(t, output, "true")
	assert.Contains(t, output, "Daily")
	assert.Contains(t, output, "test-bucket")
	assert.Contains(t, output, "prefix/")
	assert.Contains(t, output, "env=production")
	assert.Contains(t, output, "STANDARD")
	assert.Contains(t, output, "Size")
	assert.Contains(t, output, "LastModifiedDate")
}

func TestGetBucketInventory_NotFound(t *testing.T) {
	// 准备测试数据
	testID := "missing-config"

	// 打桩
	patches := gomonkey.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "GetInventory", func(_ *cos.BucketService, ctx context.Context, id string) (*cos.BucketGetInventoryResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("not found")
	})
	defer patches.Reset()

	// 捕获输出
	output := captureOutput(func() {
		err := util.GetBucketInventory(mockClient(t), testID)
		fmt.Println(err)
		assert.Error(t, err)
	})

	assert.Contains(t, output, "not found")
}

func TestDeleteBucketInventory_Success(t *testing.T) {
	// 准备测试数据
	testID := "test-config"

	// 打桩
	patches := gomonkey.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "DeleteInventory", func(_ *cos.BucketService, ctx context.Context, id string) (*cos.Response, error) {
		assert.Equal(t, testID, id)
		return &cos.Response{}, nil
	})
	defer patches.Reset()

	// 执行测试
	err := util.DeleteBucketInventory(mockClient(t), testID)

	// 验证结果
	assert.NoError(t, err)
}

func TestDeleteBucketInventory_Error(t *testing.T) {
	// 准备测试数据
	testID := "test-config"

	// 打桩
	patches := gomonkey.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "DeleteInventory", func(_ *cos.BucketService, ctx context.Context, id string) (*cos.Response, error) {
		return nil, fmt.Errorf("access denied")
	})
	defer patches.Reset()

	// 执行测试
	err := util.DeleteBucketInventory(mockClient(t), testID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestPostBucketInventory_Success(t *testing.T) {
	// 准备测试数据
	testID := "test-config"
	configContent := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><InventoryConfiguration xmlns=\"http://....\"><Id>test-config</Id><Destination><COSBucketDestination><Format>CSV</Format><AccountId>100000000002</AccountId><Bucket>qcs::cos:ap-nanjing::test-100000001</Bucket><Prefix>list4</Prefix><Encryption><SSE-COS>111</SSE-COS></Encryption></COSBucketDestination></Destination><Filter><And><Prefix>myPrefix</Prefix><Tag><Key>age</Key><Value>18</Value></Tag></And><Period><StartTime>1768688761</StartTime><EndTime>1568688762</EndTime></Period></Filter><IncludedObjectVersions>All</IncludedObjectVersions><OptionalFields><Field>Size</Field><Field>Tag</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>StorageClass</Field><Field>IsMultipartUploaded</Field></OptionalFields></InventoryConfiguration>"

	// 打桩
	patches := gomonkey.ApplyFunc(util.GetContent, func(input string) ([]byte, error) {
		return []byte(configContent), nil
	})
	defer patches.Reset()

	patches.ApplyFunc(util.ParseContent[cos.BucketPostInventoryOptions], func(content []byte, v interface{}) error {
		return json.Unmarshal(content, v)
	})

	patches.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "PostInventory", func(_ *cos.BucketService, ctx context.Context, id string, opt *cos.BucketPostInventoryOptions) (*cos.Response, error) {
		assert.Equal(t, testID, id)
		assert.Equal(t, "test-config", opt.ID)
		return &cos.Response{}, nil
	})

	// 执行测试
	err := util.PostBucketInventory(mockClient(t), testID, configContent)

	// 验证结果
	assert.NoError(t, err)
}

func TestPostBucketInventory_InvalidConfig(t *testing.T) {
	// 准备测试数据
	testID := "test-config"
	invalidConfig := "invalid json"

	// 打桩
	patches := gomonkey.ApplyFunc(util.GetContent, func(input string) ([]byte, error) {
		return []byte(invalidConfig), fmt.Errorf("invalid JSON")
	})
	defer patches.Reset()

	patches.ApplyFunc(util.ParseContent[cos.BucketPostInventoryOptions], func(content []byte, v interface{}) error {
		return fmt.Errorf("invalid JSON")
	})

	// 执行测试
	err := util.PostBucketInventory(mockClient(t), testID, invalidConfig)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestListBucketInventory_Success(t *testing.T) {
	// 准备测试数据
	config1 := generateTestListConfig("config1")
	config2 := generateTestListConfig("config2")

	// 打桩
	patches := gomonkey.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "ListInventoryConfigurations", func(_ *cos.BucketService, ctx context.Context, token string) (*cos.ListBucketInventoryConfigResult, *cos.Response, error) {
		if token == "" {
			return &cos.ListBucketInventoryConfigResult{
				InventoryConfigurations: []cos.BucketListInventoryConfiguartion{
					*config1,
				},
				IsTruncated:           true,
				NextContinuationToken: "next",
			}, &cos.Response{}, nil
		}
		return &cos.ListBucketInventoryConfigResult{
			InventoryConfigurations: []cos.BucketListInventoryConfiguartion{
				*config2,
			},
			IsTruncated: false,
		}, &cos.Response{}, nil
	})

	defer patches.Reset()

	// 捕获输出
	output := captureOutput(func() {
		err := util.ListBucketInventory(mockClient(t))
		assert.NoError(t, err)
	})

	// 验证输出
	assert.Contains(t, output, "Detailed COS Bucket Inventory Configurations")
	assert.Contains(t, output, "config1")
	assert.Contains(t, output, "config2")
	assert.Contains(t, output, "Daily")
	assert.Contains(t, output, "test-bucket")
	assert.Contains(t, output, "Total inventory configurations: 2")
}

func TestListBucketInventory_Empty(t *testing.T) {
	// 打桩
	patches := gomonkey.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "ListInventoryConfigurations", func(_ *cos.BucketService, ctx context.Context, token string) (*cos.ListBucketInventoryConfigResult, *cos.Response, error) {
		return &cos.ListBucketInventoryConfigResult{
			InventoryConfigurations: []cos.BucketListInventoryConfiguartion{},
			IsTruncated:             false,
		}, &cos.Response{}, nil
	})
	defer patches.Reset()

	// 捕获输出
	output := captureOutput(func() {
		err := util.ListBucketInventory(mockClient(t))
		assert.NoError(t, err)
	})

	// 验证输出
	assert.Contains(t, output, "Total inventory configurations: 0")
}

func TestListBucketInventory_Error(t *testing.T) {
	// 打桩
	patches := gomonkey.ApplyMethod(reflect.TypeOf((*cos.BucketService)(nil)), "ListInventoryConfigurations", func(_ *cos.BucketService, ctx context.Context, token string) (*cos.ListBucketInventoryConfigResult, *cos.Response, error) {
		return nil, nil, fmt.Errorf("access denied")
	})
	defer patches.Reset()

	// 执行测试
	err := util.ListBucketInventory(mockClient(t))

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestRenderInventoryConfigDetail(t *testing.T) {
	// 准备测试数据
	config := generateTestConfig("test-config")

	// 捕获输出
	output := captureOutput(func() {
		util.RenderInventoryConfigDetail(config)
	})

	// 验证输出
	assert.Contains(t, output, "Inventory Configuration Details")
	assert.Contains(t, output, "test-config")
	assert.Contains(t, output, "true")
	assert.Contains(t, output, "Daily")
	assert.Contains(t, output, "test-bucket")
	assert.Contains(t, output, "prefix/")
	assert.Contains(t, output, "env=production")
	assert.Contains(t, output, "STANDARD")
	assert.Contains(t, output, "Size")
	assert.Contains(t, output, "LastModifiedDate")
}

func TestRenderInventoryConfigDetail_EmptyConfig(t *testing.T) {
	// 准备测试数据
	config := &cos.BucketGetInventoryResult{
		ID: "empty-config",
	}

	// 捕获输出
	output := captureOutput(func() {
		util.RenderInventoryConfigDetail(config)
	})

	// 验证输出
	assert.Contains(t, output, "Inventory Configuration Details")
	assert.Contains(t, output, "empty-config")
	assert.Contains(t, output, "Not configured")
}
