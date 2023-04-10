BUILDER=./.builder
RULES=go
include $(BUILDER)/rules.mk
$(BUILDER)/rules.mk:
	-go run github.com/endobit/builder@latest init

build::
	cd sample && $(GO_BUILD) .

clean::
	rm -f sample/sample
