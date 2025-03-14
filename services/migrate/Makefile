ROOT_PATH := ./../..
SOURCE ?= file go_bindata github github_ee bitbucket aws_s3 google_cloud_storage godoc_vfs gitlab
DATABASE ?= postgres mysql redshift cassandra spanner cockroachdb yugabytedb clickhouse mongodb sqlserver firebird neo4j pgx pgx5 rqlite
VERSION ?= $(shell git describe --tags 2>/dev/null | cut -c 2-)

build:
	CGO_ENABLED=0 go build -ldflags='-X main.Version=$(VERSION)' -tags '$(DATABASE) $(SOURCE)' $(ROOT_PATH)/internal/app/migrate/cmd/migrate

create:
	$(ROOT_PATH)/bin/migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_users_table