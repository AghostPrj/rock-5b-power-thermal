# rock-5b-power-thermal systemd 服务使用手册 / rock-5b-power-thermal Systemd Service Manual

## 文件位置以及所需操作 / File Locations and Required Operations

首先创建文件夹`/var/lib/rock-5b-power-thermal `  
First, create the directory `/var/lib/rock-5b-power-thermal`

```text
${path_to_bin_file} -> /var/lib/rock-5b-power-thermal/rock-5b-power-thermal
rock-5b-power-thermal.yaml -> /var/lib/rock-5b-power-thermal/rock-5b-power-thermal.yaml
auorestart-rock-5b-power-thermal-service.sh -> /usr/local/bin/auorestart-rock-5b-power-thermal-service.sh
rock-5b-power-thermal.service -> /etc/systemd/system/rock-5b-power-thermal.service
rock-5b-power-thermal-autorestart.service -> /etc/systemd/system/rock-5b-power-thermal-autorestart.service
rock-5b-power-thermal-autorestart.timer -> /etc/systemd/system/rock-5b-power-thermal-autorestart.timer
```

接下来，根据你的需求编辑 `/var/lib/rock-5b-power-thermal/rock-5b-power-thermal.yaml` 文件。  
Then edit `/var/lib/rock-5b-power-thermal/rock-5b-power-thermal.yaml` as your mind.

编辑配置文件后，以root用户身份运行以下命令来启动并启用服务。  
After edit config file, run these command on root to start and enable service.

```bash
chmod a+x /usr/local/bin/auorestart-rock-5b-power-thermal-service.sh
systemctl daemon-reload 
systemctl enable rock-5b-power-thermal-autorestart.timer rock-5b-power-thermal.service
systemctl restart rock-5b-power-thermal-autorestart.timer rock-5b-power-thermal.service
```





