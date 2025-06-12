# Contributing简中简化版本

首先，克隆repo

```shell
git clone https://github.com/NodeSeekDev/alist.git
git clone --recurse-submodules https://github.com/friedHDD/alist-web.git
```
安装依赖并且初次编译前端部分

> 不要尝试使用npm i啊喂

```shell
cd alist-webui
npm i pnpm -g
pnpm i
bash release.sh
```

等pnpm编译完了

`./alist-web/dist/`全部丢到`./public/dist/`

在alist-web的目录里面`pnpm start`启动！在根目录里面`go run main.go server --dev`启动！

(初次使用先`go run main.go admin set <password> 设置密码`)

如果需要离线编译请参考

```shell
go mod vendor
go build -mod=vendor
```

大概机制就是你得给后端吃点东西，但是你实际上开发的前端是在vite上运行的那个，也就是5173端口通过proxy把api搞到5244端口去了，而5244运行的是完整的alist