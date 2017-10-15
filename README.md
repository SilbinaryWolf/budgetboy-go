## Budget Boy

A simple tool to export CSVs with total spending per week using bank CSV export data.
Built for personal use.

![subway-boy](https://user-images.githubusercontent.com/3859574/31581224-b358111e-b1b1-11e7-97ec-0617c24710bc.jpg)

## Quick Start

1) Create a new config.txt file in the same dir as Budget Boy.
2) Add items underneath categories that partially match the transaction item.
3) Place "CSVData.csv" in same dir as Budget Boy.
4) Run.
5) Check "output" folder, there will be various CSV's named like "2017_October_8th.csv"

## Config File

Setup a `config.txt` file in the same directory as the exe.

```
###########################
Earning Per Week: $125.00
Rent Per Week: $30.00
# Throw an error if it hits an item without a category?
Disallow Uncategorized: true
# Show CSV in console output
Print In Console: true
###########################

Category: Fast Food
###########################
SUBWAY
KFC
WHITE CASTLE
GRILLD

Category: Gas Station
###########################
7-ELEVEN

Category: Groceries
###########################
COLES
WOOLWORTHS

Category: Bills / Rent
###########################
Account Fee
MYKI
TELSTRA
NETFLIX

Category: _
###########################
# Ignore these transactions
Landlord Payment
Transfer from
Transfer to
```

## Expected CSVData Format

Date, Cost, Name, Balance After Cost

```
09/10/2017,"-25.20","HUNGRY JACKS SYDNEY AUS Card xx1441 Value Date: 03/10/2017","+160.87"
09/10/2017,"-19.99","WHITE CASTLE SYDNEY AUS Card xx1441 Value Date: 05/10/2017","+170.87"
09/10/2017,"-11.00","GRILL'D SYDNEY VI AUS Card xx1441 Value Date: 01/10/2017","+200.87"
```
