package service

import (
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var (
	accountTenantId string // 租户ID
	limit           = 10
	offset          = 0

	accountCmd = &cobra.Command{
		Use:               "account command <args> [flags]",
		Short:             "用户相关操作命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
可用的配置类型：
[list,add,del,password]

aigc-server account -h
`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err = prepare(cmd.Context()); err != nil {
				return errors.Wrap(err, "prepare")
			}
			return nil
		},
	}
)

var (
	accountAddCmd = &cobra.Command{
		Use:               "add <args> [flags]",
		Short:             "添加用户",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：

aigc-server account add <email> <password>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) != 2 {
				return errors.New("参数错误")
			}
			var tenant types.Tenants
			var tenantErr error
			if accountTenantId != "" {
				tenant, tenantErr = store.Tenants().FindTenantByTenantId(ctx, accountTenantId)
			} else {
				tenant, tenantErr = store.Tenants().FindTenant(ctx, 1)
			}
			if tenantErr != nil {
				_ = level.Warn(logger).Log("msg", "find tenant error", "err", tenantErr)
			}
			password, err := bcrypt.GenerateFromPassword([]byte(args[1]), bcrypt.DefaultCost)
			if err != nil {
				_ = level.Warn(logger).Log("msg", "generate password error", "err", err)
				return err
			}
			if gormDB.Model(&types.Accounts{}).Where("email = ?", args[0]).First(&types.Accounts{}).Error == nil {
				_ = level.Warn(logger).Log("msg", "account already exists", "email", args[0])
				return errors.New("account already exists")
			}
			if err = gormDB.Save(&types.Accounts{
				Email:        args[0],
				Nickname:     "新用户",
				Language:     "zh",
				IsLdap:       false,
				PasswordHash: string(password),
				Status:       true,
				Tenants:      []types.Tenants{tenant},
			}).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "save account error", "err", err)
				return errors.Wrap(err, "save account error")
			}
			return nil
		},
	}

	accountPasswordCmd = &cobra.Command{
		Use:               "password <args> [flags]",
		Short:             "修改用户密码",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server account password <email> <password>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("参数错误")
			}
			password, err := bcrypt.GenerateFromPassword([]byte(args[1]), bcrypt.DefaultCost)
			if err != nil {
				_ = level.Warn(logger).Log("msg", "generate password error", "err", err)
				return err
			}
			if err = gormDB.Model(&types.Accounts{}).Where("email = ?", args[0]).Update("password_hash", string(password)).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "update password error", "err", err)
				return errors.Wrap(err, "update password error")
			}
			return nil
		},
	}

	accountDeleteCmd = &cobra.Command{
		Use:               "del <args> [flags]",
		Short:             "删除用户",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server account del <email>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("参数错误")
			}
			if err = gormDB.Where("email = ?", args[0]).Delete(&types.Accounts{}).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "delete account error", "err", err)
				return errors.Wrap(err, "delete account error")
			}
			return nil
		},
	}

	accountListCmd = &cobra.Command{
		Use:               "list [flags]",
		Short:             "用户列表",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server account list --limit 10 --offset 0
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var accounts []types.Accounts
			if err = gormDB.WithContext(cmd.Context()).Model(&types.Accounts{}).
				Limit(limit).Offset(offset).
				Order("id DESC").Find(&accounts).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "find accounts error", "err", err)
				return errors.Wrap(err, "find accounts error")
			}
			var data [][]string
			for _, account := range accounts {
				data = append(data, []string{account.Email, account.Nickname})
			}
			renderTable([]string{"邮箱", "名称"}, data)
			return nil
		},
	}
)

func init() {
	accountAddCmd.PersistentFlags().StringVar(&accountTenantId, "tenant.id", "", "租户ID")
	accountListCmd.PersistentFlags().IntVar(&limit, "limit", 10, "limit")
	accountListCmd.PersistentFlags().IntVar(&offset, "offset", 0, "offset")

	accountCmd.AddCommand(accountAddCmd, accountPasswordCmd, accountDeleteCmd, accountListCmd)
}
