# Makefile

MAGE := $(shell command -v mage 2> /dev/null)
# MAGE_VERBOSE ?= 0

# Handle call without arguments
# Must be first
.PHONY: default
default:
ifndef MAGE
	make mage
endif
	mage

# Make silent default
# Must not be first
# ifndef VERBOSE
# .SILENT:
# endif

# Install mage binary
# https://magefile.org/
.PHONY: mage
mage:
ifndef MAGE
	go get -u -d github.com/magefile/mage
	cd $$GOPATH/src/github.com/magefile/mage; \
		go run bootstrap.go
endif

# Last-resort rule
# https://www.gnu.org/software/make/manual/make.html#Last-Resort
.DEFAULT: mage

# Pass arguments to mage
%::
ifndef MAGE
	make mage
endif
	mage $@
# ifeq ($(MAGE_VERBOSE),1)
# 	mage -v $@
# else
# 	mage $@
# endif
