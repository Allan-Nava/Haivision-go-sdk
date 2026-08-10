.PHONY: help build test fmt fmt-check vet cover backlog-lint roadmap roadmap-check release-notes release release-dry check
#
# La root del modulo non contiene file Go: il target di build è sempre ./... (mai `go build .`).
#
help:
	@echo "build         - compila tutti i package"
	@echo "test          - build + test"
	@echo "cover         - test con coverage per package"
	@echo "fmt           - gofmt -w su tutto il repo"
	@echo "fmt-check     - fallisce se qualcosa non è formattato (gate CI)"
	@echo "vet           - go vet"
	@echo "check         - fmt-check + vet + build + test + backlog-lint + roadmap-check"
	@echo "backlog-lint  - valida docs/backlog.md e la catena di versioni (gate CI)"
	@echo "roadmap       - rigenera docs/roadmap.md dal backlog"
	@echo "roadmap-check - fallisce se docs/roadmap.md non è in pari col backlog (gate CI)"
	@echo "release-notes V=vX.Y.Z - stampa la sezione CHANGELOG di quella milestone"
	@echo "release-dry    - mostra quale release verrebbe tagliata (non scrive nulla)"
	@echo "release        - versione dal backlog + CHANGELOG + gate + commit + tag (NO push)"
#
build:
	go build -v ./...
#
test: build
	go test -v ./...
#
cover:
	go test -cover ./...
#
fmt:
	gofmt -w .
#
fmt-check:
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then echo "✗ file non formattati (lancia 'make fmt'):"; echo "$$out"; exit 1; fi; \
	echo "✅ gofmt OK"
#
vet:
	go vet ./...
#
# --- backlog & milestone di versione (vedi docs/backlog.md) -------------------------
#
backlog-lint:
	python3 scripts/backlog-lint.py
#
roadmap:
	python3 scripts/generate-roadmap.py
#
roadmap-check:
	python3 scripts/generate-roadmap.py --check
#
# uso: make release-notes V=v1.1.0
release-notes:
	@test -n "$(V)" || { echo "uso: make release-notes V=vX.Y.Z"; exit 2; }
	@python3 scripts/generate-roadmap.py --release-notes $(V)
#
# La versione NON si passa: la ricava dal backlog (prima milestone pendente con 0 item open).
# `V=vX.Y.Z` la forza; ARGS passa altri flag (--no-commit, --allow-branch, --date …).
release-dry:
	python3 scripts/new-release.py --dry-run $(if $(V),--version $(V),) $(ARGS)
#
release:
	python3 scripts/new-release.py $(if $(V),--version $(V),) $(ARGS)
#
check: fmt-check vet build test backlog-lint roadmap-check
	@echo "✅ tutti i gate passati"
#
