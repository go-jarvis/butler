# butler

1. 根据配置 `config{}` 生成对应的 `default.yml` 配置文件。 
2. 读取依次配置文件 `default.yml, config.yml` + `分支配置文件.yml` + `环境变量`
    + 根据 GitlabCI, 分支配置文件 `config.xxxx.yml`
    + 如没有 CI, 读取本地文件: `local.yml`

