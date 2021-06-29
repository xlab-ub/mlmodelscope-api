cleanup_and_exit () {
    cd - 2>&1 >> /dev/null
    exit  0
}

die_with_message () {
    echo "$1"
    cleanup_and_exit 1
}
