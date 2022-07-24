

fmt:
	@gofumpt -version || go install mvdan.cc/gofumpt@latest
	gofumpt -extra -w -d .
	@gci -v || go install github.com/daixiang0/gci@latest
	#gci write -s Std -s Def -s 'Prefix($(GO_MOD_DOMAIN))' -s 'Prefix($(GO_MOD_NAME))' .
	gci write -s standard -s default -s 'Prefix(github.com/costa92/go-web)' .

lint:
	golangci-lint version
	golangci-lint run -v --color always --out-format colored-line-number

