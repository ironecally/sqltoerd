ifeq ($(output),)
	OUTPUT_ARGS :=
else
	OUTPUT_ARGS := -o $(output)
endif


init:
	@echo " >> initializing project repo"
	mkdir -p result
	mkdir -p source
	mkdir -p bin

build:
	@echo " >> building binary"
	go build -o bin/sqltoerd-cli ./cmd/cli


run:
	@make build
	@echo " >> running binary"
	@./bin/sqltoerd-cli -i $(filename) $(OUTPUT_ARGS)
	@echo " >> process complete, please check result/ folder"