# show help message
@default: help

App := 'peekbt'
Version := `grep '^const VERSION = ' cmd/main/version.go | sed "s/^VERSION = \"\(.*\)\"/\1/g"`

# show help message
@help:
    echo "Build tool for {{ App }} {{ Version }} with Just"
    echo "Usage: just <recipe>"
    echo ""
    just --list

# build the application with running tests
build: test
    go build -o peekbt cmd/main/peekbt.go

# run tests and generate the coverage report
test:
    go test -covermode=count -coverprofile=coverage.out ./...

# clean up build artifacts
clean:
    go clean
    rm -f coverage.out peekbt

# update the version if the new version is provided
update_version new_version = "":
    if [ "{{ new_version }}" != "" ]; then \
        sed 's/$VERSION/{{ new_version }}/g' .template/README.md > README.md; \
        sed 's/$VERSION/{{ new_version }}/g' .template/version.go > cmd/main/version.go; \
    fi

# build wildcherry for all platforms
make_distribution_files:
    for os in "linux" "windows" "darwin"; do \
        for arch in "amd64" "arm64"; do \
            mkdir -p dist/wildcherry-$os-$arch; \
            env GOOS=$os GOARCH=$arch go build -o dist/wildcherry-$os-$arch/wildcherry cmd/main/wildcherry.go; \
            cp README.md LICENSE dist/wildcherry-$os-$arch; \
            tar cvfz dist/wildcherry-$os-$arch.tar.gz -C dist wildcherry-$os-$arch; \
        done; \
    done

upload_assets tag:
    gh release upload --repo tamada/wildcherry {{ tag }} dist/*.tar.gz