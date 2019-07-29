package config

import "github.com/spf13/viper"

var (
	// 环境变量前缀
	ENVPREFIX = "INJECTOR"
	// bash命令的绝对路径
	BASHCMD = "BASH_CMD"
	// 需要注入的namespace的标签选择器
	NAMESPACELABELSELECTORS = "NAMESPACE_LABEL_SELECTORS"
	// 监控的configmap所在的namespace
	CONFIGMAPNAMESPACE = "COMFIG_MAP_NAMESPACE"
)

var (
	BashCmd string
	NamespaceLabelSelectors string
	ConfigMapNamespace string
)

func init() {
	viper.SetEnvPrefix(ENVPREFIX)
	bindEnv(BASHCMD, "/usr/local/bin/bash", &BashCmd)
	bindEnv(NAMESPACELABELSELECTORS, "", &NamespaceLabelSelectors)
	bindEnv(CONFIGMAPNAMESPACE, "kube-system", &ConfigMapNamespace)
}

func bindEnv(name, defaultValue string, v *string) {
	if err := viper.BindEnv(name); err != nil {
		panic(err)
	}
	viper.SetDefault(name, defaultValue)
	*v = viper.Get(name).(string)
}
