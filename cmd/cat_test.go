package cmd

import (
	"coscli/util"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCatCmd(t *testing.T) {
	fmt.Println("TestCatCmd")
	testBucket = randStr(8)
	testAlias = testBucket + "-alias"
	setUp(testBucket, testAlias, testEndpoint, false, false)
	defer tearDown(testBucket, testAlias, testEndpoint, false)
	clearCmd()
	cmd := rootCmd
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	genDir(testDir, 3)
	defer delDir(testDir)

	Convey("test coscli cat", t, func() {
		Convey("上传测试单文件", func() {
			clearCmd()
			cmd := rootCmd
			localFileName := fmt.Sprintf("%s/small-file/0", testDir)
			cosFileName := fmt.Sprintf("cos://%s/%s", testAlias, "single-small")
			args := []string{"cp", localFileName, cosFileName}
			cmd.SetArgs(args)
			e := cmd.Execute()
			So(e, ShouldBeNil)
		})
		Convey("fail", func() {
			Convey("Not enough argument", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"cat"}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("FormatUrl", func() {
				patches := ApplyFunc(util.FormatUrl, func(urlStr string) (util.StorageUrl, error) {
					return nil, fmt.Errorf("test FormatUrl error")
				})
				defer patches.Reset()
				clearCmd()
				cmd := rootCmd
				args := []string{"cat", fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("Not CosUrl", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"cat", testAlias}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("NewClient", func() {
				patches := ApplyFunc(util.NewClient, func(config *util.Config, param *util.Param, bucketName string) (client *cos.Client, err error) {
					return nil, fmt.Errorf("test NewClient error")
				})
				defer patches.Reset()
				clearCmd()
				cmd := rootCmd
				args := []string{"cat", fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("CatObject", func() {
				patches := ApplyFunc(util.CatObject, func(c *cos.Client, cosUrl util.StorageUrl) error {
					return fmt.Errorf("test CatObject error")
				})
				defer patches.Reset()
				clearCmd()
				cmd := rootCmd
				args := []string{"cat", fmt.Sprintf("cos://%s", testAlias)}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
		})
		Convey("success", func() {
			clearCmd()
			cmd := rootCmd
			args := []string{"cat", fmt.Sprintf("cos://%s/%s", testAlias, "single-small")}
			cmd.SetArgs(args)
			e := cmd.Execute()
			So(e, ShouldBeNil)
		})
		Convey("retry", func() {
			Convey("retry-2xx", func() {
				Convey("retry-2xx-CloseAutoSwitchHost", func() {
					Convey("retry-200r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/200r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-200", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/200", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-204r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/204r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-204", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/204", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-206r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/206r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-206", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/206", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
				})
				Convey("retry-2xx-AutoSwitchHost", func() {
					Convey("retry-200r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/200r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-200", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/200", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-204r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/204r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-204", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/204", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-206r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/206r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
					Convey("retry-206", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/206", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeNil)
					})
				})
			})
			Convey("retry-3xx", func() {
				Convey("retry-3xx-CloseAutoSwitchHost", func() {
					Convey("retry-301r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/301r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
					Convey("retry-301", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/301", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
					Convey("retry-/302r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/302r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
					Convey("retry-302", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/302", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
				})
				Convey("retry-3xx-AutoSwitchHost", func() {
					Convey("retry-301r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/301r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
					Convey("retry-301", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/301", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
					Convey("retry-/302r", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/302r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
					Convey("retry-302", func() {
						clearCmd()
						cmd := rootCmd
						args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/302", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
						cmd.SetArgs(args)
						e := cmd.Execute()
						fmt.Printf(" : %v", e)
						So(e, ShouldBeError)
					})
				})
			})
			//Convey("retry-4xx", func() {
			//	Convey("retry-4xx-CloseAutoSwitchHost", func() {
			//		Convey("retry-400r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/400r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-400", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/400", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-403r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/403r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-403", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/403", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-404r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/404r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-404", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/404", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//	})
			//	Convey("retry-4xx-AutoSwitchHost", func() {
			//		Convey("retry-400r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/400r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-400", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/400", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-403r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/403r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-403", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/403", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-404r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/404r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-404", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/404", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//	})
			//})
			//Convey("retry-5xx", func() {
			//	Convey("retry-5xx-CloseAutoSwitchHost", func() {
			//		Convey("retry-500r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/500r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-500", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/500", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-503r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/503r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-503", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/503", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-504r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/504r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-timeout", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/timeout", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-shutdown", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/shutdown", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "true", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//	})
			//	Convey("retry-5xx-AutoSwitchHost", func() {
			//		Convey("retry-500r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/500r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-500", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/500", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-503r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/503r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-503", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/503", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-504r", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/504r", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-504", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/504", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//		Convey("retry-timeout", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/timeout", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeError)
			//		})
			//		Convey("retry-shutdown", func() {
			//			clearCmd()
			//			cmd := rootCmd
			//			args := []string{"cat", "cos://cos-sdk-err-retry-1253960454/shutdown", "-e", "cos.ap-chengdu.myqcloud.com", "--close_auto_switch_host", "false", "-p", "http"}
			//			cmd.SetArgs(args)
			//			e := cmd.Execute()
			//			fmt.Printf(" : %v", e)
			//			So(e, ShouldBeNil)
			//		})
			//	})
			//})
		})
	})
}
