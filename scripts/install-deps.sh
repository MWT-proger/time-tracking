#!/bin/bash

# Определение дистрибутива Linux
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
else
    echo "Не удалось определить дистрибутив Linux"
    exit 1
fi

# Установка зависимостей в зависимости от дистрибутива
case $OS in
    ubuntu|debian|linuxmint)
        echo "Установка зависимостей для Ubuntu/Debian..."
        sudo apt update
        sudo apt install -y libayatana-appindicator3-dev
        ;;
    fedora)
        echo "Установка зависимостей для Fedora..."
        sudo dnf install -y libappindicator-gtk3-devel
        ;;
    arch|manjaro)
        echo "Установка зависимостей для Arch Linux..."
        sudo pacman -S --noconfirm libappindicator-gtk3
        ;;
    *)
        echo "Неподдерживаемый дистрибутив: $OS"
        echo "Пожалуйста, установите зависимости вручную:"
        echo "- libayatana-appindicator3-dev или libappindicator3-dev"
        exit 1
        ;;
esac

echo "Зависимости успешно установлены!" 