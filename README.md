# Репозиторий курса "Разработка веб-сервисов на Golang"

## Дисклеймер

Это факультатив. Запись на него - ваше сознательное решение, никто не заставлял.

Вы пришли сюда за знаниями, а не за оценками. Поэтому все домашки надо делать самостоятельно, без консультаций с другими студентами, без просмотра чужих решений, особенно если вы сами ничего не сделали еще сами, в крайнем случае - спрашивать препода.

## Основные правила

1. При обнаружении списываний или заимствований вы отчисляетесь с курса.
2. Можно ходить на лекции и не делать домашки. Можно не делать домашку если она вам не нравится и перейти к другой.
3. Преподы на вашей стороне и всегда готовы помочь. Можно задавать вопросы в телеграме в личку и в общий чат.
4. Домашки сдаём тому преподавателю, который вёл занятие.
5. Код и тем более решения домашек в паблик открывать нельзя, репозиторий должен быть приватным.
6. Иметь публичный репозиторий где-то кроме gitlab нельзя

## Прочие правила

1. Преподавательский состав оставляет за собой право дополнять правила.
2. Хардкод (код работающий под частное условие) запрещён. За первый раз - предупреждение, за последующие -1 балл. Можно спрашивать будет ли что-то хардкодом до сдачи задания. Весь код должен работать максимально универсально.
3. Домашки пишем там же где лежит вводная (например, 1/99_hw/XXX), другие папки не создаём
4. Тесты домашек править нельзя.
5. Вопросы задавать четко, конкретно: "я делаю Х, получаю Y, а хочу получить Z".
6. Студент должен иметь реальное имя-фамилию-фото в гитлабе и на портале. Реальные ФИО и фото в телеграме так же желательны.
7. Домашку надо коммитить в свою репу, создавать merge request в основную репу не надо. (про сдачу читайте чуть подробнее в одном из следующих разделов)
8. Домашки предназначены для выполнения индивидуально и самостоятельно. Это значит, что нельзя делать их группой, нельзя обсуждать как делать, нельзя показывать свои решения(это тоже карается).
9. Преподавательский состав оставляет за собой право не принимать мутные и/или некрасивые решения домашек. В этом случае необходимо поправить замечания без препирательств.
10. Консультации и проверки заданий даются в основном вечером
11. Если вы пишете в 2 ночи - не надо писать в 9 утра вопросы "а вы посмотрели?" - в подобных случаях мы скорее всего посмотрим только вечером и надо напомнить про себя после 19 часов.
12. Халявы не будет, домашки сложные, придётся работать.

## После получения доступа к репозиторию

Вам нужно форкнуть текущий репозиторий к себе. Необходимо сделать его приватным и выдать доступы для @skinass, @Ksenobait09, @vpersiyanova, @a-kuchin с уровнем доступа maintainer. Никому другому доступы давать нельзя.

[Скриншот с примером](https://s.mail.ru/7XKz/fmJyoaZMA)

## Правила сдачи ДЗ

1. На выполнение и сдачу домашнего задания даётся три недели. Дедлайном является 21:00:00 субботы.
2. Домашнее задание выполненное в срок оценивается в 10 баллов. После дедлайна - 5 баллов.
3. Домашнее задание считается выполненным после того, как преподаватель принял его. Не в момент отсылки на проверку.
4. *Дисклеймер.* Преподаватели не бывают онлайн 24/7, то есть если отправили на проверку за час до дедлайна, то
преподаватель может не успеть посмотреть ваше решение.

## Последовательность выполнения ДЗ

1. Нужно подтянуть изменения основного репозитория в свой форк.
2. Создаём новую ветку (c именем hw_X) в воем репозитории
3. Читаем задание в X/99_hw/X.md
4. Пишем своё решение в той же папке, где и лежит условие
5. Доводим код до состояния прохождения тестов
6. Не забываем форматировать код (gofmt или goimports, если ваша IDE этого не делает автоматически)
7. Создаём Merge Request из созданной ветки(п. 2) в master форкнутого репозитория
8. Стучимся в личку к преподавателю, который вёл соответствующую лекцию, с ссылкой на Merge Request и просьбой поревьювить.

### Как подтянуть изменения основного репозитория в форк

```bash
# будучи в своём репозитории
git pull https://gitlab.com/mailru-go/lectures-2021-2.git master
git push origin master
```
