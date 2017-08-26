#!/bin/bash
_wrench() 
{
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    #
    #  The basic options we'll complete.
    #
    opts="build bump config help push run version -h --help"


    #
    #  Complete the arguments to some of the basic commands.
    #
    case "${prev}" in
        build)
            local build_opts="-h --help -r --rebuild"
            COMPREPLY=($(compgen -W "${build_opts}" -- "${cur}"))
            return 0
            ;;
        bump)
            local bump_opts="major minor patch -h --help"
            COMPREPLY=($(compgen -W "${bump_opts}" -- "${cur}"))
            return 0
            ;;
        config)
            local config_opts="--format -h --help"
            COMPREPLY=($(compgen -W "${config_opts}" -- "${cur}"))
            return 0
            ;;
        push)
            local push_opts="--additional-tags -h --help"
            COMPREPLY=($(compgen -W "${push_opts}" -- "${cur}"))
            return 0
            ;;
        run)
            local wrench_run_config wrench_run_targets
            local regex='[0-9A-Za-z-]+:'
            wrench_run_config=$(wrench config --format '{{.Run}}')
            wrench_run_targets=$(echo "$wrench_run_config" | grep -oE "$regex" | cut -d':' -f1)
            COMPREPLY=($(compgen -W "${wrench_run_targets}" -- "${cur}"))
            return 0
            ;;
        *)
        ;;
    esac

   COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))  
   return 0
}
complete -o bashdefault -o default -F _wrench wrench
