all:
	cd client && go build && mv HTCPCP-client* ..
	cd server && go build && mv HTCPCP-server* ..

clean:
	rm -rf HTCPCP-* conf.json