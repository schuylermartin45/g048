##
## File:        Makefile
##
## Author:      Schuyler Martin <schuylermartin45@gmail.com>
##
## Description: Builds G048.
##

# Directories
BIN = ./bin/
SRC = ./src/

# Go Compiler
GC = go
GFLAGS = build


# Primary build directive
build:
	$(GC) $(GFLAGS) -o $(BIN)g048 $(SRC)g048

# Install dependencies
depend:
	$(GC) get github.com/gdamore/tcell

# Clean directive
clean:
	rm -rf $(BIN)*
