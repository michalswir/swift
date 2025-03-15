# Build/rebuild and run containers
docker-run:
	@docker compose up --build -d

# Shutdown containers
docker-down:
	@docker compose down

# Test the application
test:
	@echo "Running integration tests..."
	@TESTCONTAINERS_RYUK_DISABLED=true go test ./internal/database -v
	
# Fetch all data by deleting and then inserting new data
fetch-data:
	@echo "Fetching all data from Interns_2025_SWIFT_CODES"
	@curl -X DELETE "http://localhost:8080/my/deleteAll" > /dev/null 2>&1
	@curl -X POST "http://localhost:8080/my/insertAll/1iFFqsu_xruvVKzXAadAAlDBpIuU51v-pfIEU5HeGa8w" > /dev/null 2>&1

# .PHONY defines targets that are not associated with files. 
.PHONY: docker-run docker-down test fetch-data