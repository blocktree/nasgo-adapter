module github.com/blocktree/nasgo-adapter

go 1.12

require (
	github.com/astaxie/beego v1.12.0
	github.com/blocktree/go-owcdrivers v1.2.0
	github.com/blocktree/go-owcrypt v1.1.1
	github.com/blocktree/openwallet/v2 v2.0.6
	github.com/go-errors/errors v1.0.1
	github.com/imroc/req v0.2.4
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	github.com/tidwall/gjson v1.3.5
	gopkg.in/resty.v1 v1.12.0
)

//replace golang.org/x/net v0.0.0-20181220203305-927f97764cc3 => github.com/golang/net v0.0.0-20181220203305-927f97764cc3
