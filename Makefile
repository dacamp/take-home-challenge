include go-build/include.mk

bins: vendor
	@$(GO) build -o $(CURDIR)/bin/challenge-executable -i $(BUILD_FLAGS) $(BUILD_GC_FLAGS)

cover:
	@rm -f $(COVERFILE);
	@echo "mode: count" > $(COVERFILE);
	@grep -h -v "mode: " $(TMP)/*.cover >> $(COVERFILE);
	@$(GO) tool cover -html=$(COVERFILE)

test: $(TMP) vendor
	@$(foreach pkg, $(ALL_PKGS), \
		 $(GO) test -v -race $(TEST_FLAGS) -coverprofile $(TMP)/$(lastword $(subst /, ,$(pkg))).cover $(pkg) || true;)
