if [[ -z $TIMETRACE_DIR ]]; then
	cp -r autocompletion/ $HOME/.timetrace/
    TIMETRACE_DIR=$HOME/.timetrace
    cat >> $HOME/.bashrc << EOF

# >>> timetrace initialization >>>
# Configuration of timetrace for autocompletion
export TIMETRACE_DIR=\$HOME/.timetrace

if [ -f \$TIMETRACE_DIR/autocompletion/bash/timetrace.sh ]; then
    . \$TIMETRACE_DIR/autocompletion/bash/timetrace.sh
fi

# <<< timetrace initialization <<<

EOF
fi
