This is an implementation of the API needed in the code challenge https://github.com/engagetech/backend-coding-challenge for EngageTech.

The task
========

Imagine that you come back from 2 weeks of holidays on a Monday. On the team scrum board, assigned to you, two tasks await :

User story 1:
------------

As a user, i want to be able to enter my expenses and have them saved for later.

As a user, in the application UI, I can navigate to an expenses page. On this page, I can add an expense, setting :

The date of the expense

The value of the expense

The reason of the expense

When I click "Save Expense", the expense is then saved in the database. The new expense can then be seen in the list of submitted expenses.

User story 2:
------------

As a user, I want to be able to see a list of my submitted expenses.

As a user, in the application UI, i can navigate to an expenses page. On this page, I can see all the expenses I already submitted in a tabulated list. On this list, I can see :

The date of the expense

The VAT (Value added tax) associated to this expense. VAT is the UK’s sales tax. It is 20% of the value of the expense, and is included in the amount entered by the user.

The reason of the expense

Setup
-----
* Install go, you can follow [this installation guide](https://golang.org/doc/install#install).
* Install mySql, you can follow [this installation guide](https://dev.mysql.com/doc/refman/5.7/en/installing.html)
* Create database "expenses", create user and grant permissions.
> * Login to mySql
> * Create the database: 
```sql
    CREATE DATABASE expenses;
```
> * Create user: 
```sql
    CREATE USER 'userDev'@'%' IDENTIFIED by 'passDev';
``` 
> * Grant permissions 
```sql
    GRANT ALL ON expenses.* TO 'userDev'@'%';
```
* Run build.sh to build the executable.
* Run buildAndRun.sh to build and run the application.

Available endpoints
-------------------
GET  localhost:12345/app/expenses - returns the list of all expenses saved so far

POST localhost:12345/app/expenses - endpoint for adding new expenses, format is json
```json
{"date":"12/11/1988","reason":"why not?","amount":"100.1"}
```
