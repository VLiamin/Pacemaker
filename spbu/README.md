
## Структура проекта

* `cmd/cibgen`  – генератор (`golang`) тестовых данных для `pacemaker-schedulerd` и `crm_simulate`
* `cmd/agents/*`– примеры агентов ресурсов pacemaker (`golang`)`
* `vendor`      - `golang` зависимости `cmd/cibgen`
* `testdata`    - заранее предоставленные тестовые данные, сгенерированные с помощью `cmd/cibgen`
* `scripts`     - скрипты для выполнения рутинных действий
* `pacemaker`   - исходные коды pacemaker (создаётся после вызова `./scripts/pacemaker/checkout`)
* `docker`      - директория для построения образа локального `docker` контейнера для разработки

## Окружение

Для сборки и работы с `pacemaker`, `golang` утилитами и примерами требуется определённое окружение. Для упрощения работы и исключения засорения основной операционной системы это окружение предоставляется в виде `docker` контейнера, инструкция по подготовке которого приведена ниже.

### Особенности mac os

В виду того, что технология контейнерезации, использованная в `docker`, опирается на механизмы, доступные только в ос `Linux`, контейнеры`docker` на `mac os` запускаются в виртуальной машине с ос`Linux`.

Для работы с проектом корневая директория проекта монтируется в файловую систему запущенного `docker` контейнера. В `mac os` это требует проброса корневой директории проекта в данную виртуальную машину.

Технология, использованная для проброса, не позволяет получить высокой производительности при доступе к содержимому корневой директории проекта из окружения внутри `docker` контейнера.

Существуют способы обойти данную проблему, но в данном случае скорость доступа к содержимому корневой директории проекта не является блокирующим фактором и, поэтому, какие-либо способы решить проблему скорости доступа не применяются.


### Docker контейнер

Сборка локального образа `docker` контейнера для сборки `pacemaker` и запуска тестов

    ./scripts/docker/build

Образ контейнера получит тег `pcmkrtoolkit:latest`

Запуск локального `docker` контейнера для работы с с `pacemaker` и запуска тестов

    ./scripts/docker/run

После запуска контейнер будет доступен по имени `pcmkrtoolkit`

## Сборка проекта

### Сборка pacemaker (внутри docker контейнера)

Получение исходников `pacemaker` (требуется запустить один раз)

    ./scripts/pacemaker/checkout

Сборка и установка компонент `pacemaker`

    ./scripts/pacemaker/build

> Cкрипт `./scripts/pacemaker/build` соберёт `pacemaker` и установит его в текущей сессии `docker`.  В случае перезапуска локального `docker` контейнера пересборку `pacemaker` запускать не требуется, для подготовки рабочего окружения достаточно лишь запустить следующую команду
>
>     ./scripts/pacemaker/install

## Работа с pacemaker

### Пример запуска генератора тестовых данных

> 500 блочных волюмов, на 2х нодах по 4 порта каждая - **46074** действий в графе pacemaker

    go run ./cmd/cibgen/       \
        --nodes 2              \
        --ports 4              \
        --pools 2              \
        --volumes-per-pool 255 \
      > testdata/500-volumes-4-ports-2-nodes.cib.xml

> 500 блочных волюмов, на 2х нодах по 8 портов каждая - **86994** действий в графе pacemaker

    go run ./cmd/cibgen/       \
        --nodes 2              \
        --ports 8              \
        --pools 2              \
        --volumes-per-pool 255 \
      > testdata/500-volumes-8-ports-2-nodes.cib.xml

> Эту команду можно запускать вне `docker` в случае, если на системе установлен golang 1.19, в случае запуска внутри `docker`под управлением `mac os` генерация может занять продолжительное время (см. пункт "Особенности mac os").

### Запуск проверки соответствия тестовых данных XML схеме

    crm_verify -x testdata/500-volumes-4-ports-2-nodes.cib.xml -V

В случае обнаружения несоответствия тестовых данных XML схеме будет выдано описание проблемы. Описание выдаётся достаточно абстрактное в виду особенностей работы логики проверки, но примерно что пошло не так понять всё-таки можно по предоставленным идентификаторам элементов, которые не соответствуют схеме. Обычно это идентификаторы, указанные в выводе `crm_verify` первыми, за ними могут быть указаны дополнительные идентификаторы, якобы тоже не соответствующие XML схеме, но они обычно XML схеме соответствуют и их появление в выводе `crm_verify` является результатом ранее выявленного несоответствия.

### Запуск расчёта графа (со сбором flame-graph)

Расчёт графа действий по заданной конфигурации можно произвести из командной строки с помощью утилиты `crm_simulate`. Она позволяет произвести все те же самые действия, что и сервис `pacemaker-schedulerd`. В дополнение, с помощью неё можно симулировать различные события, обычно происходящие в кластере `pacemaker`, или сохранить полученный граф в `DOT` формате.

Расчёт графа действий производится функций `pcmk__schedule_actions` из файла `pacemaker/lib/pacemaker/pcmk_sched_allocate.c` и состоит из нескольких этапов (хорошо видно на flame graph):

* распаковка заданной конфигурации в представление в памяти – функции `*_unpack_*`
* назначение распределение ресурсов по нодам – функция `allocate_resources`
* создание возможных действий для каждого ресурса – функция `schedule_resource_actions`
* расчёт последовательности действий – функция `pcmk__apply_orderings`

Ниже приведён пример использования утилиты `crm_simulate`со снятием `flame graph` вызовов.

> 500 блочных волюмов, на 2х нодах по 4 порта каждая - **46074** действий в графе pacemaker

    PCMK_stderr=true \
    PCMK_trace_files=pcmk_sched_allocate.c\
    PCMK_trace_functions=allocate_resources,schedule_resource_actions,pcmk__apply_orderings \
      ./scripts/flame-graph                              \
        -o ./testdata/500-volumes-4-ports-2-nodes.flame-graph.svg \
          crm_simulate                                                      \
                -V                                                          \
                -x         ./testdata/500-volumes-4-ports-2-nodes.cib.xml   \
                -G         ./testdata/500-volumes-4-ports-2-nodes.graph.xml \
                -D         ./testdata/500-volumes-4-ports-2-nodes.graph.dot \
                --node-up  node0                                            \
                --node-up  node1

> 500 блочных волюмов, на 2х нодах по 8 портов каждая - **86994** действий в графе pacemaker

    PCMK_stderr=true \
    PCMK_trace_files=pcmk_sched_allocate.c \
    PCMK_trace_functions=allocate_resources,schedule_resource_actions,pcmk__apply_orderings \
      ./scripts/flame-graph                              \
        -o ./testdata/500-volumes-8-ports-2-nodes.flame-graph.svg \
            crm_simulate                                                    \
                -V                                                          \
                -x         ./testdata/500-volumes-8-ports-2-nodes.cib.xml   \
                -G         ./testdata/500-volumes-8-ports-2-nodes.graph.xml \
                -D         ./testdata/500-volumes-8-ports-2-nodes.graph.dot \
                --node-up   node0 \
                --node-up   node1

> flame graph сохраняется в интерактивную `svg` схему, которую можно открыть в web броузере.

### Запуск тестовых агентов ресурсов pacemaker

В директории `cmd/agents` находится исходный код тестовых агентов ресурсов pacemaker. Данные агенты не делают ничего полезного и умеют только обрабатывать команды start, stop, promote, demote, monitor

По факту, исходный код всех агентов – одинаковый. Но этого достаточно, чтобы продемонстрировать интеграцию с сервисом `pacemaker-execd`

#### Установка тестовых агентов

Для установки тестовых агентов следует запустить следующую команду

    ./scripts/agents/install

Исходный код агентов будет скомпилирован и исполняемые файлы будут установлены в директорию `/usr/lib/ocf/resource.d/yadro`.

> в случае перезапуска `docker` контейнера, установку агентов `pacemaker` нужно производить заново.

#### Запуск pacemaker-execd в фоновом режиме

Для проверки интеграции агентов с логикой `pacemaker-execd` потребуется запустить последнего в фоновом режиме в отдельном окне (вкладке терминала). Для этого нужно выполнить следующую команду **вне** `docker` контейнера `pcmkrtoolkit`

    ./scripts/docker/exec ./scripts/pacemaker/execd-run

В результате этого `pacemaker-execd` будет запущен **внутри** `docker` контейнера `pcmkrtoolkit` в фоновом режиме.

#### Выполнение операций с агентами

Для выполнения операций с агентами предназначен скрипт `./scripts/pacemaker/execd-send-command`

##### Получение метаданных агента ресурсов

    ./scripts/pacemaker/execd-send-command metadata -p yadro -k pool

##### Регистрация ресурса в pacemaker-execd

> Это обязательный шаг перед тем, как можно будет выполнять описанные далее команды

    ./scripts/pacemaker/execd-send-command register -r some-pool -p yadro -k pool

##### Получение информации о ресурсе в pacemaker-execd

    ./scripts/pacemaker/execd-send-command get_rsc_info -r some-pool

##### Выполнение операции запуска ресурса

    ./scripts/pacemaker/execd-send-command start -r some-pool

##### Выполнение операции запуска повторяющейся операции опроса состояния ресурса

> После завершения следующей команды `pacemaker-execd` будет посылать команду monitor агенту ресурса с периодичностью 10 секунд (интервал жёстко задан в скрипте `./scripts/pacemaker/execd-send-command`)

    ./scripts/pacemaker/execd-send-command monitor -r some-pool

##### Отмена повторяющейся операции опроса состояния ресурса

> После завершения следующей команды `pacemaker-execd` перестанет периодически посылать команду monitor агенту ресурса, если таковой процесс был инициирован ранее

    ./scripts/pacemaker/execd-send-command cancel -r some-pool

##### Выполнение операции остановки ресурса

> После завершения следующей команды `pacemaker-execd` перестанет периодически посылать команду monitor агенту ресурса, если таковой процесс был инициирован ранее и вызовет операцию stop агента ресурса

    ./scripts/pacemaker/execd-send-command stop -r some-pool

## Переменные окружения, влияющие на работу сервисов и утилит pacemaker

* `PCMK_stderr`      – включить отладочный вывод в `stderr` (доступные значения: `true`, `false`, `on`, `off`, `yes`, `no`, `0`, `1`)
* `PCMK_trace_files` - включить весь отладочный `trace` вывод из одного или нескольких файлов (допустимые значения: строка, содержащая имена исходных файлов, разделённых запятой)
* `PCMK_trace_functions` – включить весь отладочный `trace` вывод из одной или нескольких функций (допустимые значения: строка, содержащая имена функций исходного кода, разделённых запятой)





