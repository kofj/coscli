package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
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
func TestInitConfigFile(t *testing.T) {
	// 1. 保存原始标准输入
	originalStdin := os.Stdin
	defer func() {
		os.Stdin = originalStdin
	}()

	// 2. 保存viper全局状态
	originalViper := *viper.GetViper()
	defer func() {
		*viper.GetViper() = originalViper
	}()

	// 3. 保存全局变量状态
	originalCmdCnt := cmdCnt
	defer func() {
		cmdCnt = originalCmdCnt
	}()

	// 4. 重置viper状态，确保测试环境干净
	viper.Reset()

	// 5. 使用唯一的临时文件路径，避免文件名冲突
	tempDir := t.TempDir() // 使用testing包提供的临时目录
	configPath := filepath.Join(tempDir, "testinit.yaml")
	// 模拟用户输入
	//configPath := "testinit.yaml"
	restoreInput := mockUserInput(t, []string{
		configPath,                    // 使用默认路径
		"SecretKey",                   // 模式
		"test-secret-id",              // Secret ID
		"test-secret-key",             // Secret Key
		"",                            // Session Token (空)
		"true",                        // DisableEncryption
		"false",                       // DisableAutoFetchBucketType
		"false",                       // CloseAutoSwitchHost
		"test-bucket-1234567890",      // Bucket Name
		"cos.ap-beijing.myqcloud.com", // Endpoint
		"test-alias",                  // Bucket Alias
	})
	defer restoreInput()

	// 执行函数
	err := initConfigFile(true)
	require.NoError(t, err)

	// 验证配置文件路径
	assert.FileExists(t, configPath)

	// 读取配置文件
	v := readConfigFile(t, configPath)

	// 验证配置内容
	assert.Equal(t, "SecretKey", v.GetString("cos.base.mode"))
	assert.Equal(t, "test-secret-id", v.GetString("cos.base.secretID"))
	assert.Equal(t, "test-secret-key", v.GetString("cos.base.secretKey"))
	assert.Equal(t, "", v.GetString("cos.base.sessionToken"))
	assert.Equal(t, true, v.GetBool("cos.base.disableEncryption"))
	assert.Equal(t, false, v.GetBool("cos.base.disableAutoFetchBucketType"))
	assert.Equal(t, false, v.GetBool("cos.base.closeAutoSwitchHost"))
	assert.Equal(t, "https", v.GetString("cos.base.protocol"))
	//delYaml("testinit.yaml")
}

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
