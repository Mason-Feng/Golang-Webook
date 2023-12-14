.PHONY: docker
docker:
	#把上次编译的版本删掉
	@rm webook || true
	#运行一下 go mod tidy，防止go.sum文件不对，编译失败
	@go mod tidy
	#指定编译在ARM架构的linux操作系统上运行的可执行文件
	#后面为名字
	@GOOS=linux GOARCH=arm  go build -tags=k8s -o webook .
	#可以每次部署前增加删除之前的版本，从而简化操作
	@docker rmi -f fqw/webook:v0.0.1
	#在这里更改标签来描述相应的版本，同时对应的k8s部署里面也要修改
	@docker  build -t fqw/webook:v0.0.1 .