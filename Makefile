WEBDIR = website
SRVDIR = server

WEBDST = $(SRVDIR)/public
SRVBIN = $(SRVDIR)/bigtechfaild

.PHONY: website server

default: website server

website:
	cd $(WEBDIR) && hugo -d ../$(WEBDST)

server:
	cd $(SRVDIR) && go build

run: website server
	$(SRVBIN) -p 8080

clean:
	rm -rf $(WEBDST) $(SRVBIN)
