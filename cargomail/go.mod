module cargomail

go 1.20

require (
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.17
	golang.org/x/crypto v0.10.0
	golang.org/x/sync v0.3.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/bmizerany/pat v0.0.0-20170815010413-6226ea591a40 // indirect
	github.com/tus/tusd/v2 v2.0.0-rc9.0.20230703112035-529171da612a
)

replace github.com/tus/tusd/v2 v2.0.0-rc9.0.20230703112035-529171da612a => github.com/cargomail-org/tusd/v2 v2.0.0-rc9.0.20230716212531-a33df8f348f1
