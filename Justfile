# How to install just: https://just.systems/man/en/chapter_4.html

build:
    go build -v -o process

clean:
    rm -f process

run:
    go run main.go

get stock:
    mkdir -p data
    curl -o ./data/{{stock}}.csv "https://stooq.com/q/d/l/?s={{stock}}.us&i=d"

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
