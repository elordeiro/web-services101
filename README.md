This is a simple web service that takes a receipt in JSON format, calculates points based on items purchased and stores a unique ID that can be accessed via an endpoint.

---

To run the web service, run the following command:  
`./ReceiptProcessor`  
or  
`go run .`

---

---

To test the web service:  
While the web service is running, open a new terminal window and run  
`./runTests.sh`  
or  
use curl to send your own custom tests  
or  
open a browser and navigate to `http://localhost:8080/` and use the web interface to test the service  
(note that the time and date fields are filled in automatically with the current time and date)

---
