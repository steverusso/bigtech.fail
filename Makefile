WEBDIR = website
SRVDIR = server

WEBDST = $(SRVDIR)/public
SRVBIN = $(SRVDIR)/bigtechfaild

.PHONY: website server

default: server

website:
	cd $(WEBDIR) && hugo -d ../$(WEBDST)

server: website
	cd $(SRVDIR) && go build

run: server
	$(SRVBIN) -p 8080

install: server
	sudo sysrc bigtechfaild_enable="YES"
	sudo cp -f $(SRVBIN) /usr/local/bin/
	sudo cp -f $(SRVDIR)/contrib/bigtechfaild.rc.d /usr/local/etc/rc.d/bigtechfaild
	sudo service bigtechfaild restart

clean:
	rm -rf $(WEBDST) $(SRVBIN)
