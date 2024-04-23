package service

import (
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	tenantCmd = &cobra.Command{
		Use:               "tenant command <args> [flags]",
		Short:             "租户相关操作命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
可用的配置类型：
[list,add,del,add-user,del-user,add-models,del-models]

aigc-server tenant -h
`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err = prepare(cmd.Context()); err != nil {
				return errors.Wrap(err, "prepare")
			}
			return nil
		},
	}
)

var (
	tenantAddCmd = &cobra.Command{
		Use:               "add <args> [flags]",
		Short:             "添加租户",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：

aigc-server tenant add <name> <email>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) != 2 {
				return errors.New("参数错误")
			}
			tenant := &types.Tenants{
				Name:           args[0],
				ContactEmail:   args[1],
				PublicTenantID: uuid.New().String(),
			}
			if err = gormDB.WithContext(ctx).Create(tenant).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "create tenant error", "err", err)
				return errors.Wrap(err, "create tenant error")
			}
			return nil
		},
	}

	tenantDelCmd = &cobra.Command{
		Use:               "del <args> [flags]",
		Short:             "删除租户",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server tenant del <uuid>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("参数错误")
			}
			if err = gormDB.Where("public_tenant_id = ?", args[0]).Delete(&types.Tenants{}).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "delete tenant error", "err", err)
				return errors.Wrap(err, "delete tenant error")
			}
			return nil
		},
	}

	tenantListCmd = &cobra.Command{
		Use:               "list [flags]",
		Short:             "租户列表",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server tenant list 
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var tenants []types.Tenants
			if err = gormDB.WithContext(cmd.Context()).Order("id DESC").Find(&tenants).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "list tenant error", "err", err)
				return errors.Wrap(err, "list tenant error")
			}
			var data [][]string
			for _, tenant := range tenants {
				data = append(data, []string{tenant.PublicTenantID, tenant.Name, tenant.ContactEmail})
			}
			renderTable([]string{"租户ID", "租户名称", "联系人邮箱"}, data)
			return nil
		},
	}

	tenantAddUserCmd = &cobra.Command{
		Use:               "add-user <args> [flags]",
		Short:             "添加租户用户",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server tenant add-user <tenant_id> <email>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("参数错误")
			}
			tenant := &types.Tenants{}
			if err = gormDB.Where("public_tenant_id = ?", args[0]).First(tenant).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "find tenant error", "err", err)
				return errors.Wrap(err, "find tenant error")
			}
			account := &types.Accounts{}
			if err = gormDB.Where("email = ?", args[1]).First(account).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "find account error", "err", err)
				return errors.Wrap(err, "find account error")
			}
			if err = gormDB.Model(tenant).Association("Accounts").Append(account); err != nil {
				_ = level.Warn(logger).Log("msg", "add account to tenant error", "err", err)
				return errors.Wrap(err, "add account to tenant error")
			}
			return nil
		},
	}

	tenantDelUserCmd = &cobra.Command{
		Use:               "del-user <args> [flags]",
		Short:             "删除租户用户",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server tenant del-user <tenant_id> <email>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("参数错误")
			}
			tenant := &types.Tenants{}
			if err = gormDB.Where("public_tenant_id = ?", args[0]).First(tenant).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "find tenant error", "err", err)
				return errors.Wrap(err, "find tenant error")
			}
			account := &types.Accounts{}
			if err = gormDB.Where("email = ?", args[1]).First(account).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "find account error", "err", err)
				return errors.Wrap(err, "find account error")
			}
			if err = gormDB.Model(tenant).Association("Accounts").Delete(account); err != nil {
				_ = level.Warn(logger).Log("msg", "delete account from tenant error", "err", err)
				return errors.Wrap(err, "delete account from tenant error")
			}
			return nil
		},
	}

	tenantAddModelsCmd = &cobra.Command{
		Use:               "add-models <args> [flags]",
		Short:             "添加租户模型",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server tenant add-models <tenant_id> <model_name1> <model_name2> ...
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("参数错误")
			}
			tenant := &types.Tenants{}
			if err = gormDB.Where("public_tenant_id = ?", args[0]).First(tenant).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "find tenant error", "err", err)
				return errors.Wrap(err, "find tenant error")
			}
			for _, modelName := range args[1:] {
				model := &types.Models{}
				if err = gormDB.Where("name = ?", modelName).First(model).Error; err != nil {
					_ = level.Warn(logger).Log("msg", "find model error", "err", err)
					return errors.Wrap(err, "find model error")
				}
				if err = gormDB.Model(tenant).Association("Models").Append(model); err != nil {
					_ = level.Warn(logger).Log("msg", "add model to tenant error", "err", err)
					return errors.Wrap(err, "add model to tenant error")
				}
			}
			return nil
		},
	}

	tenantDelModelsCmd = &cobra.Command{
		Use:               "del-models <args> [flags]",
		Short:             "删除租户模型",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
Example：
aigc-server tenant del-models <tenant_id> <model_name1> <model_name2> ...
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("参数错误")
			}
			tenant := &types.Tenants{}
			if err = gormDB.Where("public_tenant_id = ?", args[0]).First(tenant).Error; err != nil {
				_ = level.Warn(logger).Log("msg", "find tenant error", "err", err)
				return errors.Wrap(err, "find tenant error")
			}
			for _, modelName := range args[1:] {
				model := &types.Models{}
				if err = gormDB.Where("name = ?", modelName).First(model).Error; err != nil {
					_ = level.Warn(logger).Log("msg", "find model error", "err", err)
					return errors.Wrap(err, "find model error")
				}
				if err = gormDB.Model(tenant).Association("Models").Delete(model); err != nil {
					_ = level.Warn(logger).Log("msg", "delete model from tenant error", "err", err)
					return errors.Wrap(err, "delete model from tenant error")
				}
			}
			return nil
		},
	}
)

func init() {
	tenantCmd.AddCommand(tenantAddCmd, tenantDelCmd, tenantListCmd, tenantAddUserCmd, tenantDelUserCmd, tenantAddModelsCmd, tenantDelModelsCmd)
}
