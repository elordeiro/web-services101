let items = [
    {
        item: "Mountain Dew 12 PK",
        price: "6.49",
    },
    {
        item: "Emils Cheese Pizza",
        price: "12.25",
    },
    {
        item: "Knorr Creamy Chicken",
        price: "1.26",
    },
    {
        item: "Doritos Nacho Cheese",
        price: "3.35",
    },
    {
        item: "Klarbrunn 12-PK 12 FL OZ",
        price: "12.00",
    },
    {
        item: "Gatorade",
        price: "2.25",
    },
    {
        item: "Pepsi - 12-oz",
        price: "1.25",
    },
    {
        item: "Dasani",
        price: "1.40",
    },
];

let stores = ["M&M Corner Market", "Target", "Walgreens", "Walmart"];

let table = document.getElementById("table");

items.forEach((item) => {
    let row = table.insertRow(table.rows.length - 2);
    let cell1 = row.insertCell(0);
    let cell2 = row.insertCell(1);
    let cell3 = row.insertCell(2);
    cell1.classList.add("row-style");
    cell2.classList.add("row-style");
    cell3.classList.add("row-style");
    cell3.style.textAlign = "center";
    cell1.innerHTML = item.item;
    cell2.innerHTML = item.price;
    let input = Object.assign(document.createElement("input"), {
        type: "number",
        value: 0,
        min: 0,
        id: item.item,
        style: "width: 50px;",
    });
    cell3.appendChild(input);
});

let storeSelect = document.getElementById("store-select");
stores.forEach((store) => {
    let option = document.createElement("option");
    option.text = store;
    storeSelect.add(option);
});

function newReceipt() {
    return {
        retailer: "",
        purchaseDate: "",
        purchaseTime: "",
        items: [],
        total: 0.0,
    };
}

function newItem(description, price) {
    return {
        shortDescription: description,
        price: price,
    };
}

let receiptIDs = {};

let submit = document.getElementById("submit");
submit.addEventListener("click", () => {
    let receipt = newReceipt();
    receipt.retailer = storeSelect.options[storeSelect.selectedIndex].text;
    let date = new Date();
    receipt.purchaseDate = date.toISOString().split("T")[0];
    receipt.purchaseTime = date.toTimeString().split(" ")[0].slice(0, 5);
    items.forEach((item) => {
        let input = document.getElementById(item.item);
        let quantity = parseInt(input.value);
        for (let i = 0; i < quantity; i++) {
            receipt.items.push(newItem(item.item, item.price));
        }
    });
    receipt.total = receipt.items
        .reduce((total, item) => total + parseFloat(item.price), 0)
        .toFixed(2);
    console.log(receipt);
    fetch("http://localhost:8080/receipts/process", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(receipt),
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Receipt was invalid.");
            }
            return response.json();
        })
        .then((data) => {
            console.log("Success:", data);
            getPoints(data.id);
        })
        .catch((error) => {
            console.error("Error:", error);
        });
});

function getPoints(id) {
    fetch(`http://localhost:8080/receipts/${id}/points`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Receipt ID was not found.");
            }
            return response.json();
        })
        .then((data) => {
            let points = data.points;
            receiptIDs[id] = points;
            console.log(receiptIDs);
            updateReceiptsTable();
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function updateReceiptsTable() {
    let table = document.getElementById("receipts-table");
    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    for (let id in receiptIDs) {
        let row = table.insertRow(-1);
        let cell1 = row.insertCell(0);
        let cell2 = row.insertCell(1);
        cell1.classList.add("row-style");
        cell2.classList.add("row-style");
        cell1.innerHTML = id;
        cell2.innerHTML = receiptIDs[id];
    }
}
