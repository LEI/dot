# Makefile

MAGE := $(shell command -v mage 2> /dev/null)
# MAGE_VERBOSE ?= 0

RUN_ARGS := $(MAKECMDGOALS)
RUN_ARGS := $(filter-out default,$(RUN_ARGS))
RUN_ARGS := $(filter-out mage,$(RUN_ARGS))
# # https://stackoverflow.com/a/14061796/7796750
# $(eval $(RUN_ARGS):;@:)

#
.PHONY: default
default:
	$(MAGE) $(filter-out $@,$(RUN_ARGS))

# # Make silent default
# # Must not be first
# ifndef VERBOSE
# .SILENT:
# endif

# Install mage binary https://magefile.org/
.PHONY: mage
mage:
ifndef MAGE
	go get -u -d github.com/magefile/mage
	cd $$GOPATH/src/github.com/magefile/mage; \
		go run bootstrap.go
endif

# ifneq (0,$(words $(RUN_ARGS)))
# https://stackoverflow.com/a/11731328/7796750
%: default
	@:
# endif

# # Last-resort rule
# # https://www.gnu.org/software/make/manual/make.html#Last-Resort
# .DEFAULT:

# # Pass arguments to mage
# %::
# 	@echo RUN_ARGS $(RUN_ARGS)
