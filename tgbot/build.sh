export CGO_ENABLED=1;

if [ ! -d "./output" ]; then
    mkdir ./output
fi

go build -ldflags "-s -w" -o ./output/tgbot
