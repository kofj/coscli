package cmd

import (
	"coscli/util"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigSetCmd(t *testing.T) {
	fmt.Println("TestConfigSetCmd")
	getConfig()
	var oldconfig util.Config = config
	secretKey, err := util.DecryptSecret(config.Base.SecretKey)
	if err == nil {
		oldconfig.Base.SecretKey = secretKey
	}
	secretId, err := util.DecryptSecret(config.Base.SecretID)
	if err == nil {
		oldconfig.Base.SecretID = secretId
	}
	sessionToken, err := util.DecryptSecret(config.Base.SessionToken)
	if err == nil {
		oldconfig.Base.SessionToken = sessionToken
	}
	copyYaml()
	defer restoreYaml()
	clearCmd()
	cmd := rootCmd
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	Convey("Test coscil config set", t, func() {
		Convey("fail", func() {
			Convey("no arguments", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"config", "set"}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
			Convey("no mode", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"config", "set", "--mode", "@"}
				cmd.SetArgs(args)
				e := cmd.Execute()
				fmt.Printf(" : %v", e)
				So(e, ShouldBeError)
			})
		})
		Convey("success", func() {
			Convey("give arguments", func() {
				clearCmd()
				cmd := rootCmd
				args := []string{"config", "set", "--secret_id", oldconfig.Base.SecretID,
					"--secret_key", oldconfig.Base.SecretKey, "--session_token", "@", "--mode", oldconfig.Base.Mode,
					"--cvm_role_name", "@", "--close_auto_switch_host", oldconfig.Base.CloseAutoSwitchHost, "--disable_encryption", oldconfig.Base.DisableEncryption}
				cmd.SetArgs(args)
				e := cmd.Execute()
				So(e, ShouldBeNil)
			})
		})

	})
}
