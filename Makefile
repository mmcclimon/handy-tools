# So it's come to this.

CMD_DIR := ./cmd
OUT_DIR := ./out
INTERNAL_DIR := ./internal

GO_CMDS := $(notdir $(shell find $(CMD_DIR) -maxdepth 1 -mindepth 1 -type d))
GO_TARGETS := $(addprefix $(OUT_DIR)/, $(GO_CMDS))
INTERNALS := $(shell find $(INTERNAL_DIR) -type f -name '*.go')

all: $(GO_TARGETS)

$(OUT_DIR)/%: $(CMD_DIR)/%/*.go $(INTERNALS)
	go build -o $@ $(CMD_DIR)/$(notdir $@)
