BUILDER=./.builder
RULES=go
include $(BUILDER)/rules.mk
$(BUILDER)/rules.mk:
	-go run endobit.io/builder@latest init
