## How to use
下载release运行 `xttDns.exe` 

## 开机启动
将`xttDns.exe`放到以下目录<br/>
`C:/Users/%username%/AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup`

## 运行配置
```yaml
# host配置文件地址
host_remote_address: http://host-server.com

# 临时文件
temp_file_path: hosts

# 重启进程路径
restart_exe_path: C:\Program Files\Internet Explorer\iexplore.exe

# 重启进程名称
restart_process_name: iexplore.exe

# 更新周期 秒
refresh_second: 10

# 运行日志
log_file: c:/xttdns/xtt.log
```