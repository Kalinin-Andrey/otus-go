Feature: notification

  Scenario: create event on a next day
    When I create event on a next day with duration in day and Title="Тест оповещения",
    Then I receive notification with ID and Title="Тест оповещения"



