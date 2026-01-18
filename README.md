# sakibox

如果你也相信，命令行的光芒可以抚平琐碎的疲惫，那么欢迎来到 sakibox。它是一只温柔的终端工具箱，陪你整理日常、守护秩序，也在必要时为你指引方向。

## 功能

- 端口管理：查看端口、查找端口占用、关闭端口进程
- 历史命令：查看、搜索并执行历史命令
- 命令收藏夹：保存常用命令并执行、删除
- 进程监控：实时进程列表、资源占用 TOP10、搜索/杀死进程
- 文件查找：按名称、扩展名、内容、大小、修改时间、全局检索
- 安装帮助：生成 Linux 工具/依赖安装命令

## 使用方式

Go 安装（需要 Go 1.24.0）:

macOS:

```bash
brew install go
```

Linux:

Ubuntu/Debian:

```bash
sudo apt update
sudo apt install -y golang
```

CentOS/RHEL:

```bash
sudo yum install -y golang
```

Arch:

```bash
sudo pacman -S go
```

Alpine:

```bash
sudo apk add go
```

通用安装（wget + 官方包）:

```bash
cd /tmp
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
```

一条命令安装:

```bash
curl -fsSL https://raw.githubusercontent.com/usukiy128/sakibox/main/install.sh | sh
```

默认安装到 `/usr/local/bin`。如需自定义路径:

```bash
SAKIBOX_INSTALL_DIR="$HOME/.local/bin" curl -fsSL https://raw.githubusercontent.com/usukiy128/sakibox/main/install.sh | sh
```

运行:

```bash
sakibox
```

更新:

```bash
sakibox
```

进入主菜单后选择“更新 sakibox”即可。

历史命令会根据当前终端尝试读取对应的历史文件（如 `~/.zsh_history` 或 `~/.bash_history`），也可以在 `~/.sakibox/config.yaml` 中配置 `history_file` 自定义路径。

## 目录结构

- cmd: CLI 入口与菜单
- internal: 功能实现
- config: 配置读取
- data: 数据文件

## 许可

Apache License 2.0
