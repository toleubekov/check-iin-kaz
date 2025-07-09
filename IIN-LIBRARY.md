# 📚 check-iin-kaz - Go библиотека для валидации ИИН Казахстана

[![Go Version](https://img.shields.io/badge/Go-1.24.2+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Легковесная и быстрая Go библиотека для валидации Индивидуальных Идентификационных Номеров (ИИН) Республики Казахстан с извлечением персональной информации.

## 🚀 Установка

```bash
go get github.com/toleubekov/check-iin-kaz
```

## 📖 Быстрый старт

### Основное использование

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/toleubekov/check-iin-kaz/iin"
)

func main() {
    // Валидация с извлечением полной информации
    info, err := iin.Validate("031231500123")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Валидный: %t\n", info.Valid)
    fmt.Printf("Пол: %s\n", info.Sex)
    fmt.Printf("Дата рождения: %s\n", info.DateOfBirth)
    fmt.Printf("Век: %d\n", info.Century)
    fmt.Printf("Региональный код: %d\n", info.RegionCode)
}
```

### Быстрая проверка валидности

```go
import "github.com/toleubekov/check-iin-kaz/iin"

if iin.IsValid("031231500123") {
    fmt.Println("ИИН корректный")
}
```

### Извлечение отдельных компонентов

```go
// Извлечение пола
sex, err := iin.ExtractSex("031231500123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Пол: %s\n", sex) // male или female

// Извлечение даты рождения
dateOfBirth, err := iin.ExtractDateOfBirth("031231500123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Дата рождения: %s\n", dateOfBirth) // DD.MM.YYYY
```

## 📋 API Reference

### Основные функции

#### `iin.Validate(iinStr string) (*IINInfo, error)`

Выполняет полную валидацию ИИН и извлекает всю доступную информацию.

**Возвращает:**
- `*IINInfo` - структура с информацией об ИИН
- `error` - ошибка, если ИИН невалидный

#### `iin.IsValid(iinStr string) bool`

Быстрая проверка корректности ИИН без извлечения дополнительной информации.

#### `iin.ValidateAndExtract(iinStr string) (bool, string, string, error)`

Совместимая функция, возвращающая: валидность, пол, дату рождения, ошибку.

### Структура IINInfo

```go
type IINInfo struct {
    Valid       bool   `json:"valid"`            // Валидность ИИН
    Sex         string `json:"sex"`              // "male" или "female"
    DateOfBirth string `json:"date_of_birth"`    // DD.MM.YYYY
    Century     int    `json:"century"`          // Век рождения (18, 19, 20)
    RegionCode  int    `json:"region_code"`      // Региональный код (1000-9999)
}
```

## 🔍 Формат ИИН

ИИН состоит из 12 цифр в формате: `YYMMDDVNNNNK`

- `YY` - год рождения (00-99)
- `MM` - месяц рождения (01-12)
- `DD` - день рождения (01-31)
- `V` - век и пол:
  - `1` - мужчина, 1800-1899
  - `2` - женщина, 1800-1899
  - `3` - мужчина, 1900-1999
  - `4` - женщина, 1900-1999
  - `5` - мужчина, 2000-2099
  - `6` - женщина, 2000-2099
- `NNNN` - порядковый номер (1000-9999)
- `K` - контрольная цифра

## 🏗️ Примеры интеграции

### HTTP сервер

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/toleubekov/check-iin-kaz/iin"
)

func validateIINHandler(w http.ResponseWriter, r *http.Request) {
    iinStr := r.URL.Query().Get("iin")
    
    info, err := iin.Validate(iinStr)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(info)
}

func main() {
    http.HandleFunc("/validate", validateIINHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Middleware для валидации

```go
func IINValidationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        iinStr := r.Header.Get("X-IIN")
        
        if iinStr != "" && !iin.IsValid(iinStr) {
            http.Error(w, "Invalid IIN", http.StatusBadRequest)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### Валидация в struct

```go
type User struct {
    Name string `json:"name"`
    IIN  string `json:"iin"`
}

func (u *User) Validate() error {
    if !iin.IsValid(u.IIN) {
        return errors.New("invalid IIN")
    }
    return nil
}
```

## ⚡ Производительность

- **Sub-millisecond** валидация ИИН
- **Zero allocation** для простых проверок
- **Thread-safe** - можно использовать в concurrent окружении
- **Малый размер** - минимальные зависимости

## 🧪 Тестирование

```go
package main

import (
    "testing"
    
    "github.com/toleubekov/check-iin-kaz/iin"
)

func TestIINValidation(t *testing.T) {
    tests := []struct {
        iin   string
        valid bool
    }{
        {"031231500123", true},
        {"400701111111", true},
        {"123456789012", false},
        {"12345678901", false},  // короткий
    }
    
    for _, test := range tests {
        result := iin.IsValid(test.iin)
        if result != test.valid {
            t.Errorf("IIN %s: expected %t, got %t", test.iin, test.valid, result)
        }
    }
}
```

## 🔧 Обработка ошибок

Библиотека возвращает понятные ошибки на русском языке:

- `"длина ИИН должна быть равна 12 символам"`
- `"ИИН должен состоять только из цифр"`
- `"некорректная контрольная сумма ИИН"`
- `"неверный месяц рождения"`
- `"неверный день рождения"`
- `"дата рождения не может быть в будущем"`

## 🤝 Совместимость

- ✅ **Go 1.18+**
- ✅ **Без внешних зависимостей**
- ✅ **Thread-safe**
- ✅ **Cross-platform**

## 📈 Benchmarks

```
BenchmarkValidate-8           1000000      1123 ns/op       0 B/op       0 allocs/op
BenchmarkIsValid-8           2000000       856 ns/op       0 B/op       0 allocs/op
BenchmarkExtractSex-8       10000000       145 ns/op       0 B/op       0 allocs/op
```

## 📄 Лицензия

MIT License - см. [LICENSE](LICENSE) файл.

## 👨‍💻 Автор

**Zhandos Toleubekov** - [GitHub](https://github.com/toleubekov)

---

**Сделано с ❤️ для Go сообщества Казахстана**