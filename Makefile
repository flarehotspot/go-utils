PLUGINS = com.adopisoft.basic-flare-theme com.flarego.basic-net-mgr com.flarego.basic-system-account

openwrt: export CGO_ENABLED=1
plugin: export CGO_ENABLED=1

default: clean
	go build -race -ldflags="-s -w" -o flarehotspot.app -tags="mono dev" main/main_mono.go
	./flarehotspot.app

plugin: clean
	cd core && make plugin
	cd main && make plugin
	bash ./plugins-action.sh "make plugin" &
	./main/app

build: clean
	go build -ldflags="-s -w" -trimpath -o flarehotspot.app -tags="mono dev" main/main_mono.go

openwrt:
	ar -rc /usr/lib/libpthread.a
	ar -rc /usr/lib/libresolv.a
	ar -rc /usr/lib/libdl.a
	rm -rf .cache public
	cd core && make plugin_prod
	cd main && make plugin
	bash ./plugins-action.sh "make plugin" &
	./main/app

sync:
	scp -O -r $(PWD)/core root@$(remote):/root/flarehotspot
	scp -O -r $(PWD)/goutils root@$(remote):/root/flarehotspot
	scp -O -r $(PWD)/sdk root@$(remote):/root/flarehotspot

sync_all:
	scp -O -r $(PWD) root@$(remote):/root/flarehotspot

plugin:
	rm -rf .cache public
	cd core && make plugin
	bash ./plugins-action.sh "make plugin"
	cd main && make plugin
	./main/app

clean:
	rm -rf .tmp .cache public *.app
	find . -name "*.so" -type f -delete
	find . -name "*.app" -type f -delete

pull:
	cd core && git pull &
	cd sdk && git pull &
	cd goutils && git pull &
	cd hardware-db && git pull &
	bash ./plugins-action.sh "git pull" &
	git pull &

push:
	cd core && git push &
	cd sdk && git push &
	cd goutils && git push &
	cd hardware-db && git push &
	bash ./plugins-action.sh "git push" &
	git push &

checkout_main:
	cd core && git checkout main &
	cd goutils && git checkout main &
	cd sdk && git checkout main &
	cd hardware-db && git checkout main
	bash ./plugins-action.sh "git checkout main" &
