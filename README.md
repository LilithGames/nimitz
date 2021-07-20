## 一个用于快速发开Kubernetes动态准入控制的框架

### 一、项目结构 
- [log] 日志组件提供日志处理的实例
- [config] 读取yaml初始化配置信息,所有配置读取只通过此处。镜像匹配的规则添加在`ImageRenameRules`中 
- [webhook/mutating] 动态准入控制Webhook的业务逻辑，如果需要Validating操作在此处添加 
- [deploy] docker和k8s的配置文件
- [cmd/webhook ] 程序入口，在此添加需要运行的webhook
- [test] 用于测试webhook是否生效的nginx deployment文件

### 二、运行说明
- 由于webhook运行需要依赖证书，因此无法直接在本地运行，可拉取镜像进行测试
- `MakeFile`中命令执行顺序为编译程序，构建镜像，部署。
- 构建镜像会把`bin`目录下编译生成的webhook二进制文件拷贝在根目录下，部署时候会自动运行`./webhook`

### 三、使用说明 
- `ImagePodMutator` 接收config.yaml中的`ImageRules`提供的pattern和replace，值得注意的是现阶段只提供正则表达式的字符串替换方式
- `Example`:  
-  `pattern: "^(.*)lilith-registry(.*)$"`  
- `replace: "\\{1}lilith-registry-vpc\\{2}"`  
  提供的image的值为`lilith-registry.cn-shanghai.cr.aliyuncs.com/avatar/code:latest`  
  最终的被更新后的值为`lilith-registry-vpc.cn-shanghai.cr.aliyuncs.com/avatar/code:latest`