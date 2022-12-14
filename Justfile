# How to install just: https://just.systems/man/en/chapter_4.html

build:
    go build -v -o process

clean:
    rm -f process

run threads mode="gain" tickers="AAPL,MSFT,GOOG,TSLA,AMZN,INTC,AMD,NVDA":
    NUM_THREADS={{threads}} go run main.go -tickers {{tickers}} -mode {{mode}}

get-all:
    #!/bin/bash
    set -euo pipefail
    while read -r stock; do
        # check if stock is already downloaded
        if [ ! -f ./data/$stock.csv ]; then
            just get $stock
        fi
        just check data/$stock.csv
    done < tickers.txt

zip-all:
    zip -r data.zip data

unzip-all:
    unzip data.zip

check file:
    #!/bin/bash
    set -euo pipefail
    # check if file contains data
    if [[ $(wc -l < {{file}}) -le 1 ]]; then
        echo "File is empty"
        rm {{file}}
    fi

fix:
    go run cmd/fix/fix.go

clean-data:
    rm -rf data

get stock:
    python scripts/get_a_stock.py {{stock}}
