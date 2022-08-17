首先介绍下devstream的核心组件

The architecture image below shows how DevStream works in general:

![img](https://docs.devstream.io/en/latest/images/architecture-overview.png)

- config manager:负责配置管理
- plugin manager：
- state manager：
- backend manager：负责存储插件的配置和状态

接下来，以插件create操作的流程为例，我们一起来解读下源码

先看下apply操作的日志文件，这个是分析源码的主线

```sh
./dtm apply -f config-harbor.yaml -y
2022-08-11 10:03:25 ℹ [INFO]  Apply started.
2022-08-11 10:03:25 ℹ [INFO]  Got Backend from config: local
2022-08-11 10:03:25 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-08-11 10:03:25 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-08-11 10:03:25 ℹ [INFO]  Tool (harbor-docker/default) found in config but does not exist in the state, will be created.
2022-08-11 10:03:25 ℹ [INFO]  Start executing the plan.
2022-08-11 10:03:25 ℹ [INFO]  Changes count: 1.
2022-08-11 10:03:25 ℹ [INFO]  -------------------- [  Processing progress: 1/1.  ] --------------------
2022-08-11 10:03:25 ℹ [INFO]  Processing: (harbor-docker/default) -> Create ...
2022-08-11 10:03:26 ℹ [INFO]  Cmd: ./install.sh.
......
Stdout: ✔ ----Harbor has been installed and started successfully.----
2022-08-11 10:03:31 ℹ [INFO]  Cmd: docker compose ls.
Stdout: NAME                STATUS              CONFIG FILES
Stdout: devstream           running(9)          /root/devstream/docker-compose.yml
2022-08-11 10:03:31 ✔ [SUCCESS]  Tool (harbor-docker/default) Create done.
2022-08-11 10:03:31 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-08-11 10:03:31 ✔ [SUCCESS]  All plugins applied successfully.
2022-08-11 10:03:31 ✔ [SUCCESS]  Apply finished.
```

然后我们根据打印来一步步跟踪源码：

首先apply会执行pluginengine.Apply

```go
// internal/pkg/pluginengine/cmd_apply.go
func applyCMDFunc(cmd *cobra.Command, args []string) {
	log.Info("Apply started.")
	if err := pluginengine.Apply(configFile, continueDirectly); err != nil {
		log.Errorf("Apply failed => %s.", err)
		os.Exit(1)
	}
	log.Success("Apply finished.")
}
```

检验成功返回true
```go
// internal/pkg/configmanager/coreconfig.go
func (c *CoreConfig) Validate() (bool, error) {
	...

}
```

configmanager登场，确定插件的存储后端，并tools文件的md5是否匹配

```go
// internal/pkg/pluginmanager/manager.go
func CheckLocalPlugins(conf *configmanager.Config) error {
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory doesn't exist")
	}

	log.Infof("Using dir <%s> to store plugins.", pluginDir)

	for _, tool := range conf.Tools {
		pluginFileName := configmanager.GetPluginFileName(&tool)
		pluginMD5FileName := configmanager.GetPluginMD5FileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("plugin %s doesn't exist", tool.Name)
			}
			return err
		}
		if _, err := os.Stat(filepath.Join(pluginDir, pluginMD5FileName)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf(".md5 file of plugin %s doesn't exist", tool.Name)
			}
			return err
		}
		if err := pluginAndMD5Matches(pluginDir, pluginFileName, pluginMD5FileName, tool.Name); err != nil {
			return err
		}
	}
	return nil
}
```

statemanager登场，确定存储方式，并读取状态配置，将其存入到statesMap

```go
// internal/pkg/statemanager/manager.go
// NewManager returns a new Manager and reads states through backend defined in config.
func NewManager(stateConfig configmanager.State) (Manager, error) {
	if m != nil {
		return m, nil
	}

	log.Debugf("The global manager m is not initialized.")

	// Get the backend from config
	if stateConfig.Backend == "local" {
		log.Infof("Using local backend. State file: %s.", stateConfig.Options.StateFile)
	} else if stateConfig.Backend == "s3" {
		log.Infof("Using s3 backend. Bucket: %s, region: %s, key: %s.", stateConfig.Options.Bucket, stateConfig.Options.Region, stateConfig.Options.Key)
	}

	b, err := backend.GetBackend(stateConfig)
	if err != nil {
		log.Errorf("Failed to get the Backend: %s.", err)
		return nil, err
	}

	m = &manager{
		Backend:   b,
		statesMap: NewStatesMap(),
	}

	// Read the initial states data from backend
	data, err := b.Read()
	if err != nil {
		log.Debugf("Failed to read data from backend: %s.", err)
		return nil, err
	}

	tmpMap := make(map[StateKey]State)
	if err = yaml.Unmarshal(data, tmpMap); err != nil {
		log.Errorf("Failed to unmarshal the state file < %s >. error: %s.", local.DefaultStateFile, err)
		log.Errorf("Reading the state file failed, it might have been compromised/modified by someone other than DTM.")
		log.Errorf("The state file is managed by DTM automatically. Please do not modify it yourself.")
		return nil, fmt.Errorf("state format error")
	}
	for k, v := range tmpMap {
		log.Debugf("Got a state from the backend: %s -> %v.", k, v)
		m.statesMap.Store(k, v)
	}

	return m, nil
}
```

接下来也是重头戏，使用之前创建好的statemanager和configmanager，对每一个插件进行排序，然后遍历每一个插件，对比配置和状态的差异，这个过程俗称调协(reconcile):
通过调协决定需要执行的操作类型为create或update，并将其记录到一个change的切片以备后用

最后将剩余的plugin删除

```go
func changesForApply(smgr statemanager.Manager, cfg *configmanager.Config) ([]*Change, error) {
	changes := make([]*Change, 0)

	// 1. create a temporary state map used to store unprocessed tools.
	tmpStatesMap := smgr.GetStatesMap().DeepCopy()

	// 2. handle dependency and sort the tools in the config into "batches" of tools
	var batchesOfTools [][]configmanager.Tool
	// the elements in batchesOfTools are sorted "batches"
	// and each element/batch is a list of tools that, in theory, can run in parallel
	// that is to say, the tools in the same batch won't depend on each other
	batchesOfTools, err := topologicalSort(cfg.Tools)
	if err != nil {
		return changes, err
	}

	// 3. generate changes for each tool
	for _, batch := range batchesOfTools {
		for _, tool := range batch {
			state := smgr.GetState(statemanager.StateKeyGenerateFunc(&tool))

			if state == nil {
				// tool not in the state, create, no need to Read resource before Create
				description := fmt.Sprintf("Tool (%s/%s) found in config but doesn't exist in the state, will be created.", tool.Name, tool.InstanceID)
				changes = append(changes, generateCreateAction(&tool, description))
			} else {
				// tool found in the state

				// first, handle possible "outputs" references in the tool's config
				// ignoring errors, since at this stage we are calculating changes, and the dependency might not have its output in the state yet
				_ = HandleOutputsReferences(smgr, tool.Options)

				if drifted(state.Options, tool.Options) {
					// tool's config differs from State's, Update
					description := fmt.Sprintf("Tool (%s/%s) config drifted from the state, will be updated.", tool.Name, tool.InstanceID)
					changes = append(changes, generateUpdateAction(&tool, description))
				} else {
					// tool's config is the same as State's

					// read resource status
					resource, err := Read(&tool)
					if err != nil {
						return changes, err
					}

					if resource == nil {
						// tool exists in the state, but resource doesn't exist, Create
						description := fmt.Sprintf("Tool (%s/%s) state found but it seems the tool isn't created, will be created.", tool.Name, tool.InstanceID)
						changes = append(changes, generateCreateAction(&tool, description))
					} else if drifted, err := ResourceDrifted(state.Resource, resource); drifted || err != nil {
						if err != nil {
							return nil, err
						}
						// resource drifted from state, Update
						description := fmt.Sprintf("Tool (%s/%s) drifted from the state, will be updated.", tool.Name, tool.InstanceID)
						changes = append(changes, generateUpdateAction(&tool, description))
					} else {
						// resource is the same as the state, do nothing
						log.Debugf("Tool (%s/%s) is the same as the state, do nothing.", tool.Name, tool.InstanceID)
					}
				}
			}

			// delete the tool from the temporary state map since it's already been processed above
			tmpStatesMap.Delete(statemanager.StateKeyGenerateFunc(&tool))
		}
	}

	// what's left in the temporary state map "tmpStatesMap" contains tools that:
	// - have a state (probably created previously)
	// - don't have a definition in the config (probably deleted by the user)
	// thus, we need to generate a "delete" change for it.
	tmpStatesMap.Range(func(key, value interface{}) bool {
		changes = append(changes, generateDeleteActionFromState(value.(statemanager.State)))
		log.Infof("Change added: %s -> %s", key.(statemanager.StateKey), statemanager.ActionDelete)
		return true
	})

	return changes, nil
}
```

前戏差不多，开始干正事了，这里开始就是create操作的具体实现逻辑

根据常见的目录，加载提前build好的插件

```go
 // internal/pkg/pluginengine/plugin.go
func Create(tool *configmanager.Tool) (map[string]interface{}, error) {
	pluginDir := getPluginDir()
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return nil, err
	}
	return p.Create(tool.Options)
}
```
通过配置文件的配置，创建相应插件，这里要创建的插件是harbordocker，加载并并执行其create方法。
这一步很关键，其定义了整个安装过程中，需要执行的操作，包括：
- 预执行操作：解析配置文件
- 执行操作：安装插件
- 状态查询操作：查询插件的状态

```go
// internal/pkg/plugin/harbordocker/create.go
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			renderConfig,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			Install,
		},
		GetStateOperation: dockerInstaller.ComposeState,
	}

	// Execute all Operations in Operator
	state, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", state)
	return state, nil
}
```

我们来依次看一下这几步操作：

首先，render配置，主要是要将var中设置的变量渲染到模板文件里，从而得到一个完整的配置文件
```go
// internal/pkg/plugin/harbordocker/harbordocker.go
func renderConfig(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts := Options{}
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	// TODO(daniel-hutao): use template wrapper here
	tmpl, err := template.New("compose").Delims("[[", "]]").Parse(HarborConfigTemplate)
	if err != nil {
		return nil, err
	}

	configFile, err := os.Create(HarborConfigFileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := configFile.Close(); err != nil {
			log.Errorf("Failed to close opened file (%s): %s.", configFile.Name(), err)
		}
	}()

	if err := tmpl.Execute(configFile, opts); err != nil {
		return nil, err
	}

	return options, err
}
```

接下来就是安装插件，主要有2步：
- writeScripts：加载脚本
```sh
func writeScripts() error {
	for name, sh := range scripts {
		err := os.WriteFile(name, []byte(sh), 0744)
		if err != nil {
			return err
		}
	}
	return nil
}
```
- ExecInSystemWithParams：执行脚本
```go
func ExecInSystemWithParams(execPath string, params []string, logsBuffer *bytes.Buffer, print bool) error {
	paramStr := strings.Join(params, " ")
	return ExecInSystem(execPath, paramStr, logsBuffer, print)
}
```

然后就是等待安装，安装成功后会有如下打印：
```sh
Stdout: ✔ ----Harbor has been installed and started successfully.----
```

安装完成后，检查插件状态
```go
func ComposeState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	state, err := op.ComposeState()
	if err != nil {
		return nil, fmt.Errorf("failed to get containers state: %s", err)
	}

	return state, nil
}
```

执行完操作后，打印结果信息并更新到statemanager。
```go
func handleResult(smgr statemanager.Manager, change *Change) error {
	log.Debugf("Start -> StatesMap now is:\n%s", string(smgr.GetStatesMap().Format()))
	defer func() {
		log.Debugf("End -> StatesMap now is:\n%s", string(smgr.GetStatesMap().Format()))
	}()

	if !change.Result.Succeeded {
		// do nothing when the change failed
		return nil
	}

	if change.ActionName == statemanager.ActionDelete {
		key := statemanager.StateKeyGenerateFunc(change.Tool)
		log.Infof("Prepare to delete '%s' from States.", key)
		err := smgr.DeleteState(key)
		if err != nil {
			log.Debugf("Failed to delete state %s: %s.", key, err)
			return err
		}
		log.Successf("Tool (%s/%s) delete done.", change.Tool.Name, change.Tool.InstanceID)
		return nil
	}

	key := statemanager.StateKeyGenerateFunc(change.Tool)
	state := statemanager.State{
		Name:       change.Tool.Name,
		InstanceID: change.Tool.InstanceID,
		DependsOn:  change.Tool.DependsOn,
		Options:    change.Tool.Options,
		Resource:   change.Result.ReturnValue,
	}
	err := smgr.AddState(key, state)
	if err != nil {
		log.Debugf("Failed to add state %s: %s.", key, err)
		return err
	}
	log.Successf("Tool (%s/%s) %s done.", change.Tool.Name, change.Tool.InstanceID, change.ActionName)
	return nil
}
```

自此，整个插件的安装过程完成。
