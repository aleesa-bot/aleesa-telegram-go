# Перечень того, что надо реализовать.

В качестве основы возьмём https://github.com/NicoNex/echotron 

Для реализации телеграммного фронтэнда нам надо сделать:

0. Парсер конфига.
   [V] Реализован начальный вариант парсера конфига.

1. Разгребальщик апдейтов и формирователь redis-ных сообщений.
   [ ] Написана заглушка для парсера апдейтов bot api телеграма

2. Разгребальщик редисных сообщений и формирователь запросов в tg bot api.
   [ ] Написан черновой вариант рагребателя

3. Обработчик группы команд admin, сохранение настроек в pebbledb.

   1 цензор сообщений определённого типа.

   2 фортунка каждое утро. (+ добавить время, относительно gmt, когда показывать фортунку).

   3 Приветствие новых участников

   4 oboobs

   5 obutts

   6 chan_msg

   7 user ban

   8 user mute

4. help

5. Readme.txt
