.PHONY: build run debug install deb desktop install-opt release release-with-deb install-deps

VERSION=$(shell grep -oP 'Version = "\K[^"]+' internal/app/version.go)
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u +"%Y-%m-%d")
LDFLAGS := -X github.com/MWT-proger/time-tracking/internal/app.BuildDate=$(BUILD_DATE) -X github.com/MWT-proger/time-tracking/internal/app.GitCommit=$(COMMIT)

build:
	go build -ldflags "$(LDFLAGS)" -o ttracker ./cmd/time-tracker

run: build
	./ttracker -data ./data.json -notify-time 5

debug: build
	./ttracker -data ./data.json -notify-time 5 -log-level debug 

install: build
	sudo cp ttracker /usr/local/bin/
	@echo "ttracker установлен в /usr/local/bin/" 

deb: build
	mkdir -p debian/DEBIAN
	mkdir -p debian/usr/local/bin
	mkdir -p debian/usr/share/applications
	mkdir -p debian/usr/share/pixmaps
	
	# Копирование исполняемого файла
	cp ttracker debian/usr/local/bin/
	
	# Копирование иконки
	cp internal/app/systray/assets/icons/clock.png debian/usr/share/pixmaps/ttracker.png
	
	# Создание .desktop файла
	echo "[Desktop Entry]" > debian/usr/share/applications/ttracker.desktop
	echo "Name=Time Tracker" >> debian/usr/share/applications/ttracker.desktop
	echo "Comment=Приложение для отслеживания времени" >> debian/usr/share/applications/ttracker.desktop
	echo "Exec=/usr/local/bin/ttracker" >> debian/usr/share/applications/ttracker.desktop
	echo "Icon=ttracker" >> debian/usr/share/applications/ttracker.desktop
	echo "Terminal=true" >> debian/usr/share/applications/ttracker.desktop
	echo "Type=Application" >> debian/usr/share/applications/ttracker.desktop
	echo "Categories=Utility;" >> debian/usr/share/applications/ttracker.desktop
	
	# Создание control файла
	echo "Package: ttracker" > debian/DEBIAN/control
	echo "Version: $(VERSION)" >> debian/DEBIAN/control
	echo "Section: utils" >> debian/DEBIAN/control
	echo "Priority: optional" >> debian/DEBIAN/control
	echo "Architecture: amd64" >> debian/DEBIAN/control
	echo "Maintainer: $(shell git config --get user.name) <$(shell git config --get user.email)>" >> debian/DEBIAN/control
	echo "Depends: libayatana-appindicator3-1 | libappindicator3-1" >> debian/DEBIAN/control
	echo "Description: Приложение для отслеживания времени, затраченного на проекты" >> debian/DEBIAN/control
	echo " Позволяет создавать проекты, отслеживать время работы," >> debian/DEBIAN/control
	echo " добавлять описания к выполненным задачам и многое другое." >> debian/DEBIAN/control
	
	# Сборка пакета
	dpkg-deb --build debian
	mv debian.deb ttracker_$(VERSION)_amd64.deb

install-opt: build
	sudo mkdir -p /opt/ttracker
	sudo cp ttracker /opt/ttracker/
	sudo cp internal/app/systray/assets/icons/clock.png /opt/ttracker/ttracker.png
	@echo "ttracker установлен в /opt/ttracker/"

desktop: install-opt
	@echo "Создание .desktop файла и установка для текущего пользователя..."
	@echo "[Desktop Entry]" > ttracker.desktop
	@echo "Name=Time Tracker" >> ttracker.desktop
	@echo "Comment=Приложение для отслеживания времени" >> ttracker.desktop
	@echo "Exec=/opt/ttracker/ttracker" >> ttracker.desktop
	@echo "Icon=/opt/ttracker/ttracker.png" >> ttracker.desktop
	@echo "Terminal=true" >> ttracker.desktop
	@echo "Type=Application" >> ttracker.desktop
	@echo "Categories=Utility;" >> ttracker.desktop
	@echo "StartupWMClass=ttracker" >> ttracker.desktop
	mkdir -p ~/.local/share/applications
	cp ttracker.desktop /usr/share/applications/
	@echo "Приложение добавлено в меню приложений для текущего пользователя" 

release:
	@echo "Создание релиза версии $(VERSION)"
	git tag v$(VERSION)
	git push origin v$(VERSION)
	@echo "Релиз v$(VERSION) создан и отправлен в GitHub"

release-with-deb: deb release
	@echo "Загрузка .deb файла в GitHub Releases..."
	# Здесь можно использовать GitHub CLI или curl для загрузки файла
	# Пример с GitHub CLI:
	# gh release upload v$(VERSION) ttracker_$(VERSION)_amd64.deb
	@echo "Релиз v$(VERSION) с .deb файлом создан"

install-deps:
	@echo "Установка зависимостей..."
	./scripts/install-deps.sh
	@echo "Зависимости установлены!"
