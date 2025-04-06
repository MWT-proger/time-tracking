.PHONY: build run debug install deb desktop install-opt

BUILD_DATE := $(shell date -u +%Y-%m-%d)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/MWT-proger/time-tracking/internal/app.BuildDate=$(BUILD_DATE) -X github.com/MWT-proger/time-tracking/internal/app.GitCommit=$(GIT_COMMIT)

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
	mkdir -p debian/usr/local/bin
	cp ttracker debian/usr/local/bin/
	dpkg-deb --build debian
	mv debian.deb ttracker.deb 

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