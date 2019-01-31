mkdir -p coverage
PKG_LIST=$(go list ./... | grep -v /vendor/)
for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "coverage/${package##*/}.cov" "$package" ;
done
tail -q -n +2 coverage/*.cov >> coverage/coverage.cov