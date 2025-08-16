# Logical Interview Prep: Vending Machine Logic

**Custom Gem**

Excellent. Let's tackle a new problem. This one is less about complex algorithms and more about modeling state and behavior in a clean, robust way.

## Scenario: Vending Machine Logic

Your task is to design and implement the logic for a simple vending machine in Go. The machine needs to manage inventory, handle money, and process user selections.

### Core Requirements

- **Inventory**: The machine stocks several products. Each product has a name (e.g., "Cola"), a price, and a quantity.

- **Money**: The machine accepts a predefined set of coin denominations (e.g., 1, 5, 10, 25). It also needs to track the total money inserted by a user for the current transaction.

- **State**: The machine must manage its internal state, including the quantity of each item and the current user's balance.

### Required Functionality

You should design a `VendingMachine` struct with the following methods:

#### `InsertCoin(coinValue int) error`

- Accepts a coin value.
- If the coin is a valid denomination, it adds it to the user's current balance.
- If the coin is invalid, it should return an error.

#### `SelectProduct(productName string) (change int, err error)`

This is the core of the operation. It should perform these checks in order:

- Does the product exist?
- Is the product in stock?
- Has the user inserted enough money?

If any check fails, return a descriptive error. The user's inserted money should remain for them to make another choice or cancel.

If all checks pass:
- Dispense the item (i.e., decrement its quantity)
- Calculate the correct change to return
- Reset the user's balance to zero
- Return the change

#### `Cancel() int`

- Cancels the current transaction.
- Returns the total amount of money the user had inserted.
- Resets the user's balance to zero.

---

### Design Considerations

How would you structure the `VendingMachine` and any related data types?

Think about:
- The fields you'd need to track the state
- How you'd handle the different success and failure paths in the `SelectProduct` method
- How to separate concerns like product management, balance tracking, and validation