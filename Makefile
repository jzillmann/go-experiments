.PHONY: dev server ui

# server & ui
dev: 
	@$(MAKE) server-dev & $(MAKE) ui

prep:
	@$(MAKE) tidy & $(MAKE) check


# server
tidy:
	@echo "Cleaning up server imports..."
	@cd server && go mod tidy


server:
	@echo "Starting the server..."
	@cd server && go run main.go

server-dev:
	@echo "Starting the server..."
	@cd server && wgo run -verbose main.go


# UI
ui:
	@echo "Starting the UI development server..."
	@cd ui && npm run dev

check:
	@cd ui && npm run check
