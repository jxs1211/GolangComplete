call chain:

    cmd/plugin/devlake/main.go:
        ->internal/pkg/plugin/devlake/read.go:getDynamicState
            ->internal/pkg/plugin/devlake/state.go:getDynamicState
                ->internal/pkg/plugininstaller/common/repo.go:GetPluginAllK8sState

    internal/pkg/plugininstaller/helm/state.go:GetPluginAllState
    ->internal/pkg/plugininstaller/common/repo.go:GetPluginAllK8sState

    internal/pkg/plugin/argocdapp/state.go:getStaticState
    ->internal/pkg/plugin/argocdapp/state.go:getArgoCDAppFromK8sAndSetState

    internal/pkg/plugininstaller/reposcaffolding/installer.go:DeleteRepo
    ->internal/pkg/plugininstaller/common/repo.go:CreateGithubClient
    ->internal/pkg/plugininstaller/common/repo.go:CreateGitlabClient

    internal/pkg/plugininstaller/reposcaffolding/option.go:PushToRemoteGitlab
        ->internal/pkg/plugininstaller/common/repo.go:CreateGitlabClient
        ->internal/pkg/plugininstaller/common/repo.go:BuildgitlabOpts

    internal/pkg/plugininstaller/reposcaffolding/option.go:PushToRemoteGithub
        ->internal/pkg/plugininstaller/common/repo.go:CreateGithubClient
    
    internal/pkg/plugininstaller/reposcaffolding/state.go:getGithubStatus
        ->internal/pkg/plugininstaller/common/repo.go:CreateGithubClient
    
    internal/pkg/plugininstaller/reposcaffolding/state.go:getGitlabStatus
        ->internal/pkg/plugininstaller/common/repo.go:CreateGitlabClient
    
    internal/pkg/plugininstaller/reposcaffolding/utils.go:walkLocalRepoPath
        ->internal/pkg/plugininstaller/common/repo.go:CreateLocalRepoPath
        ->internal/pkg/plugininstaller/common/repo.go:GenerateRenderWalker
            ->internal/pkg/plugininstaller/common/repo.go:replaceAppNameInPathStr
            ->pkg/util/file/file.go:CopyFile
    
    internal/pkg/show/config/gen_embed_var.go:copyTemplates
        ->pkg/util/file/file.go:CopyFile

    pkg/util/k8s/state.go:GetResourceStatus
        ->pkg/util/k8s/workload.go:ListDeploymentsWithLabel
            ->pkg/util/k8s/workload.go:generateLabelFilterOption
        ->pkg/util/k8s/workload.go:ListStatefulsetsWithLabel
            ->pkg/util/k8s/workload.go:generateLabelFilterOption
        ->pkg/util/k8s/workload.go:ListDaemonsetsWithLabel
            ->pkg/util/k8s/workload.go:generateLabelFilterOption
