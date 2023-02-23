.PHONY: release

release: $(verson)
	@echo "Release Version: $(version)"
	@echo $(version) > VERSION
	git tag $(version)
	git add VERSION
	git commit -m $(version)

beta: $(verson)
	@echo "Beta Version: $(version)-beta"
	@echo $(version)-beta > VERSION
	git add VERSION
	git commit -m $(VERSION)-beta