# Simply🚀Fire
В репозитории указан код релизной версии программы SimplyFire, предназначенной для автоматизации составления должностных инструкций (ДИ) сотрудников компаний.
Программа разработана на языке программирования Golang для компьютеров под управлением операционной системы Windows.
Программа разработана в мае 2024 г.

# Логика работы программы
При заполнении и выбора соответствующих полей программы и нажатии кнопки "Сформировать должностную инструкцию" программа связывается с корпоративной базой данных (БД) компании, в которой содержатся сведения о ее сотруднике, необходимые для составления ДИ, и как только произошла верификация введенных и выбранных пользователем данных с данными из БД, происходит подстановка релевантных значений о сотруднике компании из БД в ключи, указанные в шаблонизаторе ДИ, с последующей выгрузкой готовой ДИ.

# Содержание репозитория
Репозиторий состоит из двух частей:

Первая часть - это сам исходный код (main.go & go.mod & go.sum) с необходимыми вложениями к нему (SFire_Icon & Logo & шаблонизаторы ДИ (все word-файлы, начинающиеся с "ДИ...")).

ВАЖНО: все файлы первой части должны находиться в одной папке при загрузке в редактор кода.

Вторая часть - это дополнительные материалы по самой программе и по тому, как она работает ("Дополнительные артефакты программы" в формате ZIP): диаграммы последовательности и руководство (техническое) пользователя.
