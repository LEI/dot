# Makefile

MAGE := $(shell command -v mage 2> /dev/null)
# MAGE_VERBOSE ?= 0

# RUN_ARGS := $(MAKECMDGOALS)
ifeq (mage,$(firstword $(MAKECMDGOALS)))
RUN_ARGS := $(filter-out mage,$(MAKECMDGOALS))
# # https://stackoverflow.com/a/14061796/7796750
$(eval $(RUN_ARGS):;@:)
else
RUN_ARGS := $(filter-out default,$(MAKECMDGOALS))
endif

.DEFAULT: default

.PHONY: default
# $(filter-out $@,$(RUN_ARGS))
default:
	$(MAGE) $(RUN_ARGS)

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
	cd $(GOPATH)/src/github.com/magefile/mage; \
		go run bootstrap.go # $(RUN_ARGS)
endif
	type mage
	mage --version

ifneq (0,$(words $(RUN_ARGS)))
# ifneq (default,$(firstword $(RUN_ARGS)))
ifneq (mage,$(firstword $(MAKECMDGOALS)))
# https://stackoverflow.com/a/11731328/7796750
%: default
	@:
endif
# endif
endif

# # https://www.gnu.org/software/make/manual/make.html#Last-Resort
# .DEFAULT:

# %::
# 	mage $(RUN_ARGS)
