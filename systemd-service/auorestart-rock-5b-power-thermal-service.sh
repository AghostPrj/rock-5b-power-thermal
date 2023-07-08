#!/bin/bash

# 检查"rock-5b-power-thermal"进程的数量
count=$(ps -aux | grep rock-5b-power-thermal | grep -v grep | wc -l)

# 如果数量为0，则重启服务
if [ "$count" -eq 0 ]; then
    systemctl restart rock-5b-power-thermal
fi