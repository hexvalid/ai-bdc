# ai-bdc
monorepo for ai-bdc

## Generating Data

1. Deploy a Digitalocean's CPU optimized Debian 9 droplet
2. Convert dropted to Arch Linux by using this script: https://github.com/gh2o/digitalocean-debian-to-arch
3. Install dependencies: 
```bash
pacman -S screen bash-completion go php php-cgi php-gd libjpeg libpng fontconfig
git clone https://github.com/hexvalid/ai-bdc
go get github.com/cheggaaa/pb/v3
cd ai-bdc/gen
mkdir out
```

4. Prepare PHP-CGI:
```bash
# Open up 10 php-cgi server via screen
php -c php/php.ini -t php/ -S 127.0.0.1:9000
php -c php/php.ini -t php/ -S 127.0.0.1:9001
php -c php/php.ini -t php/ -S 127.0.0.1:9002
......
php -c php/php.ini -t php/ -S 127.0.0.1:9009
```

5. Run Generation Script:
```bash
# Check mode and count
# Run mode = 1 and mode = 2
# As twice!
nano gen.go

# run via screen
go run gen.go
```

6. Package:
```bash
tar cfJ gendata-250k-white.tar.xz out/*
```

7. Clean:
```bash
rm /tmp/bdc_pipe && touch /tmp/bdc_pipe && rm /tmp/bdc_void/*
```
