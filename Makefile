VERSION ?= 0.0.4
NAME ?= HONUA-DATABASE

release:
	go mod tidy
	git add .
	git commit -m "[RELEASE] ${NAME}: changes for v${VERSION}"
	git tag v${VERSION}
	git push origin v${VERSION}
	git push