#!/usr/bin/env python3
import time
import json
import os
import click
from datetime import datetime, timedelta

# Путь к файлу для хранения данных — в домашней директории пользователя
DATA_FILE = os.path.join(os.path.expanduser("~"), "учет_времени.json")


# Функция загрузки данных из файла
def load_data():
    if os.path.exists(DATA_FILE):
        with open(DATA_FILE, "r") as file:
            return json.load(file)
    else:
        return {}


# Функция сохранения данных в файл
def save_data(data):
    with open(DATA_FILE, "w") as file:
        json.dump(data, file, ensure_ascii=False, indent=4)


# Команда для создания нового проекта
@click.command()
@click.option(
    "--project", prompt="Название проекта", help="Название проекта для создания"
)
def create(project):
    data = load_data()
    if project in data:
        click.echo(f'Проект "{project}" уже существует.')
    else:
        data[project] = {"entries": []}
        save_data(data)
        click.echo(f'Проект "{project}" успешно создан.')


# Команда для начала отслеживания времени на проекте
@click.command()
@click.option("--project", help="Название проекта для отслеживания времени")
def start(project):
    data = load_data()
    if not data:
        click.echo("Нет доступных проектов. Сначала создайте проект командой `create`.")
        return

    # Если проект не указан, предложить выбрать из доступных
    if not project:
        project = click.prompt("Выберите проект", type=click.Choice(data.keys()))
    elif project not in data:
        click.echo(f'Проект "{project}" не найден. Проверьте правильность названия.')
        return

    # Проверка, что отслеживание не запущено
    if "start_time" in data[project]:
        click.echo(f'Отслеживание времени для проекта "{project}" уже запущено.')
        return

    # Запуск отслеживания
    data[project]["start_time"] = time.time()
    save_data(data)
    click.echo(f'Начато отслеживание времени для проекта "{project}".')


# Команда для остановки отслеживания времени на проекте
@click.command()
@click.option("--project", help="Название проекта для остановки отслеживания времени")
def stop(project):
    data = load_data()
    if not data:
        click.echo("Нет доступных проектов. Сначала создайте проект командой `create`.")
        return

    # Если проект не указан, предложить выбрать из доступных
    if not project:
        project = click.prompt("Выберите проект", type=click.Choice(data.keys()))
    elif project not in data:
        click.echo(f'Проект "{project}" не найден. Проверьте правильность названия.')
        return

    # Проверка, что отслеживание было запущено
    if "start_time" not in data[project]:
        click.echo(f'Нет активного отслеживания времени для проекта "{project}".')
        return

    # Запрос описания работы
    description = click.prompt("Что сделано", default="")

    # Остановка отслеживания и сохранение данных
    start_time = data[project].pop("start_time")
    elapsed_time = time.time() - start_time
    entry = {
        "time_spent": int(elapsed_time),
        "description": description,
        "date": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
    }
    data[project]["entries"].append(entry)
    save_data(data)

    elapsed_str = str(timedelta(seconds=int(elapsed_time)))
    click.echo(
        f'Остановлено отслеживание для проекта "{project}". Время: {elapsed_str}. Описание: {description}'
    )


# Команда для вывода сводки по проектам
@click.command()
def summary():
    data = load_data()
    if not data:
        click.echo("Нет данных для отображения.")
        return

    for project, info in data.items():
        click.echo(f'\nПроект "{project}":')
        total_time = sum(entry["time_spent"] for entry in info.get("entries", []))
        total_str = str(timedelta(seconds=int(total_time)))
        click.echo(f"  Общее время: {total_str}")

        # Вывод всех записей по проекту
        for entry in info.get("entries", []):
            entry_time = str(timedelta(seconds=entry["time_spent"]))
            click.echo(f'    {entry["date"]} - {entry_time}: {entry["description"]}')


# Главная функция CLI для объединения команд
@click.group()
def cli():
    pass


cli.add_command(create)
cli.add_command(start)
cli.add_command(stop)
cli.add_command(summary)

if __name__ == "__main__":
    cli()
