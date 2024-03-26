# Saving Data

To save your plugin data like plugin settings, configuration and statistics, we are going to use the [ConfigApi.Plugin.Save](../api/config-api.md#plugin) method:

```go
type MyPluginConfig struct {
    MySetting       string  `json:"my_setting"`
    OtherSetting    int     `json:"other_setting"`
}

my_key := "my_config"
my_config := MyPluginConfig{
    MySetting:      "my_value",
    OtherSetting:   123,
}

err := api.Config().Plugin(config_key).Save(my_config)
```

Plugin configuration is separated into different keys for ease of management.

To get your plugin data for a specific key, use the [ConfigApi.Plugin.Get](../api/config-api.md#plugin) method:
```go

my_config, err := api.Config().Plugin(my_key).Get()
```

