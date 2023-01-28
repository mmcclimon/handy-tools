# So it's come to this.

CMD_DIR := ./cmd
OUT_DIR := ./out
INTERNAL_DIR := ./internal
BIN_DIR := $(HOME)/bin

INTERNALS := $(shell find $(INTERNAL_DIR) -type f -name '*.go')
GO_CMDS := $(notdir $(shell find $(CMD_DIR) -maxdepth 1 -mindepth 1 -type d))
GO_TARGETS := $(addprefix $(OUT_DIR)/, $(GO_CMDS))

INSTALL_TARGETS := $(addprefix $(BIN_DIR)/, $(GO_CMDS))

all: $(GO_TARGETS)

$(OUT_DIR)/%: $(CMD_DIR)/%/*.go $(INTERNALS)
	go build -o $@ $(CMD_DIR)/$(notdir $@)

install: $(INSTALL_TARGETS)

$(BIN_DIR)/%: $(OUT_DIR)/%
	cp $< $@
