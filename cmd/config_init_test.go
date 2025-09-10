package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"os"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInit(t *testing.T) {
	fmt.Println("TestConfigInit")
	//Convey("success", t, func() {
	//	cmd := rootCmd
	//	cmd.SilenceErrors = true
	//	cmd.SilenceUsage = true
	//	args := []string{"config", "init"}
	//	cmd.SetArgs(args)
	//	e := cmd.Execute()
	//	fmt.Printf(" : %v", e)
	//	So(e, ShouldBeNil)
	//})
	Convey("fail", t, func() {
		patches := ApplyFunc(initConfigFile, func(cfgFlag bool) error {
			return fmt.Errorf("test initConfigFile error")
		})
		defer patches.Reset()
		clearCmd()
		cmd := rootCmd
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		args := []string{"config", "init"}
		cmd.SetArgs(args)
		e := cmd.Execute()
		fmt.Printf(" : %v", e)
		So(e, ShouldBeError)
	})
}

//
//func TestInitConfigFile(t *testing.T) {
//	// 模拟用户输入
//	configPath := "testinit.yaml"
//	restoreInput := mockUserInput(t, []string{
//		configPath, // 使用默认路径
//		"SecretKey", // 模式
//		"test-secret-id", // Secret ID
//		"test-secret-key", // Secret Key
//		"", // Session Token (空)
//		"true", // DisableEncryption
//		"false", // DisableAutoFetchBucketType
//		"false", // CloseAutoSwitchHost
//		"test-bucket-1234567890", // Bucket Name
//		"cos.ap-beijing.myqcloud.com", // Endpoint
//		"test-alias", // Bucket Alias
//	})
//	defer restoreInput()
//
//	// 执行函数
//	err := initConfigFile(true)
//	require.NoError(t, err)
//
//	// 验证配置文件路径
//	assert.FileExists(t, configPath)
//
//	// 读取配置文件
//	v := readConfigFile(t, configPath)
//
//	// 验证配置内容
//	assert.Equal(t, "SecretKey", v.GetString("cos.base.mode"))
//	assert.Equal(t, "test-secret-id", v.GetString("cos.base.secretID"))
//	assert.Equal(t, "test-secret-key", v.GetString("cos.base.secretKey"))
//	assert.Equal(t, "", v.GetString("cos.base.sessionToken"))
//	assert.Equal(t, true, v.GetBool("cos.base.disableEncryption"))
//	assert.Equal(t, false, v.GetBool("cos.base.disableAutoFetchBucketType"))
//	assert.Equal(t, false, v.GetBool("cos.base.closeAutoSwitchHost"))
//	assert.Equal(t, "https", v.GetString("cos.base.protocol"))
//	delYaml("testinit.yaml")
//}

// mockUserInput 模拟用户输入
func mockUserInput(t *testing.T, inputs []string) func() {
	t.Helper()

	// 创建输入缓冲区
	inputBuffer := bytes.NewBufferString("")
	for _, input := range inputs {
		inputBuffer.WriteString(input + "\n")
	}

	// 保存原始标准输入
	oldStdin := os.Stdin

	// 替换标准输入
	r, w, _ := os.Pipe()
	os.Stdin = r

	// 写入模拟输入
	go func() {
		defer w.Close()
		io.Copy(w, inputBuffer)
	}()

	// 返回恢复函数
	return func() {
		os.Stdin = oldStdin
	}
}

// readConfigFile 读取配置文件内容
func readConfigFile(t *testing.T, path string) *viper.Viper {
	t.Helper()

	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	return v
}
