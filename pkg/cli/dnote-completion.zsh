#compdef dnote

_dnote() {
  local -a _1st_arguments

  _1st_arguments=(
    'add:add a new note'
    'view:list books, notes, or view a content'
    'edit:edit a note or a book'
    'remove:remove a note or a book'
    'find:find notes by keywords'
    'sync:sync data with the server'
    'login:login to the dnote server'
    'logout:logout from the dnote server'
    'version:print the current version'
    'help:get help about any command'
  )

  _get_booknames() {
    local names
    names=$(dnote view --name-only 2>/dev/null)
    local -a ret

    while IFS= read -r line; do
      [[ -n "$line" ]] && ret+=("${line}")
    done <<< "$names"

    echo "${ret[@]}"
  }

  if (( CURRENT == 2 )); then
    _describe -t commands "dnote subcommand" _1st_arguments
    return
  elif (( CURRENT == 3 )); then
    case "$words[2]" in
      v|view|a|add)
        _alternative \
          "names:book names:($(_get_booknames))"
        ;;
    esac
  fi
}

compdef _dnote dnote 2>/dev/null || true

if [ "$funcstack[1]" = "_dnote" ]; then
  _dnote "$@"
fi
