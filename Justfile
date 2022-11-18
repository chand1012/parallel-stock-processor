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

fix file:
    #!/bin/bash
    set -euo pipefail
    # read each line and make sure it has 6 columns
    # if the line has less than 6 columns, delete the line
    CMD="d"
    while read -r line; do
        if [[ $(echo $line | awk -F, '{print NF}') -ne 6 ]]; then
            echo "Line has less than 6 columns"
            sed "/$line/d" {{file}} > {{file}}.tmp
            mv {{file}}.tmp {{file}}
        fi
    done < {{file}}

fix-all:
    #!/bin/bash
    set -euo pipefail
    # get all files in data directory
    files=$(find ./data -type f)
    echo "Fixing $(ls -1 data | wc -l) files"
    for file in $files; do
        echo "Fixing $file"
        just fix $file
    done

clean-data:
    rm -rf data
