GO = go
GOBIN = $(GOPATH)/bin
GODEP = $(GOBIN)/dep
GOLINT = $(GOBIN)/golint
GOLIST = *.go $(shell $(GO) list ./... | grep -v vendor) 
OVERALLS = $(GOBIN)/overalls
GOVERALLS = $(GOBIN)/goveralls

PKG = github.com/propervillain/moistpetal
DIR = $(GOPATH)/src/$(PKG)
CMD = $(shell $(GO) list ./... | grep -v vendor | grep /cmd)

.PHONY: all dep update gen fmt vet lint test build install list profile coveralls clean

all: gen check build

check: fmt vet lint test

dep: $(GODEP)
	@$(GODEP) ensure

update: $(GODEP)
	@$(GODEP) ensure -update

gen:
	@echo "[ go generate ]"
	@$(foreach p,$(GOLIST), \
		echo $p; \
		$(GO) generate $p || exit 1; \
	)

fmt:
	@echo "[ go fmt ]"
	@$(foreach p,$(GOLIST), \
		echo $p; \
		$(GO) fmt $p; \
	)

vet:
	@echo "[ go vet ]"
	@$(foreach p,$(GOLIST), \
		echo $p; \
		$(GO) vet $p || exit 1;\
	)

lint: $(GOLINT)
	@echo "[ golint ]"
	@$(foreach p,$(GOLIST), \
		echo $p; \
		$(GOLINT) $p; \
	)

test: 
	@$(foreach p,$(GOLIST), \
		$(GO) test -short -race -v $p || exit 1;\
	)

build: 
	@echo "[ go build ]"
	@$(foreach p,$(CMD), \
		echo $p; \
		$(GO) build $p; \
	)

install: 
	@echo "[ go install ]"
	@$(foreach p,$(CMD), \
		echo $p; \
		$(GO) install $p; \
	)

list:
	@$(foreach p,$(GOLIST), \
		echo $p; \
	)

profile: $(OVERALLS)
	$(OVERALLS) -project=$(PKG) -covermode=atomic \
		-ignore=.git,.github,vendor -debug \
		-- -short -race -v \

coveralls: $(GOVERALLS)
	$(GOVERALLS) -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $(COVERALLS_TOKEN)

clean:
	@$(foreach p,$(CMD), \
		echo "rm -f" $(shell basename "$p"); \
		basename "$p" | xargs rm -f; \
	)
	rm -rf ./log/testdata

$(GODEP):
	go get -u github.com/golang/dep/cmd/dep

$(GOLINT):
	go get -u github.com/golang/lint/golint

$(OVERALLS):
	go get -u github.com/go-playground/overalls

$(GOVERALLS):
	go get -u github.com/mattn/goveralls
