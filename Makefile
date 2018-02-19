
help:
	@echo "Go CI 0.0.4:"
	@echo "Available commands:"
	@echo "\tmake install			Install dependencies."
	@echo "\tmake run			Run default command."
	@echo "\tmake test			Run tests."
	@echo "\tmake coverage			Show coverage in html."
	@echo "\tmake clean			Clean build files."

install:
	@echo "Make: Install"
	./scripts/install.sh

run:
	@echo "Make: Run"
	./scripts/run.sh

.PHONY: test
test:
	@echo "Make: Test"
	./scripts/test.sh

coverage:
	@echo "Make: Coverage"
	./scripts/cover.sh

clean:
	@echo "Make: Clean"
	./scripts/clean.sh


## Container related commands.

container-install:
	@echo "Make: Install"
	./container/install.sh

.PHONY: test
container-test:
	@echo "Make: Test"
	./container/test.sh

container-coverage:
	@echo "Make: Coverage"
	./container/cover.sh

container-clean:
	@echo "Make: Clean"
	./container/clean.sh
