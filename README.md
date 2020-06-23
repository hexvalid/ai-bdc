# Training

## Server Init
```bash
sudo apt update
sudo apt upgrade
sudo apt install nvidia-410
...

```

## Train
- Install dependencies:
```bash
pip3 install -r requirements.txt
```

```bash
python3 main.py --mode train --cuda True
python3 main.py --mode test --cuda True
python3 main.py --mode train --cuda True --warm-up True
```

## Train trick
- 1e-4 lr train 30-epoch for warm-up
- 2e-5 lr train about 10-epoch for fine tuning

### References:
- https://github.com/lsvih/Simple-Gimp-Captcha-Resolver
- https://github.com/ypwhs/captcha_break

# Generating Data

1. Deploy a Digitalocean's CPU optimized Debian 9 droplet
2. Convert dropted to Arch Linux by using this script: https://github.com/gh2o/digitalocean-debian-to-arch
3. Install dependencies: 
```bash
pacman -S screen bash-completion go php php-cgi php-gd libjpeg libpng fontconfig
git clone https://github.com/hexvalid/ai-bdc
go get github.com/cheggaaa/pb/v3
cd ai-bdc/gen
mkdir out && mkdir /tmp/bdc_void/
```

4. Change style (optional):
```bash
# Edit $STYLE=
# -1 mean random
nano php/botdetect-captcha-lib/botdetect/CaptchaIncludes.php
```

5. Prepare PHP-CGI:
```bash
# Open up 10 php-cgi server via screen
php -c php/php.ini -t php/ -S 127.0.0.1:9000
php -c php/php.ini -t php/ -S 127.0.0.1:9001
php -c php/php.ini -t php/ -S 127.0.0.1:9002
......
php -c php/php.ini -t php/ -S 127.0.0.1:9009
```

6. Run Generation Script:
```bash
# Check mode and count
# Run mode = 1 and mode = 2
# As twice!
nano gen.go

# run via screen
go run gen.go
```

7. Package:
```bash
tar cfJ gendata-count250k-style18.tar.xz out/
```

8. Clean:
```bash
  rm /tmp/bdc_pipe && touch /tmp/bdc_pipe && rm -rf /tmp/bdc_void/ && mkdir /tmp/bdc_void/
```

9. Publish
```bash
go get github.com/github-release/github-release
export GITHUB_TOKEN=...
go/bin/github-release upload --user hexvalid --repo ai-bdc --tag gendatas --name "gendata-count250k-style18.tar.xz" --file gendata-count250k-style18.tar.xz
```

