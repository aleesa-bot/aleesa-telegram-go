#!/usr/bin/env gmake -f

GOOPTS=CGO_ENABLED=0
BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath -buildvcs=false
MYNAME1=aleesa-telegram-go
BINARY1=${MYNAME1}
UNIX_BINARY1=${MYNAME1}
WINDOWS_BINARY1=${MYNAME1}.exe
MYNAME2=settings-migrator
BINARY2=${MYNAME2}
UNIX_BINARY2=${MYNAME2}
WINDOWS_BINARY2=${MYNAME2}.exe
RMCMD=rm -rf

# On windows binary name can depend not only on platform selected, but target too, we need no .exe suffix for linux bin.
ifeq ($(OS),Windows_NT)
ifdef GOOS
ifeq ($(GOOS),windows)
BINARY1=${WINDOWS_BINARY1}
BINARY2=${WINDOWS_BINARY2}
else  # not ifeq ($(GOOS),windows)
BINARY1=${MYNAME1}
BINARY2=${MYNAME2}
endif # ifeq ($(GOOS),windows)
else  # not ifdef GOOS
BINARY1=${WINDOWS_BINARY1}
BINARY2=${WINDOWS_BINARY2}
endif # ifdef GOOS
ifeq ($(SHELL), sh.exe)
RMCMD=DEL /Q /F
endif
endif

# Set newline symbol explicitly to aboud undefined behaviour in windows.
define IFS

endef


all: clean build


build:
ifeq ($(OS),Windows_NT)
# Looks like on windows gnu make explicitly set SHELL to sh.exe, if it was not set.
ifeq ($(SHELL), sh.exe)
#       # Vanilla cmd.exe / powershell.
	SET "CGO_ENABLED=0"
	go build ${BUILDOPTS} -o ${BINARY1} ./cmd/${MYNAME1}
else ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # git-bash
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY1} ./cmd/${MYNAME1}
else  # not ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # Some other shell.
#       # TODO: handle it.
	$(info "-- Dunno how to handle this shell: ${SHELL}")
endif # ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
else  # not  ($(OS),Windows_NT)
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${BINARY1} ./cmd/${MYNAME1}
endif # ifeq ($(OS),Windows_NT)


clean:
ifeq ($(OS),Windows_NT)
ifeq ($(SHELL),sh.exe)
#	# Vanilla cmd.exe / powershell.
	if exist ${WINDOWS_BINARY1} ${RMCMD} ${WINDOWS_BINARY1}
	if exist ${UNIX_BINARY1} ${RMCMD} ${UNIX_BINARY1}
	if exist ${WINDOWS_BINARY2} ${RMCMD} ${WINDOWS_BINARY2}
	if exist ${UNIX_BINARY2} ${RMCMD} ${UNIX_BINARY2}
else  # not ifeq ($(SHELL),sh.exe)
	${RMCMD} ./${WINDOWS_BINARY1}
	${RMCMD} ./${UNIX_BINARY1}
	${RMCMD} ./${WINDOWS_BINARY2}
	${RMCMD} ./${UNIX_BINARY2}
endif # ifeq ($(SHELL),sh.exe)
else  # not ifeq ($(OS),Windows_NT)
	${RMCMD} ./${UNIX_BINARY1}
	${RMCMD} ./${UNIX_BINARY2}
endif


upgrade:
	go get -u ./...
	go mod tidy
	go mod vendor


# Build settings-migrator binary.
${MYNAME2}: clean
ifeq ($(OS),Windows_NT)
# Looks like on windows gnu make explicitly set SHELL to sh.exe, if it was not set.
ifeq ($(SHELL), sh.exe)
#       # Vanilla cmd.exe / powershell.
	SET "CGO_ENABLED=0"
	go build ${BUILDOPTS} -o ${WINDOWS_BINARY2} ./cmd/${MYNAME2}
else ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # git-bash
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${WINDOWS_BINARY2} ./cmd/${MYNAME2}
else  # not ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
#       # Some other shell.
#       # TODO: handle it.
	$(info "-- Dunno how to handle this shell: ${SHELL}")
endif # ifeq (,$(findstring(Git/usr/bin/sh.exe, $(SHELL))))
else  # not  ($(OS),Windows_NT)
	CGO_ENABLED=0 go build ${BUILDOPTS} -o ${UNIX_BINARY2} ./cmd/${MYNAME2}
endif # ifeq ($(OS),Windows_NT)

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
