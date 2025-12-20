_APPLICATION_NAME          := slumber
_APPLICATION_VERSION_MAJOR := 0
_APPLICATION_VERSION_MINOR := 1
_APPLICATION_VERSION_PATCH := 0

# The path to the main package
_MAIN_PACKAGE_PATH  := .

# The path to the binary
_BINARY_NAME        := $(_APPLICATION_NAME)

.DEFAULT_GOAL := help

include help.mk
include go.mk
-include custom.mk
