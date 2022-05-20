kill $(cat traymenu_pid) >/dev/null 2>&1;

nohup ./traymenu > traymenu.log 2>&1 &
echo $! > traymenu_pid
