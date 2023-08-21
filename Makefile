VERSION ?= 1.0.13
NAME ?= HONUA-DATABASE

release:
	go mod tidy
	git add .
	git commit -m "[RELEASE] ${NAME}: changes for v${VERSION}"
	git tag v${VERSION}
	git push origin v${VERSION}
	git push