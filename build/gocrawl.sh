echo "Copy environment file to env"
yes | cp -rf build/gocrawler-env /root/go/env/gocrawler-env
echo "Build go application service"
GOOS=linux GOARCH=amd64 go build -o gocrawler main.go
echo "Restart service"
systemctl restart gocrawler
