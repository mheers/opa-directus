module github.com/mheers/opa-directus

go 1.24.1

require (
	github.com/altipla-consulting/directus-go/v2 v2.2.0
	github.com/coder/websocket v1.8.13
	github.com/open-policy-agent/opa v1.3.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	golang.org/x/sys v0.31.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace github.com/altipla-consulting/directus-go/v2 => github.com/mheers/directus-go/v2 v2.0.0-20250416142358-8b3394c5a9e6
