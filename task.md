
# Задание: Игра без права на ошибку

## Описание
Мы пишем игру, где любая ошибка приводит к трагическому финалу. В игре есть команды, объекты игрового мира, шаги, игрок и сама игра. Ваша задача — модернизировать обработку ошибок, добавив структуру и соответствующие типы, а также реализовать советы для игрока на основе возникших ошибок.



## Команды, объекты, шаги
### Команды
```go
// label - уникальное наименование
type label string

// command - команда, которую можно выполнять в игре
type command label

// список доступных команд
var (
    eat  = command("eat")
    take = command("take")
    talk = command("talk to")
)
```

### Объекты игрового мира
```go
// thing - объект, который существует в игре
type thing struct {
    name    label
    actions map[command]string
}

// supports() возвращает true, если объект поддерживает команду action
func (t thing) supports(action command) bool {
    _, ok := t.actions[action]
    return ok
}

// полный список объектов в игре
var (
    apple = thing{"apple", map[command]string{
        eat:  "mmm, delicious!",
        take: "you have an apple now",
    }}
    bob = thing{"bob", map[command]string{
        talk: "Bob says hello",
    }}
    coin = thing{"coin", map[command]string{
        take: "you have a coin now",
    }}
    mirror = thing{"mirror", map[command]string{
        take: "you have a mirror now",
        talk: "mirror does not answer",
    }}
    mushroom = thing{"mushroom", map[command]string{
        eat:  "tastes funny",
        take: "you have a mushroom now",
    }}
)
```

### Шаг игры
```go
// step описывает шаг игры: сочетание команды и объекта
type step struct {
    cmd command
    obj thing
}

// isValid() возвращает true, если объект совместим с командой
func (s step) isValid() bool {
    return s.obj.supports(s.cmd)
}
```

---

## Игрок
### Тип player
```go
// player - игрок
type player struct {
    nEaten     int      // количество съеденного
    nDialogs   int      // количество диалогов
    inventory  []thing  // инвентарь
}

// has() возвращает true, если у игрока в инвентаре есть предмет obj
func (p *player) has(obj thing) bool {
    for _, got := range p.inventory {
        if got.name == obj.name {
            return true
        }
    }
    return false
}

// do() выполняет команду cmd над объектом obj от имени игрока
func (p *player) do(cmd command, obj thing) error {
    switch cmd {
    case eat:
        if p.nEaten > 1 {
            return errors.New("you don't want to eat anymore")
        }
        p.nEaten++
    case take:
        if p.has(obj) {
            return fmt.Errorf("you already have a %s", obj)
        }
        p.inventory = append(p.inventory, obj)
    case talk:
        if p.nDialogs > 0 {
            return errors.New("you don't want to talk anymore")
        }
        p.nDialogs++
    }
    return nil
}
```

---

## Игра
### Тип game
```go
// game описывает игру
type game struct {
    player  *player
    things  map[label]int
    nSteps  int
}

// has() проверяет, остались ли в игровом мире указанные предметы
func (g *game) has(obj thing) bool {
    count := g.things[obj.name]
    return count > 0
}

// execute() выполняет шаг step
func (g *game) execute(st step) error {
    if !st.isValid() {
        return fmt.Errorf("cannot %s", st)
    }

    if st.cmd == take || st.cmd == eat {
        if !g.has(st.obj) {
            return fmt.Errorf("there are no %ss left", st.obj)
        }
        g.things[st.obj.name]--
    }

    if err := g.player.do(st.cmd, st.obj); err != nil {
        return err
    }

    g.nSteps++
    return nil
}

// newGame() создает новую игру
func newGame() *game {
    p := newPlayer()
    things := map[label]int{
        apple.name:    2,
        coin.name:     3,
        mirror.name:   1,
        mushroom.name: 1,
    }
    return &game{p, things, 0}
}
```

---

## Новые требования
### Типы ошибок
```go
// invalidStepError - ошибка для несовместимой команды и объекта
type invalidStepError struct {}

// notEnoughObjectsError - ошибка для нехватки объектов в игре
type notEnoughObjectsError struct {}

// commandLimitExceededError - ошибка для превышения лимита действий
type commandLimitExceededError struct {}

// objectLimitExceededError - ошибка для превышения лимита объектов в инвентаре
type objectLimitExceededError struct {}

// gameOverError - ошибка уровня игры
type gameOverError struct {
    nSteps int
    // вложенная ошибка
    innerErr error
}
```

### Функция советов
```go
// giveAdvice() дает совет игроку на основе ошибки
func giveAdvice(err error) string {
    // ...
}
```

---

## Пример использования
```go
func main() {
    gm := newGame()
    steps := []step{
        {eat, apple},
        {talk, bob},
        {take, coin},
        {eat, mushroom},
    }

    for _, st := range steps {
        if err := tryStep(gm, st); err != nil {
            fmt.Println(err)
            fmt.Println(giveAdvice(err))
            os.Exit(1)
        }
    }
    fmt.Println("You win!")
}
```

# Задание: Игра без права на ошибку

## Условие задачи

В данной задаче вы работаете с программой-игрой, где любая ошибка игрока приводит к завершению игры. Основной целью является добавление структурированных типов ошибок, улучшение обработчика ошибок и написание функции для предоставления советов игроку.

---

## Детали

### Команда, объект, шаг

#### Типы данных

```go
// label - уникальное наименование
type label string

// command - команда, которую можно выполнять в игре
type command label

// список доступных команд
var (
    eat  = command("eat")
    take = command("take")
    talk = command("talk to")
)

// thing - объект, который существует в игре
type thing struct {
    name    label
    actions map[command]string
}

// step описывает шаг игры: сочетание команды и объекта
type step struct {
    cmd command
    obj thing
}

```
---

### Игрок

#### Типы данных

```go
// player - игрок
type player struct {
    nEaten     int       // количество съеденного
    nDialogs   int       // количество диалогов
    inventory  []thing   // инвентарь
}

// has() возвращает true, если у игрока
// в инвентаре есть предмет obj
func (p *player) has(obj thing) bool {
    for _, got := range p.inventory {
        if got.name == obj.name {
            return true
        }
    }
    return false
}

// do() выполняет команду cmd над объектом obj
func (p *player) do(cmd command, obj thing) error {
    switch cmd {
    case eat:
        if p.nEaten > 1 {
            return errors.New("you don't want to eat anymore")
        }
        p.nEaten++
    case take:
        if p.has(obj) {
            return fmt.Errorf("you already have a %s", obj)
        }
        p.inventory = append(p.inventory, obj)
    case talk:
        if p.nDialogs > 0 {
            return errors.New("you don't want to talk anymore")
        }
        p.nDialogs++
    }
    return nil
}
```

---

### Игра

#### Типы данных

```go
// game описывает игру
type game struct {
    player *player       // игрок
    things map[label]int // объекты игрового мира
    nSteps int           // количество успешно выполненных шагов
}

// execute() выполняет шаг игры
func (g *game) execute(st step) error {
    if !st.isValid() {
        return fmt.Errorf("cannot %s", st)
    }

    if st.cmd == take || st.cmd == eat {
        if !g.has(st.obj) {
            return fmt.Errorf("there are no %ss left", st.obj)
        }
        g.things[st.obj.name]--
    }

    if err := g.player.do(st.cmd, st.obj); err != nil {
        return err
    }

    g.nSteps++
    return nil
}

// newGame() создает новую игру
func newGame() *game {
    p := newPlayer()
    things := map[label]int{
        apple.name:    2,
        coin.name:     3,
        mirror.name:   1,
        mushroom.name: 1,
    }
    return &game{p, things, 0}
}
```

---

## Задание

### Требуется:

1. **Создать новые типы ошибок:**
    - `invalidStepError`
    - `notEnoughObjectsError`
    - `commandLimitExceededError`
    - `objectLimitExceededError`

2. **Создать тип `gameOverError` для ошибок верхнего уровня:**
    - Содержит количество успешно выполненных шагов.
    - Оборачивает другие ошибки.

3. **Реализовать функцию `giveAdvice(err error) string`:**
    - Возвращает советы игроку на основе типа ошибки.

### Пример:

#### Вывод:

```plaintext
trying to eat apple... OK
trying to talk to bob... OK
trying to take coin... OK
trying to eat mushroom... OK
You win!
```

#### Ошибки:

```plaintext
trying to talk to bob... FAIL
you don't want to talk anymore
exit status 1
```

#### Советы:

```plaintext
things like 'eat bob' are impossible
be careful with scarce mirrors
don't be greedy, 1 apple is enough
```

