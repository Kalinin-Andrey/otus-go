# go-rti-testing


##### Установка
```
go get github.com/Kalinin-Andrey/otus-go/calendar
```
##### Запуск в docker-контейнерах
```
docker-compose -f ./deployment/docker-compose.yml up
```
##### Запуск тестов
```
go test -v ./internal/test
```
`+` конфиг функциональных тестов АПИ для Apache jMeter
```
./internal/test/rti-testing.jmx
```

##### Ниже исходное задание

У одного продукта(product) могут быть различные конфигурации компонентов(component). Для его продажи необходимо сформировать предложение(offer). Продукт возможно продать, если у всех обязательных компонентов сформировалась одна простая цена(price). Формирование цен происходит на основании правил применимости(rule applicabilities) из входящих условий (condition). Необходимо разработать сервис, который на вход будет принимать условия, а на выход отдавать предложение. В задании есть тесты, если все тесты успешно проходят, то предложение корректно сформировалось. 

##### Будет плюсом:
1. Реализовать логирование ошибок приложения и запросов к реализованным REST сервисам. 
2. Реализовать по Чистой архитектуре (https://habr.com/ru/post/269589/)
3. Написать Dockerfile сборки и деплоя приложения

##### Технические требования:
1. Приложение должно быть написано на Golang
2. Приложение не должно быть написано с помощью какого-либо фреймворка, однако можно устанавливать для него различные пакеты
3. Результаты запросов должны быть представлены в формате JSON
4. Результат задания должен быть выложен на github, должна быть инструкция по запуску проект

##### Operator(Логические операторы):
* EQ (равно)
* LTE (меньше или равно)
* GTE (больше или равно)

##### Component(Компонент):
В одном продукте может быть несколько компонентов.  
В одном компоненте должна быть только 1 цена с типом COST, если их больше или меньше, то этот компонент невозможно продать.  
Если у компонента задан признак isMain - то в случае его невалидности, продукт невозможно продать.

##### Price(Цена):
Цена, имеет типы COST и DISCOUNT. COST - простая цена, DISCOUNT - скидка в %.  
В компоненте может быть несколько цен, но компонент возможно продать если по результатам расчета осталась только одна 
с типом COST.  
В итоговое предложение необходимо вывести цену на компонент с учетом скидки(DISCOUNT), если такая осталась по результатам расчета.    
Если осталось больше одной скидки, то применяется максимальная.

##### RuleApplicabilities(Правила применимости): 
Правило считается валидным, только в том случае, если все его значения удовлетворяют условиям.

##### Condition(Условия):
Условия для расчета правил применимости Condition.RuleName=RuleApplicability.CodeName

##### Offer(Итоговое предложение): 
В итоговое предложение должны попасть все компоненты, которые имеют одну простую цену.
Если есть скидки их небходимо применить к цене того компонента, на котором она задана. Если скидок больше 1, то применяется максимальная.  
В TotalCost должна попасть сумма всех простых цен с учетом скидки.

#### Пример  
##### Запрос:  
```json
{
  "name": "Игровой",
  "components": [
    {
      "isMain": true,
      "name": "Интернет",
      "prices": [
        {
          "cost": 100,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "adsl"
            },
            {
              "codeName": "internetSpeed",
              "operator": "EQ",
              "value": "10"
            }
          ]
        },
        {
          "cost": 150,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "adsl"
            },
            {
              "codeName": "internetSpeed",
              "operator": "EQ",
              "value": "15"
            }
          ]
        },
        {
          "cost": 500,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "xpon"
            },
            {
              "codeName": "internetSpeed",
              "operator": "EQ",
              "value": "100"
            }
          ]
        },
        {
          "cost": 900,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "xpon"
            },
            {
              "codeName": "internetSpeed",
              "operator": "EQ",
              "value": "200"
            }
          ]
        },
        {
          "cost": 200,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "fttb"
            },
            {
              "codeName": "internetSpeed",
              "operator": "EQ",
              "value": "30"
            }
          ]
        },
        {
          "cost": 400,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "fttb"
            },
            {
              "codeName": "internetSpeed",
              "operator": "EQ",
              "value": "50"
            }
          ]
        },
        {
          "cost": 600,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "fttb"
            },
            {
              "codeName": "internetSpeed",
              "operator": "EQ",
              "value": "200"
            }
          ]
        },
        {
          "cost": 10,
          "priceType": "DISCOUNT",
          "ruleApplicabilities": [
            {
              "codeName": "internetSpeed",
              "operator": "GTE",
              "value": "50"
            }
          ]
        }
      ]
    },
    {
      "name": "ADSL Модем",
      "prices": [
        {
          "cost": 300,
          "priceType": "COST",
          "ruleApplicabilities": [
            {
              "codeName": "technology",
              "operator": "EQ",
              "value": "adsl"
            }
          ]
        }
      ]
    }
  ]
}
```  
##### Условия:  
```json
[
  {
    "ruleName": "technology",
    "value": "xpon"
  },
  {
    "ruleName": "internetSpeed",
    "value": "200"
  }
]
```  
##### Результат расчета:  
```json
{
  "name": "Игровой",
  "components": [
    {
      "name": "Интернет",
      "isMain": true,
      "prices": [
        {
          "cost": 765
        }
      ]
    }
  ],
  "totalCost": {
    "cost": 765
  }
}
```   

