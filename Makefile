SRC=./init.go
BIN=./nanowiki

all: $(BIN)

$(BIN): $(SRC)
	go build -o $(BIN) $(SRC)
