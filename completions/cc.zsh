#compdef cc

_cc() {
    local -a commands
    commands=(
        'install:Deploy plugin commands and agents to a target project'
        'update:Update deployed plugin to latest embedded version'
        'status:Show deployment status and detect drift'
        'list:List available plugin assets'
        'uninstall:Remove deployed plugin commands and agents'
        'version:Print version'
        'completion:Generate shell completion scripts'
    )

    local -a global_flags
    global_flags=(
        '--help[Show help]'
    )

    _arguments -C \
        '1:command:->command' \
        '*::arg:->args'

    case $state in
    command)
        _describe -t commands 'cc commands' commands
        ;;
    args)
        case $words[1] in
        install)
            local -a install_args
            install_args=(
                'commands:Install only commands'
                'agents:Install only agents'
                'all:Install everything'
            )
            _arguments \
                '1:category:->category' \
                '--scope[Deployment scope]:scope:(user project)' \
                '--target[Target project root]:directory:_directories' \
                '--dry-run[Show what would be installed]' \
                '--plugin[Plugin to install]:plugin:(research-plan-implement-validate)' \
                '--help[Show help]'
            case $state in
            category)
                _describe -t install_args 'install category' install_args
                ;;
            esac
            ;;
        update)
            _arguments \
                '--scope[Deployment scope]:scope:(user project)' \
                '--target[Target project root]:directory:_directories' \
                '--diff[Show changes before applying]' \
                '--force[Overwrite locally modified files]' \
                '--help[Show help]'
            ;;
        status)
            _arguments \
                '--scope[Deployment scope]:scope:(user project)' \
                '--target[Target project root]:directory:_directories' \
                '--help[Show help]'
            ;;
        list)
            local -a list_args
            list_args=(
                'commands:List available commands'
                'agents:List available agents'
                'plugins:List available plugins'
            )
            _arguments \
                '1:category:->category' \
                '--help[Show help]'
            case $state in
            category)
                _describe -t list_args 'list category' list_args
                ;;
            esac
            ;;
        uninstall)
            _arguments \
                '--scope[Deployment scope]:scope:(user project)' \
                '--target[Target project root]:directory:_directories' \
                '--help[Show help]'
            ;;
        completion)
            local -a shells
            shells=(
                'bash:Generate bash completions'
                'zsh:Generate zsh completions'
            )
            _arguments '1:shell:->shell'
            case $state in
            shell)
                _describe -t shells 'shell' shells
                ;;
            esac
            ;;
        version)
            _arguments '--help[Show help]'
            ;;
        esac
        ;;
    esac
}

_cc "$@"
