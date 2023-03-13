module github.com/cargomail-org/grip/poc/smtpd

go 1.20

require github.com/cargomail-org/smtpd-grip v0.0.0-20230313102159-db2477f0ab9d

replace (
	github.com/cargomail-org/smtpd-grip => ../../../smtpd-grip
)
