TAG := 1.21.11
REPO := rithwik0604/minecraft-server

buildimage:
	docker build -t $(REPO):$(TAG) image/
pushimage:
	docker push $(REPO):$(TAG)
build:
	export GOOS=linux; export GOARCH=amd64; go build .

# clean everything after build
postbuildclean:
	rm -rf .git .gitignore Dockerfile go.* index.html main.go makefile