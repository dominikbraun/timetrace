# file: timetrace
# timetrace parameter autocompletion


_complete_timetrace () {
    # regexp that match the first arguments
    local cmd="${1##*/}"
    local word=${COMP_WORDS[COMP_CWORD]}
    local xpat
    local commands="create delete edit get list help start status stop version"
    local no_completion="help start status stop version"
    local flags="-h -v"

    #create, delete, edit 
    local common_subcommand="project"
    #get
    local get_subcommands="project record"
    #list resources
    local list_subcommands="projects records"

    # Check to see what command is been executed
    cur_command=$(_find_command)
    case "$cmd" in
    timetrace)
        if [[ "$(_get_argc)" -eq 1 ]]; then
            xpat="$commands"
        elif [[ "$(_get_argc)" -ge 2 ]]; then
            case "$cur_command" in
            # create, delete and edit subcommands
            create | delete | edit)
                xpat="$common_subcommand"
                ;;

            # get subcommands
            get)
                xpat="$get_subcommands"
                ;;

            # list subcommands
            list)
                xpat="$list_subcommands"
                ;;

            # commands that support only flags
            help | start | status | stop | version)
                xpat="$flags"
                ;;

            *)
                if [[ "$(_get_argc)" -eq 2 ]]; then
                    xpat="$commands"
                fi
                ;;
            esac
        fi
        ;;

    esac

    COMPREPLY=($(compgen -W "$xpat" -- "${word}"))
}

_find_command () {
    local line=${COMP_LINE}
    IFS=' ' read -r -a array <<< "$line"
    len="${#array[@]}"
    
    # check if there's no subcommand
    if [[ $len -ge 2 ]]; then
        echo "${array[1]}"
    fi
}

_get_argc () {
    local line=${COMP_LINE}
    IFS=' ' read -r -a array <<< "$line"
    len="${#array[@]}"
    echo "$len"
}

complete -F _complete_timetrace timetrace
