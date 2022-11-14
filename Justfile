# How to install just: https://just.systems/man/en/chapter_4.html

build:
    go build -v -o process

clean:
    rm -f process

run:
    go run main.go

get stock:
    if [ ! -d ./data ]; then mkdir data; fi
    curl -o ./data/{{stock}}.csv "https://stooq.com/q/d/l/?s={{stock}}.us&i=d"
