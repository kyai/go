#!/bin/bash
echo "KTV Golang updater is running..."

base_path="./"
pid_file=$base_path"run.pid"
run_file=$base_path"message_server"

# svn update
if [ "$1"x == "u"x ]; then
    svn up
fi

# get and kill pid
if [ -f "$pid_file" ]; then
    pid=`cat "$pid_file"`
    # echo $pid
    kill $pid
    echo "kill $pid successed"
fi

# restart server
if [ -f "$run_file" ]; then
    nohup $run_file &
fi

echo "Done"
exit 0
