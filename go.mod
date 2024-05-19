module github.com/FUMCGarland/hvac

go 1.22.0

require (
	github.com/briandowns/openweathermap v0.19.0
	github.com/brutella/dnssd v1.2.10
	github.com/eclipse/paho.golang v0.21.0
	github.com/go-co-op/gocron/v2 v2.5.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lestrrat-go/jwx/v2 v2.0.21
	github.com/mochi-mqtt/server/v2 v2.6.3
	github.com/warthog618/go-gpiocdev v0.9.0
	golang.org/x/crypto v0.23.0
	golang.org/x/time v0.5.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.5 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/miekg/dns v1.1.59 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/lestrrat-go/jwx/v2 v2.0.21 => ./deps/jwx
